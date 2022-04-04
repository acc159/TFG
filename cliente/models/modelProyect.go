package models

import (
	"bytes"
	"cliente/config"
	"cliente/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Proyect struct {
	ID          string   `bson:"_id,omitempty"`
	Name        string   `bson:"name,omitempty"`
	Description string   `bson:"description,omitempty"`
	Users       []string `bson:"users,omitempty"`
}

type ProyectCipher struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Cipherdata []byte             `bson:"cipherdata,omitempty"`
	Users      []string           `bson:"users,omitempty"`
	Lists      []string           `bson:"lists,omitempty"`
}

//Recupero un proyecto dado su ID
func GetProyect(proyectID string) ProyectCipher {
	resp, err := http.Get(config.URLbase + "proyects/" + proyectID)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	//Compruebo si no hay ningun usuario
	if resp.StatusCode == 404 {
		fmt.Println("Proyecto no encontrado")
		return ProyectCipher{}
	} else {
		var proyect ProyectCipher
		json.NewDecoder(resp.Body).Decode(&proyect)
		//Descifro
		return proyect
	}
}

//Creo un proyecto
func CreateProyect(newProyect Proyect) bool {
	//Añadimos el email del usuario que esta creando el proyecto
	newProyect.Users = append(newProyect.Users, UserSesion.Email)

	//Ciframos el proyecto
	//1.Generamos la clave aleatoria que se utilizara en el cifrado AES
	Krandom, IVrandom := utils.GenerateKeyIV()
	proyectCipher := CifrarProyecto(newProyect, Krandom, IVrandom)

	//Enviamos el proyecto cifrado al servidor
	proyectJSON, err := json.Marshal(proyectCipher)
	if err != nil {
		fmt.Println(err)
	}
	url := config.URLbase + "proyect"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(proyectJSON))
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
		fmt.Println("El proyecto no pudo ser creado")
		return false
	}
	//Si el proyecto se crea con exito nos devuelve el ID del proyecto creado
	var proyectID string
	json.NewDecoder(resp.Body).Decode(&proyectID)
	//Creamos la relacion para el usuario que crea el proyecto y para cada uno de los usuarios del campo user
	for i := 0; i < len(newProyect.Users); i++ {
		//Recupero para cada usuario su Public Key y la uso para cifrar la Krandom
		publicKeyUser := GetUserByEmail(newProyect.Users[i]).PublicKey
		publicKey := utils.PemToPublicKey(publicKeyUser)
		KrandomCipher := utils.CifrarRSA(publicKey, Krandom)
		CreateRelation(newProyect.Users[i], proyectID, KrandomCipher)
	}
	return true
}

//Eliminar un proyecto
func DeleteProyect(proyectID string) bool {
	url := config.URLbase + "proyects/" + proyectID
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
		fmt.Println("El proyecto no pudo ser borrado")
		return false
	} else {
		return true
	}
}

//Recuperar los usuarios de un proyecto
func GetUsersProyect(proyectID string) []string {
	resp, err := http.Get(config.URLbase + "proyect/users/" + proyectID)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var responseObject []string
	if resp.StatusCode == 404 {
		return responseObject
	} else {
		json.NewDecoder(resp.Body).Decode(&responseObject)
		return responseObject
	}
}

//Elimino al usuario del array Users del proyecto
func DeleteUserProyect(proyectID string, userEmail string) bool {

	//Recupero la relacion para quitar tambien al usuario de las listas del proyecto donde este
	relation := GetRelationUserProyect(userEmail, proyectID)
	//Para cada lista que tiene el proyecto elimino al usuario de dicha lista
	for i := 0; i < len(relation.Lists); i++ {
		DeleteUserList(relation.Lists[i].ListID, userEmail)
	}

	url := config.URLbase + "proyect/users/" + proyectID + "/" + userEmail
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
		fmt.Println("El usuario no pudo ser eliminado del proyecto")
		return false
	} else {
		return true
	}
}

//Añadir un usuario a un proyecto
func AddUserProyect(proyectIDstring string, userEmail string) bool {

	//Recupero la relacion para el usuario actual para poder descifrar la clave del proyecto que necesitare cifrar para añadir a la relacion del nuevo usuario
	relationUser := GetRelationUserProyect(UserSesion.Email, proyectIDstring)

	//Desciframos y obtenemos una nueva
	proyectKey := relationUser.ProyectKey
	CreateRelation(userEmail, proyectIDstring, proyectKey)

	//Añado el usuario al array de usuarios del proyecto
	url := config.URLbase + "proyect/users/" + proyectIDstring + "/" + userEmail
	req, err := http.NewRequest("POST", url, nil)
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
		fmt.Println("El usuario no pudo ser añadido al proyecto")
		return false
	} else {
		return true
	}
}

//Cifrado y Descifrado

func DescifrarProyecto(proyectCipher ProyectCipher, key []byte) Proyect {

	descifradoBytes := utils.DescifrarAES(key, proyectCipher.Cipherdata)
	proyect := BytesToProyect(descifradoBytes)
	proyect.ID = proyectCipher.ID.Hex()
	proyect.Users = proyectCipher.Users
	return proyect
}

func CifrarProyecto(proyect Proyect, key []byte, IV []byte) ProyectCipher {
	//Paso el proyecto a []byte
	proyectBytes := ProyectToBytes(proyect)
	//Cifro
	proyectCipherBytes := utils.CifrarAES(key, IV, proyectBytes)

	proyectCipher := ProyectCipher{
		Cipherdata: proyectCipherBytes,
		Users:      proyect.Users,
	}
	return proyectCipher
}

func ProyectToBytes(proyect Proyect) []byte {
	proyectBytes, _ := json.Marshal(proyect)
	return proyectBytes
}

func BytesToProyect(datos []byte) Proyect {
	var proyect Proyect
	err := json.Unmarshal(datos, &proyect)
	if err != nil {
		fmt.Println("error:", err)
	}
	return proyect
}
