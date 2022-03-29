package routes

import (
	"servidor/handlers"

	"github.com/gorilla/mux"
)

func Proyect(r *mux.Router) {

	r.HandleFunc("/proyects/ids", handlers.GetProyectsByIDs).Methods("GET")

	r.HandleFunc("/proyects", handlers.GetProyects).Methods("GET")
	r.HandleFunc("/proyects/{id}", handlers.GetProyect).Methods("GET")
	r.HandleFunc("/proyect", handlers.CreateProyect).Methods("POST")
	r.HandleFunc("/proyects/{id}", handlers.UpdateProyect).Methods("PUT")
	r.HandleFunc("/proyects/{id}", handlers.DeleteProyect).Methods("DELETE")

	r.HandleFunc("/proyects/users/{id}", handlers.AddUserProyect).Methods("POST")
	r.HandleFunc("/proyects/users/{id}", handlers.DeleteUserProyect).Methods("DELETE")

}
