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
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Email      string             `bson:"email"`
	ServerKey  []byte             `bson:"server_key"`
	PublicKey  []byte             `bson:"public_key"`
	PrivateKey []byte             `bson:"private_key"`
	Rol        string             `bson:"rol"`
	Kaes       []byte             `bson:"Kaes"`
}

type DatosUser struct {
	Proyecto Proyect
	Listas   []List
}

//Contiene los proyectos y listas del usuario
var DatosUsuario []DatosUser

//Contiene los datos del usuario
var UserSesion User

//Descifra la clave privada del usuario
func GetPrivateKeyUser() *rsa.PrivateKey {
	return utils.PemToPrivateKey(utils.DescifrarAES(UserSesion.Kaes, UserSesion.PrivateKey))
}

//Generamos a partir de un hash del usuario y contraseña :  Kservidor 16 Bytes, IV 16 Bytes y Kaes 32 Bytes
func HashUser(user_pass []byte) ([]byte, []byte, []byte) {

	hash := sha512.Sum512(user_pass)
	Kservidor := hash[:16]
	IV := hash[aes.BlockSize : aes.BlockSize*2]
	Kaes := hash[aes.BlockSize*2:]
	return Kservidor, IV, Kaes
}

//Registro del usuario
func Register(email string, password string) bool {
	user_pass := []byte(email + password)
	Kservidor, IV, Kaes := HashUser(user_pass)
	KservidorHash, err := bcrypt.GenerateFromPassword(Kservidor, 12)
	if err != nil {
		fmt.Println("Error al hashear")
	}
	//Asigno los valores del usuario
	UserSesion.Email = email
	//La clave del servidor
	UserSesion.ServerKey = KservidorHash
	//Obtenemos las Claves RSA
	privateKey, publicKey := utils.GeneratePrivatePublicKeys()
	//La clave publica la almaceno directamente
	UserSesion.PublicKey = utils.PublicKeyToPem(&publicKey)
	//La clave privada primero la paso a []byte
	privateKeyPem := utils.PrivateKeyToPem(privateKey)
	//La cifro con AES
	privateKeyCipher := utils.CifrarAES(Kaes, IV, privateKeyPem)
	//La almaceno cifrada
	UserSesion.PrivateKey = privateKeyCipher
	//Guardo Kaes para la sesion del usuario
	UserSesion.Kaes = Kaes

	//Enviamos los datos al servidor
	userIDstring := RegisterServer(UserSesion)
	if userIDstring == "" {
		UserSesion = User{}
		return false
	} else {
		id, _ := primitive.ObjectIDFromHex(userIDstring)
		UserSesion.ID = id
		return true
	}
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
	client := &http.Client{}

	//Realizo la peticion
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	//Respuesta
	//En caso de fallo del registro del usuario en el servidor
	if resp.StatusCode == 400 {
		fmt.Println("El usuario no se ha registrado")
		return ""
	} else {
		//Si todo fue correcto en el servidor devuelvo el id del usuario creado
		var resultado string
		json.NewDecoder(resp.Body).Decode(&resultado)
		return resultado
	}
}

//Login del usuario
func LogIn(email string, password string) bool {
	user_pass := []byte(email + password)
	Kservidor, _, Kaes := HashUser(user_pass)
	UserSesion.Email = email
	UserSesion = LogInServer()
	err := bcrypt.CompareHashAndPassword(UserSesion.ServerKey, Kservidor)
	if err != nil {
		UserSesion = User{}
		return false
	} else {
		fmt.Println("La contraseña coincide")
		//Guardo Kaes para la sesion del usuario
		UserSesion.Kaes = Kaes
		return true
	}
}

//Login del usuario en el servidor
func LogInServer() User {
	userJSON, err := json.Marshal(UserSesion)
	if err != nil {
		fmt.Println(err)
	}

	url := config.URLbase + "login"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(userJSON))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	//Realizo la peticion
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		fmt.Println("El usuario no esta registrado")
		return User{}
	} else {
		//Si todo fue correcto en el servidor devuelvo el id del usuario creado
		var resultado User
		json.NewDecoder(resp.Body).Decode(&resultado)
		return resultado
	}
}

//Cerrar la sesion
func LogOut() {
	UserSesion = User{}
	DatosUsuario = []DatosUser{}
}

//Recupero todos los usuarios
func GetUsers() []User {
	var usersResponse []User
	resp, err := http.Get(config.URLbase + "users")
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

func GetUserByEmail(userEmail string) User {
	var usersResponse User
	resp, err := http.Get(config.URLbase + "users/" + userEmail)
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
	DatosUsuario = []DatosUser{}
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
			for j := 0; j < len(relations[i].Lists); j++ {
				list := GetUserList(relations[i].Lists[j].ListID, relations[i].Lists[j].ListKey, privateKey)
				lists = append(lists, list)
			}

			//Listas VIEJO BORRAR TRAS REVISAR
			// var ListsIDs []string
			// var listKeysCipher [][]byte
			// for j := 0; j < len(relations[i].Lists); j++ {
			// 	ListsIDs = append(ListsIDs, relations[i].Lists[j].ListID)
			// 	listKeysCipher = append(listKeysCipher, relations[i].Lists[j].ListKey)
			// }
			// var lists []ListCipher
			// if len(ListsIDs) > 0 {
			// 	lists = GetListsByIDs(ListsIDs)
			// }

			// //Desciframos las listas del proyecto
			// var listsDescifradas []List
			// for index := 0; index < len(lists); index++ {
			// 	//Descifro la Key de la lista antes de descifrar la lista en si
			// 	listKey := utils.DescifrarRSA(privateKey, listKeysCipher[index])
			// 	listsDescifradas = append(listsDescifradas, DescifrarLista(lists[index], listKey))
			// }

			datos := DatosUser{
				Proyecto: proyectoDescifrado,
				Listas:   lists,
			}
			DatosUsuario = append(DatosUsuario, datos)
		}
	}
}

//Elimina a un usuario del sistema borrandolo de todo proyectos, listas y tareas
func DeleteUser(userEmail string) {
	DatosUsuario = []DatosUser{}
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
					//Quito al usuario del array Users de la Lista
					DeleteUserList(DatosUsuario[i].Listas[j].ID, userEmail)

					//Traer las tareas de la lista y quitar al usuario de todas ellas

				}
			}
		}
	}
	//Borro las relaciones
	DeleteUserRelation(userEmail)
	//Borro al usuario
	DeleteUserByEmail(userEmail)
}

//Elimino al usuario del sistema
func DeleteUserByEmail(userEmail string) bool {
	url := config.URLbase + "users/" + userEmail
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
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

/*
func PostUser() {
	//Creo el usuario que voy a mandar al servidor
	user := UserCipher{Cipherdata: "enviandoDesdeCliente"}

	//Convertimos el user de tipo objeto GO a un JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}

	//Preparo la peticion POST
	url := config.URLbase + "user"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(userJSON))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	//Realizo la peticion
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	//Respuesta
	var response string
	json.NewDecoder(resp.Body).Decode(&response)
	fmt.Println(response)
}
*/
