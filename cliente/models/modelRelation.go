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
	ListID  string `bson:"listID,omitempty"`
	ListKey string `bson:"listKey,omitempty"`
}

//Recupero las relaciones para un usuario dado su email
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

//Creo una relacion sin listas FALTA RELLENAR EL CAMPO PROYECT KEY
func CreateRelation(userEmail string, proyectStringID string) bool {
	//userID, _ := primitive.ObjectIDFromHex(userStringID)
	proyectID, _ := primitive.ObjectIDFromHex(proyectStringID)
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

//Añado una lista a la relacion dada
func AddListToRelation(proyectID string, listIDstring string, userEmail string) bool {
	//Creo la Relacion Lista
	// listID, _ := primitive.ObjectIDFromHex(listIDstring)
	relationList := RelationLists{
		ListID:  listIDstring,
		ListKey: "sadfsadfsadfsad",
	}
	relationListJSON, err := json.Marshal(relationList)
	if err != nil {
		fmt.Println(err)
	}
	url := config.URLbase + "relations/list/" + proyectID + "/" + userEmail
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(relationListJSON))
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

//Elimino una lista dado su id en una relacion
func DeleteRelationList(proyectStringID string, listStringID string, userEmail string) {
	datos := []string{userEmail, proyectStringID, listStringID}
	relationJSON, err := json.Marshal(datos)
	if err != nil {
		fmt.Println(err)
	}
	url := config.URLbase + "relations/list"
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
		fmt.Println("La lista en la relacion no pudo ser borrada")
	} else {
		var resultado string
		json.NewDecoder(resp.Body).Decode(&resultado)
		fmt.Println(resultado)
	}
}
