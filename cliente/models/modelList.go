package models

import (
	"bytes"
	"cliente/config"
	"cliente/utils"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/hex"
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
	Check       string   `bson:"check"`
}

type ListCipher struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Cipherdata  []byte             `bson:"cipherdata,omitempty"`
	Users       []string           `bson:"users,omitempty"`
	ProyectID   primitive.ObjectID `bson:"proyectID,omitempty"`
	Check       string             `bson:"check"`
	UpdateCheck string             `bson:"updateCheck"`
}

//Creo una lista con el proyectID correspondiente
func CreateList(list List, proyectIDstring string) (bool, bool) {
	list.Users = append(list.Users, UserSesion.Email)
	//1.Generamos la clave aleatoria que se utilizara en el cifrado AES
	Krandom, IVrandom := utils.GenerateKeyIV()
	//Ciframos la lista
	listCipher := CifrarLista(list, Krandom, IVrandom)
	proyectID, _ := primitive.ObjectIDFromHex(proyectIDstring)
	listCipher.ProyectID = proyectID

	h := sha1.New()
	h.Write(listCipher.Cipherdata)
	listCipher.Check = hex.EncodeToString(h.Sum(nil))

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
	req = AddTokenHeader(req)
	client := utils.GetClientHTTPS()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 400:
		fmt.Println("El proyecto no pudo ser creado")
		return false, false
	case 401:
		fmt.Println("Token Expirado")
		return false, true
	default:
		var listID string
		json.NewDecoder(resp.Body).Decode(&listID)
		//Añado la lista a la relacion de cada usuario miembro de la lista
		CreateListRelations(listID, proyectIDstring, Krandom, list.Users)
		return true, false
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
	req = AddTokenHeader(req)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := utils.GetClientHTTPS()
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
func DeleteList(listsID string) (bool, bool) {
	//Eliminar la lista
	url := config.URLbase + "lists/" + listsID
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
	switch resp.StatusCode {
	case 400:
		fmt.Println("La lista no pudo ser borrada")
		return false, false
	case 401:
		fmt.Println("Token Expirado")
		return false, true
	default:
		return true, false
	}
}

//Recupero los usuarios de una lista
func GetUsersList(listID string) []string {
	url := config.URLbase + "list/users/" + listID
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

func ExistList(listID string) bool {
	list := GetList(listID)
	return !list.ID.IsZero()
}

//Recupero una lista dado su ID
func GetList(listID string) ListCipher {

	url := config.URLbase + "list/" + listID
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
	req = AddTokenHeader(req)
	client := utils.GetClientHTTPS()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		fmt.Println("El usuario no pudo ser eliminado de la lista")
		return false
	} else {
		for i := 0; i < len(DatosUsuario); i++ {
			for j := 0; j < len(DatosUsuario[i].Listas); j++ {
				if DatosUsuario[i].Listas[j].ID == listID {
					DatosUsuario[i].Listas[j].Users = utils.FindAndDelete(DatosUsuario[i].Listas[j].Users, userEmail)
				}
			}
		}
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
	req = AddTokenHeader(req)
	client := utils.GetClientHTTPS()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		fmt.Println("El usuario no pudo ser añadido a la lista")
		return false
	} else {
		//Lo añado al usuario en local
		for i := 0; i < len(DatosUsuario); i++ {
			for j := 0; j < len(DatosUsuario[i].Listas); j++ {
				if DatosUsuario[i].Listas[j].ID == listID {
					DatosUsuario[i].Listas[j].Users = append(DatosUsuario[i].Listas[j].Users, userEmail)
				}
			}
		}
		return true
	}
}

//Devuelve una lista descifrada dado el id de una lista cifrada junto a la clave de la lista y la clave privada del usuario para descifrar dicha clave de la lista
func GetUserList(listID string, listKeyCipher []byte, privateKey *rsa.PrivateKey) List {
	listKey := utils.DescifrarRSA(privateKey, listKeyCipher)
	listCipher := GetList(listID)
	return DescifrarLista(listCipher, listKey)
}

func UpdateList(newList List) (string, string) {
	relation, _ := GetRelationUserProyect(UserSesion.Email, newList.ProyectID)
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

	//En updateCheck pongo el hash de los datos anteriores
	listCipher.UpdateCheck = newList.Check
	h := sha1.New()
	h.Write(listCipher.Cipherdata)
	listCipher.Check = hex.EncodeToString(h.Sum(nil))

	//Actualizo la lista en el servidor
	listJSON, err := json.Marshal(listCipher)
	if err != nil {
		fmt.Println(err)
	}
	url := config.URLbase + "lists/" + newList.ID
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(listJSON))
	req = AddTokenHeader(req)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := utils.GetClientHTTPS()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	// switch resp.StatusCode {
	// case 400:
	// 	fmt.Println("La tarea no pudo ser borrada")
	// 	return "Error"
	// case 470:
	// 	fmt.Println("Token Expirado")
	// 	return "Ya actualizada"
	// default:
	// 	return "OK"
	// }

	switch resp.StatusCode {
	case 400:
		fmt.Println("La lista no pudo ser actualizada")
		return "Error", "false"
	case 401:
		fmt.Println("Token Expirado")
		return "", "true"
	case 470:
		return "Ya actualizada", "false"
	default:
		return "OK", "false"
	}
}

func DescifrarLista(listCipher ListCipher, key []byte) List {

	descifradoBytes := utils.DescifrarAES(key, listCipher.Cipherdata)
	list := BytesToList(descifradoBytes)
	list.ID = listCipher.ID.Hex()
	list.Users = listCipher.Users
	list.ProyectID = listCipher.ProyectID.Hex()
	list.Check = listCipher.Check
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
