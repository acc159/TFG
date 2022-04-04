package models

import (
	"bytes"
	"cliente/config"
	"cliente/utils"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type List struct {
	ID          string   `bson:"_id,omitempty"`
	Name        string   `bson:"name,omitempty"`
	Description string   `bson:"description,omitempty"`
	Users       []string `bson:"users,omitempty"`
}

type ListCipher struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Cipherdata []byte             `bson:"cipherdata,omitempty"`
	Users      []string           `bson:"users,omitempty"`
	ProyectID  primitive.ObjectID `bson:"proyectID,omitempty"`
}

//Creo una lista con el proyectID correspondiente
func CreateList(list List, proyectIDstring string) string {

	//1.Generamos la clave aleatoria que se utilizara en el cifrado AES
	Krandom, IVrandom := utils.GenerateKeyIV()

	//Ciframos la lista
	listCipher := CifrarLista(list, Krandom, IVrandom)
	proyectID, _ := primitive.ObjectIDFromHex(proyectIDstring)
	listCipher.ProyectID = proyectID
	//Enviamos la lista cifrada
	listJSON, err := json.Marshal(listCipher)
	if err != nil {
		fmt.Println(err)
	}
	url := config.URLbase + "list"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(listJSON))
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
		return ""
	} else {
		var listID string
		json.NewDecoder(resp.Body).Decode(&listID)

		//Añado la lista a la relacion de cada usuario miembro de la lista
		for i := 0; i < len(list.Users); i++ {
			//Recupero para cada usuario su Public Key y la uso para cifrar la Krandom
			publicKeyUser := GetUserByEmail(list.Users[i]).PublicKey
			publicKey := utils.PemToPublicKey(publicKeyUser)
			KrandomCipher := utils.CifrarRSA(publicKey, Krandom)
			AddListToRelation(proyectIDstring, listID, list.Users[i], KrandomCipher)
		}

		return listID
	}
}

//Recupero todas las listas cuyos ids sean los que contiene el array de ids pasado como parametro
func GetListsByIDs(stringsIDs []string) []ListCipher {
	relationJSON, err := json.Marshal(stringsIDs)
	if err != nil {
		fmt.Println(err)
	}
	url := config.URLbase + "lists/ids"
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(relationJSON))
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
	//Compruebo si no hay ningun usuario
	var responseObject []ListCipher
	if resp.StatusCode == 404 {
		fmt.Println("No hay ninguna lista")
	} else {
		json.NewDecoder(resp.Body).Decode(&responseObject)
	}
	return responseObject
}

//Eliminar una lista
func DeleteList(listsID string) bool {
	//Eliminar la lista
	url := config.URLbase + "lists/" + listsID
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
		return false
	} else {
		return true
	}

}

//Recupero los usuarios de una lista
func GetUsersList(listID string) []string {
	resp, err := http.Get(config.URLbase + "list/users/" + listID)
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

//Recupero una lista dado su ID
func GetList(listID string) ListCipher {
	resp, err := http.Get(config.URLbase + "list/" + listID)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	//Compruebo si no hay ningun usuario
	if resp.StatusCode == 404 {
		fmt.Println("Lista no encontrada")
		return ListCipher{}
	} else {
		var list ListCipher
		json.NewDecoder(resp.Body).Decode(&list)
		//Descifro
		return list
	}
}

//Elimino al usuario del array Users de la lista
func DeleteUserList(listID string, userEmail string) bool {
	url := config.URLbase + "list/users/" + listID + "/" + userEmail
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
		fmt.Println("El usuario no pudo ser eliminado de la lista")
		return false
	} else {
		return true
	}
}

//Añadir usuario a una lista
func AddUserList(userEmail string, proyectID string, listID string) bool {
	//Recupero la relacion para el usuario actual para poder descifrar la clave del proyecto que necesitare cifrar para añadir a la relacion del nuevo usuario
	relationUser := GetRelationUserProyect(UserSesion.Email, proyectID)
	//Obtengo la clave para la lista determinada
	var listKey []byte
	for i := 0; i < len(relationUser.Lists); i++ {
		if relationUser.Lists[i].ListID == listID {
			listKey = relationUser.Lists[i].ListKey
		}
	}

	//Cifro para el nuevo usuario

	//Añado la lista a la relacion
	AddListToRelation(proyectID, listID, userEmail, listKey)

	//Añado el usuario al array Users de la lista
	url := config.URLbase + "list/users/" + listID + "/" + userEmail
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
		fmt.Println("El usuario no pudo ser añadido a la lista")
		return false
	} else {
		return true
	}
}

func GetUserList(listID string, listKeyCipher []byte, privateKey *rsa.PrivateKey) List {
	listKey := utils.DescifrarRSA(privateKey, listKeyCipher)
	listCipher := GetList(listID)
	return DescifrarLista(listCipher, listKey)
}

//Cifrado y Descifrado

// func DescifrarLista(listCipher ListCipher) List {
// 	list := List{
// 		ID:          listCipher.ID.Hex(),
// 		Name:        "Nombre de la lista",
// 		Description: "Descripcion de la lista",
// 		Users:       listCipher.Users,
// 	}
// 	return list
// }

// func CifrarLista(listCipher List) ListCipher {
// 	return ListCipher{
// 		Cipherdata: "DASFSDFASDFSDF",
// 		Users:      listCipher.Users,
// 	}
// }

func DescifrarLista(listCipher ListCipher, key []byte) List {

	descifradoBytes := utils.DescifrarAES(key, listCipher.Cipherdata)
	list := BytesToList(descifradoBytes)
	list.ID = listCipher.ID.Hex()
	list.Users = listCipher.Users
	return list
}

func CifrarLista(list List, key []byte, IV []byte) ListCipher {
	//Paso el proyecto a []byte
	listBytes := ListToBytes(list)
	//Cifro
	listCipherBytes := utils.CifrarAES(key, IV, listBytes)

	listCipher := ListCipher{
		Cipherdata: listCipherBytes,
		Users:      list.Users,
	}
	return listCipher
}

func ListToBytes(list List) []byte {
	listBytes, _ := json.Marshal(list)
	return listBytes
}

func BytesToList(datos []byte) List {
	var list List
	err := json.Unmarshal(datos, &list)
	if err != nil {
		fmt.Println("error:", err)
	}
	return list
}
