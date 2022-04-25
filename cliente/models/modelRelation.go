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

type Relation struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserEmail  string             `bson:"userEmail,omitempty"`
	ProyectID  primitive.ObjectID `bson:"proyectID,omitempty"`
	ProyectKey []byte             `bson:"proyectKey,omitempty"`
	Lists      []RelationLists    `bson:"lists,omitempty"`
	Rol        string             `bson:"rol"`
}

type RelationLists struct {
	ListID  string `bson:"listID,omitempty"`
	ListKey []byte `bson:"listKey,omitempty"`
}

var RelationsLocal []Relation

//Recupero las relaciones para un usuario dado su email
func GetProyectsListsByUser(userEmail string) []Relation {

	url := config.URLbase + "relations/" + userEmail
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
	var relations []Relation
	if resp.StatusCode == 400 {
		return relations
	} else {
		json.NewDecoder(resp.Body).Decode(&relations)
		return relations
	}
}

//Creo una relacion sin listas
func CreateRelation(userEmail string, proyectStringID string, proyectKey []byte) bool {
	//userID, _ := primitive.ObjectIDFromHex(userStringID)
	proyectID, _ := primitive.ObjectIDFromHex(proyectStringID)
	//Creo la relacion a enviar
	var relation Relation

	if userEmail == UserSesion.Email {
		relation = Relation{
			UserEmail:  userEmail,
			ProyectID:  proyectID,
			ProyectKey: proyectKey,
			Rol:        "Admin",
		}
	} else {
		relation = Relation{
			UserEmail:  userEmail,
			ProyectID:  proyectID,
			ProyectKey: proyectKey,
			Rol:        "User",
		}
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
	req = AddTokenHeader(req)
	client := utils.GetClientHTTPS()
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
	req = AddTokenHeader(req)
	client := utils.GetClientHTTPS()
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
func DeleteRelationList(proyectStringID string, listStringID string, userEmail string) bool {
	url := config.URLbase + "relations/list/" + userEmail + "/" + proyectStringID + "/" + listStringID
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
		fmt.Println("La lista en la relacion no pudo ser borrada")
		return false
	} else {
		return true
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
	req = AddTokenHeader(req)
	client := utils.GetClientHTTPS()
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

//Elimino la relacion con el proyectID y el userEmail dados
func DeleteUserProyectRelation(userEmail string, proyectID string) bool {
	url := config.URLbase + "relations/" + proyectID + "/" + userEmail
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
		fmt.Println("La relacion usuario-proyecto no pudo ser eliminada")
		return false
	} else {
		return true
	}
}

//Devolver una relacion para un usuario y proyecto dado
func GetRelationUserProyect(userEmail string, proyectID string) (Relation, bool) {
	url := config.URLbase + "relations/" + userEmail + "/" + proyectID
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
	var relation Relation

	switch resp.StatusCode {
	case 400:
		return relation, false
	case 401:
		fmt.Println("Token Expirado")
		return relation, true
	default:
		json.NewDecoder(resp.Body).Decode(&relation)
		return relation, false
	}
}

//Devolver una relacion de tipo lista para un usuario y lista dado
func GetRelationListByUser(userEmail string, listID string) RelationLists {
	url := config.URLbase + "relations/list/" + userEmail + "/" + listID
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
	var relationList RelationLists
	if resp.StatusCode == 400 {
		fmt.Println("No existe la lista para dicho usuario")
		return relationList
	} else {
		json.NewDecoder(resp.Body).Decode(&relationList)
		return relationList
	}
}

func CheckChanges() bool {
	relations := GetProyectsListsByUser(UserSesion.Email)
	if len(relations) != len(RelationsLocal) {
		return true
	} else {
		for i := 0; i < len(relations); i++ {
			if len(relations[i].Lists) != len(RelationsLocal[i].Lists) {
				return true
			}
		}
	}
	return false
}
