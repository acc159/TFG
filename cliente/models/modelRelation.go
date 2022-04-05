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
	ProyectKey []byte             `bson:"proyectKey,omitempty"`
	Lists      []RelationLists    `bson:"lists,omitempty"`
}

type RelationLists struct {
	ListID  string `bson:"listID,omitempty"`
	ListKey []byte `bson:"listKey,omitempty"`
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
func CreateRelation(userEmail string, proyectStringID string, proyectKey []byte) bool {
	//userID, _ := primitive.ObjectIDFromHex(userStringID)
	proyectID, _ := primitive.ObjectIDFromHex(proyectStringID)
	//Creo la relacion a enviar
	relation := Relation{
		UserEmail:  userEmail,
		ProyectID:  proyectID,
		ProyectKey: proyectKey,
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
func AddListToRelation(proyectID string, listIDstring string, userEmail string, listKey []byte) bool {
	//Creo la Relacion Lista
	// listID, _ := primitive.ObjectIDFromHex(listIDstring)
	relationList := RelationLists{
		ListID:  listIDstring,
		ListKey: []byte(listKey),
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
	url := config.URLbase + "relations/list/" + userEmail + "/" + proyectStringID + "/" + listStringID
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
		fmt.Println("La lista en la relacion no pudo ser borrada")
	} else {
		var resultado string
		json.NewDecoder(resp.Body).Decode(&resultado)
		fmt.Println(resultado)
	}
}

//Elimino las relaciones donde aparece el usuario
func DeleteUserRelation(userEmail string) bool {
	url := config.URLbase + "relations/user/" + userEmail
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
		fmt.Println("Las relaciones del usuario no pudieron ser eliminadas")
		return false
	} else {
		return true
	}
}

//Elimino la relacion con el proyectID y el userEmail dados
func DeleteUserProyectRelation(userEmail string, proyectID string) bool {
	url := config.URLbase + "relations/" + proyectID + "/" + userEmail
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
		fmt.Println("La relacion usuario-proyecto no pudo ser eliminada")
		return false
	} else {
		return true
	}
}

//Devolver una relacion para un usuario y proyecto dado
func GetRelationUserProyect(userEmail string, proyectID string) Relation {
	resp, err := http.Get(config.URLbase + "relations/" + userEmail + "/" + proyectID)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var relation Relation
	if resp.StatusCode == 400 {
		fmt.Println("Ningun proyecto ni lista para dicho usuario")
		return relation
	} else {
		json.NewDecoder(resp.Body).Decode(&relation)
		return relation
	}
}

//Devolver una relaciond de tipo lista para un usuario y lista dado
func GetRelationListByUser(userEmail string, listID string) RelationLists {
	resp, err := http.Get(config.URLbase + "relations/list/" + userEmail + "/" + listID)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var relationList RelationLists
	if resp.StatusCode == 400 {
		fmt.Println("No existe la lista para dicho usuario")
		return relationList
	} else {
		json.NewDecoder(resp.Body).Decode(&relationList)
		return relationList
	}
}
