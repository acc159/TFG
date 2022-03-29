package models

import (
	"cliente/config"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Relation struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"userID,omitempty"`
	ProyectID  primitive.ObjectID `bson:"proyectID,omitempty"`
	ProyectKey string             `bson:"proyectKey,omitempty"`
	Lists      []RelationLists    `bson:"lists,omitempty"`
}

type RelationLists struct {
	ListID  primitive.ObjectID `bson:"listID,omitempty"`
	ListKey string             `bson:"listKey,omitempty"`
}

var Relations []Relation

func GetProyectsListsByUser(userID string) {

	resp, err := http.Get(config.URLbase + "relations/" + userID)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		fmt.Println("Ningun proyecto ni lista para dicho usuario")
	} else {
		json.NewDecoder(resp.Body).Decode(&Relations)
	}
}

/*
func CreateRelation(userStringID string, proyectStringID string, listStringID string) {
	userID, _ := primitive.ObjectIDFromHex(userStringID)
	listID, _ := primitive.ObjectIDFromHex(proyectStringID)
	proyectID, _ := primitive.ObjectIDFromHex(listStringID)

	//Creo la relacion a enviar
	relation := Relation{
		UserID:     userID,
		ListID:     listID,
		ProyectID:  proyectID,
		ListKey:    "sdafdasf",
		ProyectKey: "asdfasdfasdfsdafsdf",
	}

	//Pasamos el tipo Relation a JSON
	relationJSON, err := json.Marshal(relation)
	if err != nil {
		fmt.Println(err)
	}

	//Peticion POST
	url := config.URLbase + "relation"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(relationJSON))
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
		fmt.Println("La relacion no pudo ser creada")
	} else {
		var responseObject string
		json.NewDecoder(resp.Body).Decode(&responseObject)
		fmt.Println(responseObject)
	}

}

func DeleteRelation(userStringID string, proyectStringID string, listStringID string) {

	userID, _ := primitive.ObjectIDFromHex(userStringID)
	proyectID, _ := primitive.ObjectIDFromHex(proyectStringID)
	listID, _ := primitive.ObjectIDFromHex(listStringID)

	relation := Relation{
		UserID:    userID,
		ListID:    listID,
		ProyectID: proyectID,
	}

	relationJSON, err := json.Marshal(relation)
	if err != nil {
		fmt.Println(err)
	}

	url := config.URLbase + "relations"
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(relationJSON))
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
		fmt.Println("La relacion no pudo ser borrada")
	} else {
		var resultado string
		json.NewDecoder(resp.Body).Decode(&resultado)
		fmt.Println(resultado)
	}
}
*/
