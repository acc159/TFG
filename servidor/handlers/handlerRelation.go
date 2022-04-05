package handlers

import (
	"encoding/json"
	"net/http"
	"servidor/models"

	"github.com/gorilla/mux"
)

func GetRelations(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	email := params["email"]
	relations := models.GetRelationsbyUserEmail(email)
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
	params := mux.Vars(r)
	proyectID := params["proyectID"]
	userEmail := params["userEmail"]
	resultado := models.DeleteRelation(userEmail, proyectID)
	if resultado {
		respuesta := "La relacion fue borrada"
		json.NewEncoder(w).Encode(respuesta)
	} else {
		w.WriteHeader(400)
		respuesta := "No se pudo borrar la relacion"
		json.NewEncoder(w).Encode(respuesta)
	}
}

func DeleteListRelation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	proyectID := params["proyectID"]
	userEmail := params["userEmail"]
	listID := params["listID"]
	resultado := models.DeleteListRelation(userEmail, proyectID, listID)
	if resultado {
		respuesta := "La lista en la relacion fue borrada"
		json.NewEncoder(w).Encode(respuesta)
	} else {
		w.WriteHeader(400)
		respuesta := "No se pudo borrar la lista en la relacion"
		json.NewEncoder(w).Encode(respuesta)
	}
}

func AddListToRelation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	proyectID := params["proyectID"]
	userEmail := params["userEmail"]

	var list models.RelationLists
	json.NewDecoder(r.Body).Decode(&list)

	resultado := models.AddListToRelation(userEmail, proyectID, list)
	if !resultado {
		w.WriteHeader(400)
		w.Write([]byte("La lista no fue añadida a la relacion"))
	} else {
		w.Write([]byte("Lista añadida a la relacion"))
	}
}

func UpdateRelationList(w http.ResponseWriter, r *http.Request) {
	var relation models.Relation
	json.NewDecoder(r.Body).Decode(&relation)
	resultado := models.UpdateRelationList(relation)
	if resultado {
		respuesta := "La lista en la relacion fue actualizada"
		json.NewEncoder(w).Encode(respuesta)
	} else {
		w.WriteHeader(400)
		respuesta := "No se actualizar"
		json.NewEncoder(w).Encode(respuesta)
	}
}

func DeleteRelationByUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userEmail := params["userEmail"]
	resultado := models.DeleteRelationByUser(userEmail)
	if !resultado {
		w.WriteHeader(400)
		w.Write([]byte("Las relaciones no se han borrado"))
	} else {
		w.Write([]byte("Todas las relaciones del usuario han sido borradas"))
	}
}

func GetRelationsByUserProyect(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	email := params["email"]
	proyectID := params["proyectID"]
	relation := models.GetRelationsByUserProyect(email, proyectID)
	if relation.ID.Hex() == "000000000000000000000000" {
		w.WriteHeader(400)
		w.Write([]byte("Ninguna relacion para dicho usuario y proyecto"))
		return
	} else {
		json.NewEncoder(w).Encode(relation)
	}
}

func GetRelationsByUserList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	email := params["userEmail"]
	listID := params["listID"]
	relationList := models.GetRelationsByUserList(email, listID)
	if relationList.ListID.Hex() == "000000000000000000000000" {
		w.WriteHeader(400)
		w.Write([]byte("No existe dicha lista en la relacion para ese proyecto y usuario"))
		return
	} else {
		json.NewEncoder(w).Encode(relationList)
	}
}
