package routes

import (
	"servidor/handlers"

	"github.com/gorilla/mux"
)

func User(r *mux.Router) {

	//Aqui defino las rutas para Usuario
	r.HandleFunc("/signup", handlers.Signup).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	r.HandleFunc("/user", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/books/{id}", handlers.UpdateUser).Methods("PUT")
	r.HandleFunc("/books/{id}", handlers.DeleteUser).Methods("DELETE")
}
