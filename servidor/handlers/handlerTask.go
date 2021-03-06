package handlers

import (
	"encoding/json"
	"net/http"
	"servidor/models"
	"servidor/utils"

	"github.com/gorilla/mux"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	w = utils.SetRefreshToken(w, r)
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

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	w = utils.SetRefreshToken(w, r)
	params := mux.Vars(r)
	taskID := params["taskID"]
	task := models.GetTaskByID(taskID)
	if task.ID.Hex() == "000000000000000000000000" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(task)
	} else {
		json.NewEncoder(w).Encode(task)
	}
}

func GetTasksByList(w http.ResponseWriter, r *http.Request) {
	w = utils.SetRefreshToken(w, r)
	params := mux.Vars(r)
	listID := params["listID"]
	tasks := models.GetTasksByList(listID)
	if len(tasks) == 0 {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(tasks)
	} else {
		json.NewEncoder(w).Encode(tasks)
	}
}

func DeleteTasksByList(w http.ResponseWriter, r *http.Request) {
	w = utils.SetRefreshToken(w, r)
	params := mux.Vars(r)
	listID := params["listID"]
	result := models.DeleteTasksByListID(listID)
	if !result {
		w.WriteHeader(400)
		respuesta := "No se borraron"
		json.NewEncoder(w).Encode(respuesta)
	} else {
		respuesta := "Se borraron"
		json.NewEncoder(w).Encode(respuesta)
	}
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	w = utils.SetRefreshToken(w, r)
	params := mux.Vars(r)
	id := params["id"]
	var newTask models.Task
	json.NewDecoder(r.Body).Decode(&newTask)
	modificado := models.UpdateTask(newTask, id)
	if modificado == "OK" {
		respuesta := "Tarea actualizada"
		json.NewEncoder(w).Encode(respuesta)
	} else if modificado == "Error" {
		w.WriteHeader(400)
		respuesta := "No se pudo modificar la tarea"
		json.NewEncoder(w).Encode(respuesta)
	} else {
		w.WriteHeader(470)
		respuesta := "Tarea ya actualizada"
		json.NewEncoder(w).Encode(respuesta)
	}
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w = utils.SetRefreshToken(w, r)
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
