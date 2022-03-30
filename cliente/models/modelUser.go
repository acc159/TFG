package models

import (
	"bytes"
	"cliente/config"
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
	PublicKey  string             `bson:"public_key"`
	PrivateKey string             `bson:"private_key"`
	Rol        string             `bson:"rol"`
}

var UserSesion User

type DatosUser struct {
	Proyecto Proyect
	Listas   []List
}

var DatosUsuario []DatosUser

func Register(email string, password string) bool {

	//Pruebas
	/*
		prueba := []byte(email + password)
		sum := sha256.Sum256(prueba)
		sumString := string(sum[:])
		pass := hex.EncodeToString(sum[:]) //String

		fmt.Println(sum)
		fmt.Println(sum[0:16])
		fmt.Println(sum[16:32])
		fmt.Println(sumString)
		fmt.Println(pass)

		hash2, _ := bcrypt.GenerateFromPassword([]byte(pass), 12)
		fmt.Println(hash2)
	*/
	//

	user_pass := email + password
	SAL := 12 //Con esto generara un segmento aleatorio cuanto mayor sea el numero mas robusto que usara como SAL
	hash, err := bcrypt.GenerateFromPassword([]byte(user_pass), SAL)
	if err != nil {
		fmt.Println("Error al hashear")
	}

	UserSesion.Email = email
	UserSesion = User{
		Email: email,
	}

	//La clave del servidor
	UserSesion.ServerKey = hash
	//Obtenemos las Claves RSA
	UserSesion.PrivateKey = "ClavePrivada"
	UserSesion.PublicKey = "ClavePublica"

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
	/*
		err = bcrypt.CompareHashAndPassword(hash, []byte(user_pass))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("La contraseña coincide")
		}
	*/

	//El usuario accede
}

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

func LogIn(email string, password string) bool {

	user_pass := email + password
	UserSesion.Email = email

	UserSesion = LogInServer()

	err := bcrypt.CompareHashAndPassword(UserSesion.ServerKey, []byte(user_pass))
	if err != nil {
		fmt.Println(err)
		UserSesion = User{}
		return false
	} else {
		fmt.Println("La contraseña coincide")
		return true
	}

}

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

func LogOut() {
	UserSesion = User{}
	DatosUsuario = []DatosUser{}
}

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

//Obtengo los proyectos y sus listas asociadas para el usuario
func GetUserProyectsLists() {

	//Limpio los datos del usuario
	DatosUsuario = []DatosUser{}

	//Recupero las relaciones
	relations := GetProyectsListsByUser(UserSesion.Email)
	if len(relations) > 0 {
		//Proyectos
		for i := 0; i < len(relations); i++ {
			proyecto := GetProyect(relations[i].ProyectID.Hex())

			//Listas
			var ListsIDs []string
			for j := 0; j < len(relations[i].Lists); j++ {
				ListsIDs = append(ListsIDs, relations[i].Lists[j].ListID.Hex())
			}
			var lists []ListCipher
			if len(ListsIDs) > 0 {
				lists = GetListsByIDs(ListsIDs)
			}

			//Desciframos el proyecto
			proyectoDescifrado := DescifrarProyecto(proyecto)

			//Desciframos las listas del proyecto
			var listsDescifradas []List
			for index := 0; index < len(lists); index++ {
				listsDescifradas = append(listsDescifradas, DescifrarLista(lists[index]))
			}

			//Desciframos las listas

			datos := DatosUser{
				Proyecto: proyectoDescifrado,
				Listas:   listsDescifradas,
			}
			DatosUsuario = append(DatosUsuario, datos)
		}
	}
	//return DatosUsuario
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
