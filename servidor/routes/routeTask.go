package routes

import (
	"servidor/handlers"

	"github.com/gorilla/mux"
)

func Task(r *mux.Router) {
	r.HandleFunc("/task", handlers.CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{listID}", handlers.GetTasksByList).Methods("GET")
	r.HandleFunc("/tasks/{id}", handlers.UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", handlers.DeleteTask).Methods("DELETE")
}
