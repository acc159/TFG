package models

import (
	"bytes"
	"cliente/config"
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
	Cipherdata string             `bson:"cipherdata,omitempty"`
	Users      []string           `bson:"users,omitempty"`
	ProyectID  string             `bson:"proyectID,omitempty"`
}

func CreateList(list List) string {
	//Ciframos la lista
	listCipher := CifrarLista(list)
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
		return listID
	}

}

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

func DescifrarLista(listCipher ListCipher) List {
	list := List{
		ID:          listCipher.ID.Hex(),
		Name:        "Nombre de la lista",
		Description: "Descripcion de la lista",
		Users:       []string{"pepito@gmail.com", "juanito@gmail.com"},
	}
	return list
}

func CifrarLista(listCipher List) ListCipher {
	return ListCipher{
		Cipherdata: "DASFSDFASDFSDF",
		Users:      []string{"sadfds", "sadfasd"},
	}
}
