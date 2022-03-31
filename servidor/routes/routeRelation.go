package routes

import (
	"servidor/handlers"

	"github.com/gorilla/mux"
)

func Relation(r *mux.Router) {

	//Recupero las relaciones para el email pasado como parametro
	r.HandleFunc("/relations/{email}", handlers.GetRelations).Methods("GET")
	//Creo una nueva relacion
	r.HandleFunc("/relation", handlers.CreateRelation).Methods("POST")
	//Borro una relacion
	r.HandleFunc("/relations", handlers.DeleteRelation).Methods("DELETE")
	//Borro la lista de una relacion
	r.HandleFunc("/relations/list", handlers.DeleteListRelation).Methods("DELETE")
	//Añado una lista nueva al campo lists de una relacion
	r.HandleFunc("/relations/list/{proyectID}/{userEmail}", handlers.AddListToRelation).Methods("PUT")
	//Borro todas las relaciones para un usuario dado
	r.HandleFunc("/relations/user/{userEmail}", handlers.DeleteRelationByUser).Methods("DELETE")

	//Sin usar
	r.HandleFunc("/relations/list", handlers.UpdateRelationList).Methods("PUT")

}
