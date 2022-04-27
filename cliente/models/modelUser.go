package models

import (
	"bytes"
	"cliente/config"
	"cliente/utils"
	"crypto"
	"crypto/aes"
	"crypto/rsa"
	"crypto/sha256"
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
	Kaes       []byte             `bson:"Kaes"`
	Token      string             `bson:"token,omitempty"`
	Status     string             `bson:"status,omitempty"`
}

type DataUser struct {
	Proyecto Proyect
	Listas   []List
}

type UserCertificate struct {
	Certificate   []byte `json:"certificate"`
	PublicKeyAC   []byte `json:"publicKeyAC"`
	PublicKeyUser []byte `json:"publicKeyUser"`
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
	UserSesion.Token = resp.Header.Get("refreshToken")
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
func LogIn(email string, password string) string {
	user_pass := []byte(email + password)
	Kservidor, _, Kaes := HashUser(user_pass)
	userLogin := User{
		Email:     email,
		ServerKey: Kservidor,
	}
	resultado := LogInServer(userLogin)

	if resultado == "OK" {
		UserSesion.Kaes = Kaes
	}
	return resultado

	// switch resultado {
	// case "Error":

	// case "Bloqueado":

	// default :
	// 	UserSesion.Kaes = Kaes
	// 	return true
	// }

	// if UserSesion.Email != "" {
	// 	UserSesion.Kaes = Kaes
	// 	return true
	// } else {
	// 	return false
	// }
}

//Login del usuario en el servidor
func LogInServer(userLogin User) string {
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
	UserSesion.Token = resp.Header.Get("refreshToken")
	if resp.StatusCode == 400 {
		return "Error"
	} else {
		//Si todo fue correcto en el servidor devuelvo el id del usuario creado
		var resultado User
		json.NewDecoder(resp.Body).Decode(&resultado)
		if resultado.Email == "" {
			return "Bloqueado"
		} else {
			//Asigno el token que genero el servidor
			token := resp.Header.Get("token")
			fmt.Println(token)
			UserSesion = resultado
			UserSesion.Token = token
			return "OK"
		}
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
	UserSesion.Token = resp.Header.Get("refreshToken")
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
func GetUserByEmail(userEmail string) (User, bool) {
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
	UserSesion.Token = resp.Header.Get("refreshToken")
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

//Obtengo las relaciones junto a los proyectos y sus listas asociadas para el usuario
func GetUserProyectsLists() {
	//Limpio los datos del usuario
	DatosUsuario = []DataUser{}
	//Obtengo mi clave privada
	privateKey := GetPrivateKeyUser()
	//Recupero las relaciones
	relations := GetProyectsListsByUser(UserSesion.Email)
	RelationsLocal = relations
	if len(relations) > 0 {
		//Proyectos
		for i := 0; i < len(relations); i++ {
			proyecto := GetProyect(relations[i].ProyectID.Hex())
			//Descifro la clave del proyecto con la clave privada del usuario
			proyectKey := utils.DescifrarRSA(privateKey, relations[i].ProyectKey)
			//Desciframos el proyecto
			proyectoDescifrado := DescifrarProyecto(proyecto, proyectKey)
			proyectoDescifrado.Rol = GetUserProyectRol(proyectoDescifrado)
			//Listas
			var lists []List
			//Por cada lista del proyecto la recupero descifrada usando mi clave privada para descifrar la clave de descifrado de la lista
			for j := 0; j < len(relations[i].Lists); j++ {
				list := GetUserList(relations[i].Lists[j].ListID, relations[i].Lists[j].ListKey, privateKey)
				list.Rol = GetUserListRol(list)
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

func GetUserProyectRol(proyect Proyect) string {
	for i := 0; i < len(proyect.Users); i++ {
		if proyect.Users[i].User == UserSesion.Email {
			return proyect.Users[i].Rol
		}
	}
	return "User"
}

func GetUserListRol(list List) string {
	for i := 0; i < len(list.Users); i++ {
		if list.Users[i].User == UserSesion.Email {
			return list.Users[i].Rol
		}
	}
	return "User"
}

// //Elimina a un usuario del sistema borrandolo de todo proyectos, listas y tareas
// func DeleteUser(userEmail string) bool {
// 	DatosUsuario = []DataUser{}
// 	//Recupero las relaciones junto a los proyectos y las listas Mejorable el pensar en llamar a una funcion que no descifre todo porque no lo necesitamos
// 	GetUserProyectsLists()
// 	for i := 0; i < len(DatosUsuario); i++ {
// 		DeleteUserProyect(DatosUsuario[i].Proyecto.ID, userEmail)
// 		for j := 0; j < len(DatosUsuario[i].Listas); j++ {
// 			//Quito al usuario del array Users de la Lista y de las tareas
// 			DeleteUserList(DatosUsuario[i].Listas[j].ID, userEmail)
// 		}
// 	}
// 	//Borro las relaciones
// 	DeleteUserRelation(userEmail)
// 	//Borro al usuario
// 	DeleteUserByEmail(userEmail)
// 	return true
// }

func UpdateStatus(userEmail string, status string) (bool, bool) {
	var user User
	user.Status = status
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}
	url := config.URLbase + "users/" + userEmail
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(userJSON))
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
	UserSesion.Token = resp.Header.Get("refreshToken")
	switch resp.StatusCode {
	case 400:
		return false, false
	case 401:
		fmt.Println("Token Expirado")
		return false, true
	default:
		return true, false
	}
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
	UserSesion.Token = resp.Header.Get("refreshToken")
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
		if users[i].Status != "Bloqueado" {
			usersEmails = append(usersEmails, users[i].Email)
		}
	}
	return usersEmails, tokenExpire
}

//Recupero la clave publica de un usuario
func GetPublicKey(userEmail string) (*rsa.PublicKey, bool) {
	user, tokenExpire := GetUserByEmail(userEmail)
	if tokenExpire {
		return nil, tokenExpire
	}
	publicKeyUserPem := user.PublicKey
	publicKey := utils.PemToPublicKey(publicKeyUserPem)
	return publicKey, false
}

func GetCertificateUser(userEmail string) UserCertificate {
	url := config.URLbase + "user/certificate/" + userEmail
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
	var userCertificate UserCertificate
	defer resp.Body.Close()
	UserSesion.Token = resp.Header.Get("refreshToken")
	if resp.StatusCode == 400 {
		fmt.Println("El certificado no se encontro")
		return userCertificate
	} else {
		json.NewDecoder(resp.Body).Decode(&userCertificate)
		return userCertificate
	}
}

func VerifyCertificateSign(userCertificate UserCertificate) bool {
	certificate := utils.PemToCertificate(userCertificate.Certificate)
	publicKeyAC := utils.PemToPublicKey(userCertificate.PublicKeyAC)

	//Verificar Firma del certificado con la C.Publica de la AC
	h := sha256.New()
	h.Write(certificate.RawTBSCertificate)
	hash_data := h.Sum(nil)

	//Si coincide la firma devuelvo true, en caso contrario habra un error y dara false
	err := rsa.VerifyPKCS1v15(publicKeyAC, crypto.SHA256, hash_data, certificate.Signature)
	return err == nil
}

func VerifyPublicKeyWithCertificate(userCertificate UserCertificate) bool {
	certificate := utils.PemToCertificate(userCertificate.Certificate)
	userPublicKey := utils.PemToPublicKey(userCertificate.PublicKeyUser)
	//Verificar si la clave publica del certificado coincide con la clave publica del usuario
	if certificate.PublicKey.(*rsa.PublicKey).N.Cmp(userPublicKey.N) == 0 && userPublicKey.E == certificate.PublicKey.(*rsa.PublicKey).E {
		return true
	} else {
		return false
	}
}

func RefreshTokenUser() bool {
	url := config.URLbase + "user/refresh"
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
		return false
	case 401:
		return false
	default:
		UserSesion.Token = resp.Header.Get("refreshToken")
		return true
	}
}

func AddTokenHeader(req *http.Request) *http.Request {
	req.Header.Add("Authorization", "Bearer "+UserSesion.Token)
	return req
}
