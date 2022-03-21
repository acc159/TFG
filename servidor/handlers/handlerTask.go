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
	models.CreateTask(task)

	w.Write([]byte("Tarea creada"))
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
		w.Write([]byte("Tarea actualizada"))
		return
	} else {
		w.WriteHeader(400)
		w.Write([]byte("No se pudo modificar la tarea"))
	}
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	borrada := models.DeleteTask(id)
	if borrada {
		w.Write([]byte("Tarea borrada"))
		return
	} else {
		w.WriteHeader(400)
		w.Write([]byte("No pudo ser borrada la tarea"))
	}

}
