package routes

import (
	"servidor/handlers"

	"github.com/gorilla/mux"
)

func List(r *mux.Router) {

	//Recupero todas las listas con los ids pasados en el body de la peticion
	r.HandleFunc("/lists/ids", handlers.GetListsByIDs).Methods("GET")
	//Crear una lista
	r.HandleFunc("/list", handlers.CreateList).Methods("POST")
	//Eliminar una lista
	r.HandleFunc("/lists/{id}", handlers.DeleteList).Methods("DELETE")
	//Recupero los usuarios de una lista
	r.HandleFunc("/list/users/{id}", handlers.GetUsersList).Methods("GET")
	//Recupero una lista por su ID
	r.HandleFunc("/list/{id}", handlers.GetList).Methods("GET")

	//Sin Usar
	//Actualizar una lista
	r.HandleFunc("/lists/{id}", handlers.UpdateList).Methods("PUT")

	//Añado un usuario al campo Users de la lista
	r.HandleFunc("/list/users/{id}/{email}", handlers.AddUserList).Methods("POST")
	//Borro un usuario del campo Users de la lista
	r.HandleFunc("/list/users/{id}/{email}", handlers.DeleteUserList).Methods("DELETE")
}
