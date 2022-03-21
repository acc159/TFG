package routes

import (
	"servidor/handlers"

	"github.com/gorilla/mux"
)

func List(r *mux.Router) {
	r.HandleFunc("/list", handlers.CreateList).Methods("POST")
	r.HandleFunc("/lists/{id}", handlers.UpdateList).Methods("PUT")
	r.HandleFunc("/lists/{id}", handlers.DeleteList).Methods("DELETE")
}
