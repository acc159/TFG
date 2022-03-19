package routes

import (
	"servidor/handlers"

	"github.com/gorilla/mux"
)

func Proyect(r *mux.Router) {
	r.HandleFunc("/proyects", handlers.GetProyects).Methods("GET")
	r.HandleFunc("/proyect", handlers.CreateProyect).Methods("POST")
	r.HandleFunc("/proyect/:id", handlers.UpdateProyect).Methods("PUT")
	r.HandleFunc("/proyect/:id", handlers.DeleteProyect).Methods("DELETE")
}
