package routes

import (
	"servidor/handlers"

	"github.com/gorilla/mux"
)

//Rutas Usuario
func User(r *mux.Router) {

	//Registro
	r.HandleFunc("/signup", handlers.Signup).Methods("POST")
	//Login
	r.HandleFunc("/login", handlers.Login).Methods("POST")
	//Borrar Usuario
	r.HandleFunc("/users/{userEmail}", handlers.DeleteUser).Methods("DELETE")

	//SIN USAR
	//Recuperar Usuarios
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	//Recuperar Usuario por ID
	r.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	//Crear Usuario
	r.HandleFunc("/user", handlers.CreateUser).Methods("POST")
	//Actualizar Usuario
	r.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")

}
