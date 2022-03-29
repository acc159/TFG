package handlers

import (
	"encoding/json"
	"net/http"
	"servidor/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetRelations(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	relations := models.GetRelationsbyUserID(id)
	if len(relations) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("Ninguna relacion para dicho usuario"))
		return
	} else {
		json.NewEncoder(w).Encode(relations)
	}
}

func CreateRelation(w http.ResponseWriter, r *http.Request) {
	var relation models.Relation
	json.NewDecoder(r.Body).Decode(&relation)
	resultado := models.CreateRelation(relation)
	if resultado {
		respuesta := "Relacion creada"
		json.NewEncoder(w).Encode(respuesta)
		return
	} else {
		w.WriteHeader(400)
		respuesta := "No se pudo crear la relacion"
		json.NewEncoder(w).Encode(respuesta)
	}

}

func DeleteRelation(w http.ResponseWriter, r *http.Request) {
	var relation models.Relation
	json.NewDecoder(r.Body).Decode(&relation)
	resultado := models.DeleteRelation(relation.UserID, relation.ProyectID)
	if resultado {
		respuesta := "La relacion fue borrada"
		json.NewEncoder(w).Encode(respuesta)
	} else {
		w.WriteHeader(400)
		respuesta := "No se pudo borrar la relacion"
		json.NewEncoder(w).Encode(respuesta)
	}
}

type RelationJSON struct {
	UserID    primitive.ObjectID `bson:"userID,omitempty"`
	ProyectID primitive.ObjectID `bson:"proyectID,omitempty"`
	ListID    primitive.ObjectID `bson:"listID,omitempty"`
}

func DeleteRelationList(w http.ResponseWriter, r *http.Request) {
	var relation RelationJSON
	json.NewDecoder(r.Body).Decode(&relation)
	resultado := models.DeleteRelationList(relation.UserID, relation.ProyectID, relation.ListID)
	if resultado {
		respuesta := "La lista en la relacion fue borrada"
		json.NewEncoder(w).Encode(respuesta)
	} else {
		w.WriteHeader(400)
		respuesta := "No se pudo borrar la lista en la relacion"
		json.NewEncoder(w).Encode(respuesta)
	}
}

func UpdateRelationList(w http.ResponseWriter, r *http.Request) {
	var relation models.Relation
	json.NewDecoder(r.Body).Decode(&relation)
	resultado := models.UpdateRelationList(relation)
	if resultado {
		respuesta := "La lista en la relacion fue borrada"
		json.NewEncoder(w).Encode(respuesta)
	} else {
		w.WriteHeader(400)
		respuesta := "No se pudo borrar la lista en la relacion"
		json.NewEncoder(w).Encode(respuesta)
	}
}
