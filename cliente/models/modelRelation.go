package models

import (
	"bytes"
	"cliente/config"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Relation struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserEmail  string             `bson:"userEmail,omitempty"`
	ProyectID  primitive.ObjectID `bson:"proyectID,omitempty"`
	ProyectKey string             `bson:"proyectKey,omitempty"`
	Lists      []RelationLists    `bson:"lists,omitempty"`
}

type RelationLists struct {
	ListID  primitive.ObjectID `bson:"listID,omitempty"`
	ListKey string             `bson:"listKey,omitempty"`
}

func GetProyectsListsByUser(userEmail string) []Relation {

	resp, err := http.Get(config.URLbase + "relations/" + userEmail)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var relations []Relation
	if resp.StatusCode == 400 {
		fmt.Println("Ningun proyecto ni lista para dicho usuario")
		return relations
	} else {
		json.NewDecoder(resp.Body).Decode(&relations)
		return relations
	}
}

func CreateRelation(userEmail string, proyectStringID string, listStringID string) bool {
	//userID, _ := primitive.ObjectIDFromHex(userStringID)
	proyectID, _ := primitive.ObjectIDFromHex(proyectStringID)

	//Compruebo si es una relacion de proyecto sin lista o con lista

	//Creo la relacion a enviar
	relation := Relation{
		UserEmail:  userEmail,
		ProyectID:  proyectID,
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
		return false
	} else {
		var responseObject string
		json.NewDecoder(resp.Body).Decode(&responseObject)
		fmt.Println(responseObject)
		return true
	}

}

func AddListToRelation(proyectID string, listID string) {
	//Recupero la relacion

	//Actualizo
}

/*
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
	defer resp.Body.Close()git a

	if resp.StatusCode == 400 {
		fmt.Println("La relacion no pudo ser borrada")
	} else {
		var resultado string
		json.NewDecoder(resp.Body).Decode(&resultado)
		fmt.Println(resultado)
	}
}
*/
