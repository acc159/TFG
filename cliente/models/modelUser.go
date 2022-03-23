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
	ID         primitive.ObjectID `json:"_id,omitempty"`
	Nombre     string             `json:"nombre"`
	Apellidos  string             `json:"apellidos"`
	Email      string             `json:"email"`
	PublicKey  string             `json:"public_key"`
	PrivateKey string             `json:"private_key"`
}

type UserCipher struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Cipherdata string             `bson:"cipherdata"`
}

func Register(usuario string, password string, user User) {
	//Obtenemos Kaes , IV, Kservidor
	user_pass := usuario + password
	SAL := 12 //Con esto generara un segmento aleatorio cuanto mayor sea el numero mas robusto que usara como SAL
	hash, err := bcrypt.GenerateFromPassword([]byte(user_pass), SAL)
	if err != nil {
		fmt.Println("Error al hashear")
	}
	fmt.Println(hash)

	err = bcrypt.CompareHashAndPassword(hash, []byte("sadfdsf"))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("La contraseña coincide")
	}

	//Enviamos los datos al servidor

	//El usuario accede

}

func LogIn() {

}

func GetUsers() {
	resp, err := http.Get(config.URLbase + "users")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	//Compruebo si no hay ningun usuario
	if resp.StatusCode == 404 {
		fmt.Println("Ningun usuario encontrado")
	} else {
		var responseObject []UserCipher
		json.NewDecoder(resp.Body).Decode(&responseObject)
		fmt.Println(responseObject)
	}
}

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
