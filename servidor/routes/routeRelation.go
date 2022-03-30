package routes

import (
	"servidor/handlers"

	"github.com/gorilla/mux"
)

func Relation(r *mux.Router) {
	r.HandleFunc("/relations/{email}", handlers.GetRelations).Methods("GET")
	r.HandleFunc("/relation", handlers.CreateRelation).Methods("POST")
	r.HandleFunc("/relations", handlers.DeleteRelation).Methods("DELETE")
	//r.HandleFunc("/relations/list", handlers.DeleteRelationList).Methods("DELETE")
	r.HandleFunc("/relations/list", handlers.UpdateRelationList).Methods("PUT")

	r.HandleFunc("/relations/list/{proyectID}/{userEmail}", handlers.AddListToRelation).Methods("PUT")

}
