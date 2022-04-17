package models

import (
	"bytes"
	"cliente/config"
	"cliente/utils"
	"crypto/aes"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Email      string             `bson:"email"`
	ServerKey  []byte             `bson:"server_key"`
	PublicKey  []byte             `bson:"public_key"`
	PrivateKey []byte             `bson:"private_key"`
	Rol        string             `bson:"rol"`
	Kaes       []byte             `bson:"Kaes"`
	Token      string             `bson:"token,omitempty"`
}

type DataUser struct {
	Proyecto Proyect
	Listas   []List
}

//Contiene los proyectos y listas del usuario
var DatosUsuario []DataUser

//Contiene los datos del usuario
var UserSesion User

//Descifra la clave privada del usuario
func GetPrivateKeyUser() *rsa.PrivateKey {
	return utils.PemToPrivateKey(utils.DescifrarAES(UserSesion.Kaes, UserSesion.PrivateKey))
}

//Generamos a partir de un hash del usuario y contraseña:  Kservidor 16 Bytes, IV 16 Bytes y Kaes 32 Bytes
func HashUser(user_pass []byte) ([]byte, []byte, []byte) {
	hash := sha512.Sum512(user_pass)
	Kservidor := hash[:16]
	IV := hash[aes.BlockSize : aes.BlockSize*2]
	Kaes := hash[aes.BlockSize*2:]
	return Kservidor, IV, Kaes
}

// func CheckUserExist(email string) bool {
// 	url := config.URLbase + "users/" + email
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	client := utils.GetClientHTTPS()
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode == 400 {
// 		return false
// 	} else {
// 		return true
// 	}
// }

//Registro del usuario
func Register(email string, password string) string {

	//Envio los datos del registro
	user_pass := []byte(email + password)
	Kservidor, IV, Kaes := HashUser(user_pass)
	//Obtenemos las Claves RSA
	privateKey, publicKey := utils.GeneratePrivatePublicKeys()
	//La clave publica la almaceno directamente en formato PEM
	publicKeyPEM := utils.PublicKeyToPem(&publicKey)
	//La clave privada primero la paso a PEM
	privateKeyPem := utils.PrivateKeyToPem(privateKey)
	//La cifro con AES
	privateKeyCipher := utils.CifrarAES(Kaes, IV, privateKeyPem)
	//La almaceno cifrada
	UserSesion.PrivateKey = privateKeyCipher
	//Guardo Kaes para la sesion del usuario
	UserSesion.Kaes = Kaes
	user := User{
		Email:      email,
		ServerKey:  Kservidor,
		PublicKey:  publicKeyPEM,
		PrivateKey: privateKeyCipher,
	}
	//Enviamos los datos al servidor
	resultado := RegisterServer(user)
	return resultado
	// switch userIDstring {
	// case "serverOFF":
	// 	return userIDstring, true
	// case "Error":
	// 	return userIDstring, true
	// case "Duplicado":
	// 	return userIDstring, true
	// default:
	// 	id, _ := primitive.ObjectIDFromHex(userIDstring)
	// 	UserSesion.ID = id
	// 	return userIDstring, false
	// }
}

//Registro del usuario en el servidor
func RegisterServer(user User) string {
	//Convertimos el user de tipo objeto GO a un JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}
	//Preparo la peticion POST
	url := config.URLbase + "signup"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(userJSON))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := utils.GetClientHTTPS()
	//Realizo la peticion
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "serverOFF"
	}
	defer resp.Body.Close()
	//En caso de fallo del registro del usuario en el servidor
	if resp.StatusCode == 400 {
		return "Error"
	} else if resp.StatusCode == 409 {
		return "Duplicado"
	} else {
		//Si todo fue correcto en el servidor devuelvo el id del usuario creado
		var resultado string
		json.NewDecoder(resp.Body).Decode(&resultado)
		return resultado
	}
}

//Login del usuario. Para el login solo envio al servidor el email y la Kservidor la cual se comprueba alli y si es correcta me devuelve todos los datos del usuario
func LogIn(email string, password string) bool {
	user_pass := []byte(email + password)
	Kservidor, _, Kaes := HashUser(user_pass)
	userLogin := User{
		Email:     email,
		ServerKey: Kservidor,
	}
	UserSesion = LogInServer(userLogin)
	if UserSesion.Email != "" {
		UserSesion.Kaes = Kaes
		return true
	} else {
		return false
	}
}

//Login del usuario en el servidor
func LogInServer(userLogin User) User {
	userJSON, err := json.Marshal(userLogin)
	if err != nil {
		fmt.Println(err)
	}
	url := config.URLbase + "login"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(userJSON))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := utils.GetClientHTTPS()
	//Realizo la peticion
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		return User{}
	} else {
		//Si todo fue correcto en el servidor devuelvo el id del usuario creado
		var resultado User
		json.NewDecoder(resp.Body).Decode(&resultado)

		//Asigno el token que genero el servidor
		token := resp.Header.Get("token")
		fmt.Println(token)
		UserSesion.Token = token
		return resultado
	}
}

//Cerrar la sesion
func LogOut() {
	UserSesion = User{}
	DatosUsuario = []DataUser{}
}

