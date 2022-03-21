package handlers

import (
	"encoding/json"
	"net/http"
	"servidor/models"

	"github.com/gorilla/mux"
)

func CreateList(w http.ResponseWriter, r *http.Request) {
	var list models.List
	json.NewDecoder(r.Body).Decode(&list)
	resultado := models.CreateList(list)

	if resultado {
		w.Write([]byte("Lista Creada"))
		return
	} else {
		w.WriteHeader(400)
		w.Write([]byte("No se pudo crear la tarea"))
	}
}

func DeleteList(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	borrada := models.DeleteList(id)
	if borrada {
		w.Write([]byte("Tarea borrada"))
		return
	} else {
		w.WriteHeader(400)
		w.Write([]byte("No pudo ser borrada la tarea"))
	}

}

func UpdateList(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]
	var list models.List
	json.NewDecoder(r.Body).Decode(&list)
	modificado := models.UpdateList(list, id)
	if modificado {
		w.Write([]byte("Tarea actualizada"))
		return
	} else {
		w.WriteHeader(400)
		w.Write([]byte("No se pudo modificar la tarea"))
	}
}
