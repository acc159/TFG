package routes

import (
	"servidor/handlers"

	"github.com/gorilla/mux"
)

func Proyect(r *mux.Router) {
	r.HandleFunc("/proyects", handlers.GetProyects).Methods("GET")
	r.HandleFunc("/proyects/{id}", handlers.GetProyect).Methods("GET")
	r.HandleFunc("/proyect", handlers.CreateProyect).Methods("POST")
	r.HandleFunc("/proyects/{id}", handlers.UpdateProyect).Methods("PUT")
	r.HandleFunc("/proyects/{id}", handlers.DeleteProyect).Methods("DELETE")
}
