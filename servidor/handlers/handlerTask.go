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
