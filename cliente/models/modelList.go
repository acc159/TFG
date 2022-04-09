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
	ProyectID   string   `bson:"proyectID,omitempty"`
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
		CreateListRelations(listID, proyectIDstring, Krandom, list.Users)
		return listID
	}
}

func CreateListRelations(listID string, proyectID string, Krandom []byte, users []string) {
	for i := 0; i < len(users); i++ {
		publicKeyUser := GetPublicKey(users[i])
		KrandomCipher := utils.EncryptKeyWithPublicKey(publicKeyUser, Krandom)
		AddListToRelation(proyectID, listID, users[i], KrandomCipher)
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
		var listCipher ListCipher
		json.NewDecoder(resp.Body).Decode(&listCipher)
		//Descifro
		return listCipher
	}
}

//Elimino al usuario del array Users de la lista
func DeleteUserList(listID string, userEmail string) bool {
	//Elimino al usuario como usuario asignado de todas las tareas de la lista
	tasks := GetTasksByList(listID)
	for i := 0; i < len(tasks); i++ {
		tasks[i].Users = utils.FindAndDelete(tasks[i].Users, userEmail)
		UpdateTask(listID, tasks[i])
	}

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
	//Recupero la clave publica que usare para cifrar la clave del proyecto
	publicKey := GetPublicKey(userEmail)
	//Recupero la clave de la lista descifrada
	listKey := GetListKey(listID)
	//Cifro la clave de la lista con la clave publica del usuario nuevo añadido
	newListKey := utils.CifrarRSA(publicKey, listKey)
	//Añado la lista a la relacion
	AddListToRelation(proyectID, listID, userEmail, newListKey)
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

//Devuelve una lista descifrada dado el id de una lista cifrada junto a la clave de la lista y la clave privada del usuario para descifrar dicha clave de la lista
func GetUserList(listID string, listKeyCipher []byte, privateKey *rsa.PrivateKey) List {
	listKey := utils.DescifrarRSA(privateKey, listKeyCipher)
	listCipher := GetList(listID)
	return DescifrarLista(listCipher, listKey)
}

func UpdateList(newList List) bool {
	relation := GetRelationUserProyect(UserSesion.Email, newList.ProyectID)
	var listKeyCipher []byte
	//Busco la Clave cifrada de la lista
	for i := 0; i < len(relation.Lists); i++ {
		if relation.Lists[i].ListID == newList.ID {
			listKeyCipher = relation.Lists[i].ListKey
		}
	}
	//Descifro la clave con la clave del usuario actual
	privateKey := GetPrivateKeyUser()
	listKey := utils.DescifrarRSA(privateKey, listKeyCipher)
	//Genero un nuevo IV
	_, IV := utils.GenerateKeyIV()
	//Cifro la nueva lista
	listCipher := CifrarLista(newList, listKey, IV)
	listCipher.ID, _ = primitive.ObjectIDFromHex(newList.ID)

	//Actualizo la lista en el servidor
	listJSON, err := json.Marshal(listCipher)
	if err != nil {
		fmt.Println(err)
	}
	url := config.URLbase + "lists/" + newList.ID
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(listJSON))
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
		fmt.Println("La lista no pudo ser actualizada")
		return false
	} else {
		return true
	}
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
	list.ProyectID = listCipher.ProyectID.Hex()
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

func GetListKey(stringListID string) []byte {
	listKeyCipher := GetRelationListByUser(UserSesion.Email, stringListID).ListKey
	//Descifro la clave
	privateKey := GetPrivateKeyUser()
	listKey := utils.DescifrarRSA(privateKey, listKeyCipher)
	return listKey
}
