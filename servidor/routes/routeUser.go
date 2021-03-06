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
	//Recuperar Usuario por userEmail
	r.HandleFunc("/users/{email}", handlers.GetUser).Methods("GET")

	//SIN USAR
	//Recuperar Usuarios
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")

	//Crear Usuario
	r.HandleFunc("/user", handlers.CreateUser).Methods("POST")

	r.HandleFunc("/user/refresh", handlers.GetRefreshToken).Methods("GET")

	//Actualizar Usuario
	r.HandleFunc("/users/{email}", handlers.UpdateUser).Methods("PUT")

	//Certificado
	r.HandleFunc("/user/certificate/{userEmail}", handlers.GetUserCertificate).Methods("GET")
}