//Recupero todos los usuarios
func GetUsers() ([]User, bool) {
	var usersResponse []User
	url := config.URLbase + "users"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req = AddTokenHeader(req)
	client := utils.GetClientHTTPS()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 400:
		fmt.Println("Ningun usuario encontrado")
		return usersResponse, false
	case 401:
		fmt.Println("Token Expirado")
		return usersResponse, true
	default:
		json.NewDecoder(resp.Body).Decode(&usersResponse)
		return usersResponse, false
	}
}

//Recupero un usuario por su email
func GetUserByEmail(userEmail string) User {
	var usersResponse User
	url := config.URLbase + "users/" + userEmail
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req = AddTokenHeader(req)
	client := utils.GetClientHTTPS()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	//Compruebo si no hay ningun usuario
	if resp.StatusCode == 404 {
		fmt.Println("Ningun usuario encontrado")
		return usersResponse
	} else {
		json.NewDecoder(resp.Body).Decode(&usersResponse)
		return usersResponse
	}
}

//Obtengo las relaciones junto a los proyectos y sus listas asociadas para el usuario
func GetUserProyectsLists() {
	//Limpio los datos del usuario
	DatosUsuario = []DataUser{}
	//Obtengo mi clave privada
	privateKey := GetPrivateKeyUser()
	//Recupero las relaciones
	relations := GetProyectsListsByUser(UserSesion.Email)
	if len(relations) > 0 {
		//Proyectos
		for i := 0; i < len(relations); i++ {
			proyecto := GetProyect(relations[i].ProyectID.Hex())
			//Descifro la clave del proyecto con la clave privada del usuario
			proyectKey := utils.DescifrarRSA(privateKey, relations[i].ProyectKey)
			//Desciframos el proyecto
			proyectoDescifrado := DescifrarProyecto(proyecto, proyectKey)

			//Listas
			var lists []List
			//Por cada lista del proyecto la recupero descifrada usando mi clave privada para descifrar la clave de descifrado de la lista
			for j := 0; j < len(relations[i].Lists); j++ {
				list := GetUserList(relations[i].Lists[j].ListID, relations[i].Lists[j].ListKey, privateKey)
				lists = append(lists, list)
			}
			datos := DataUser{
				Proyecto: proyectoDescifrado,
				Listas:   lists,
			}
			DatosUsuario = append(DatosUsuario, datos)
		}
	}
}

//Elimina a un usuario del sistema borrandolo de todo proyectos, listas y tareas
func DeleteUser(userEmail string) bool {
	DatosUsuario = []DataUser{}
	//Recupero las relaciones junto a los proyectos y las listas Mejorable el pensar en llamar a una funcion que no descifre todo porque no lo necesitamos
	GetUserProyectsLists()
	for i := 0; i < len(DatosUsuario); i++ {
		//Para cada Proyecto miro si el proyecto no tiene mas usuarios que el usuario a borrar
		if len(DatosUsuario[i].Proyecto.Users) == 1 {
			//Borro el proyecto
			DeleteProyect(DatosUsuario[i].Proyecto.ID)
		} else {
			//Quito al usuario del array Users del proyecto
			DeleteUserProyect(DatosUsuario[i].Proyecto.ID, userEmail)
			//Por cada una de las listas del proyecto
			for j := 0; j < len(DatosUsuario[i].Listas); j++ {
				//Compruebo si la lista solo tiene de usuario a dicho usuario
				if len(DatosUsuario[i].Listas[j].Users) == 1 {
					//Borro la lista
					DeleteList(DatosUsuario[i].Listas[j].ID)
				} else {
					//Quito al usuario del array Users de la Lista y de las tareas
					DeleteUserList(DatosUsuario[i].Listas[j].ID, userEmail)
				}
			}
		}
	}
	//Borro las relaciones
	DeleteUserRelation(userEmail)
	//Borro al usuario
	DeleteUserByEmail(userEmail)
	return true
}

//Elimino al usuario del sistema
func DeleteUserByEmail(userEmail string) bool {
	url := config.URLbase + "users/" + userEmail
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req = AddTokenHeader(req)
	client := utils.GetClientHTTPS()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		fmt.Println("El usuario no pudo ser eliminado sistema")
		return false
	} else {
		return true
	}
}

//Recupero los emails de todos los usuarios del sistema
func GetEmails() ([]string, bool) {
	users, tokenExpire := GetUsers()
	if tokenExpire {
		return []string{}, tokenExpire
	}
	var usersEmails []string
	for i := 0; i < len(users); i++ {
		usersEmails = append(usersEmails, users[i].Email)
	}
	return usersEmails, tokenExpire
}

//Recupero la clave publica de un usuario
func GetPublicKey(userEmail string) *rsa.PublicKey {
	publicKeyUserPem := GetUserByEmail(userEmail).PublicKey
	publicKey := utils.PemToPublicKey(publicKeyUserPem)
	return publicKey
}

func AddTokenHeader(req *http.Request) *http.Request {
	req.Header.Add("Authorization", "Bearer "+UserSesion.Token)
	return req
}
