package routes

import (
	"servidor/handlers"

	"github.com/gorilla/mux"
)

func Relation(r *mux.Router) {
	r.HandleFunc("/relations/{id}", handlers.GetRelations).Methods("GET")
	r.HandleFunc("/relation", handlers.CreateRelation).Methods("POST")
	r.HandleFunc("/relations", handlers.DeleteRelation).Methods("DELETE")
}
