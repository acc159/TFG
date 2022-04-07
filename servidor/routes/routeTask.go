package routes

import (
	"servidor/handlers"

	"github.com/gorilla/mux"
)

func Task(r *mux.Router) {
	//Creo una tarea
	r.HandleFunc("/task", handlers.CreateTask).Methods("POST")
	//Recupero una tarea por su ID
	r.HandleFunc("/tasks/{taskID}", handlers.GetTaskByID).Methods("GET")
	//Recupero las tareas pertenecientes a una lista dado el ID de la lista
	r.HandleFunc("/tasks/list/{listID}", handlers.GetTasksByList).Methods("GET")
	//Actualizo una tarea
	r.HandleFunc("/tasks/{id}", handlers.UpdateTask).Methods("PUT")
	//Borro una tarea
	r.HandleFunc("/tasks/{id}", handlers.DeleteTask).Methods("DELETE")
}
