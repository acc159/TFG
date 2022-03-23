package handlers

import (
	"encoding/json"
	"net/http"
	"servidor/models"

	"github.com/gorilla/mux"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	json.NewDecoder(r.Body).Decode(&task)
	respuesta := models.CreateTask(task)

	if respuesta == "000000000000000000000000" {
		w.WriteHeader(400)
		respuesta := "No se creo la tarea"
		json.NewEncoder(w).Encode(respuesta)
	} else {
		json.NewEncoder(w).Encode(respuesta)
	}
}

func GetTasksByList(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	listID := params["listID"]

	tasks := models.GetTasksByList(listID)
	json.NewEncoder(w).Encode(tasks)

}

func UpdateTask(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]
	var newTask models.Task
	json.NewDecoder(r.Body).Decode(&newTask)
	modificado := models.UpdateTask(newTask, id)
	if modificado {
		respuesta := "Tarea actualizada"
		json.NewEncoder(w).Encode(respuesta)
	} else {
		w.WriteHeader(400)
		respuesta := "No se pudo modificar la tarea"
		json.NewEncoder(w).Encode(respuesta)
	}
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	borrada := models.DeleteTask(id)
	if borrada {
		respuesta := "Se borro la tarea"
		json.NewEncoder(w).Encode(respuesta)
		return
	} else {
		w.WriteHeader(400)
		respuesta := "No se pudo borrar la tarea"
		json.NewEncoder(w).Encode(respuesta)
	}

}
