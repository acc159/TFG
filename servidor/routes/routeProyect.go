package routes

import (
	"servidor/handlers"
	"servidor/middlewares"

	"github.com/gorilla/mux"
)

func Proyect(r *mux.Router) {
	middlewares.ValidateTokenMiddleware(r.ServeHTTP)

	//Recupero todos los proyectos que tengan los IDs pasados en el body de la peticion
	r.HandleFunc("/proyects/ids", handlers.GetProyectsByIDs).Methods("GET")
	//Crear un proyecto
	r.HandleFunc("/proyect", handlers.CreateProyect).Methods("POST")
	//Actualizar un proyecto
	r.HandleFunc("/proyects/{id}", handlers.UpdateProyect).Methods("PUT")
	//Borrar un proyecto
	r.HandleFunc("/proyects/{id}", handlers.DeleteProyect).Methods("DELETE")
	//Recupero un proyecto por su ID
	r.HandleFunc("/proyects/{id}", handlers.GetProyect).Methods("GET")
	//Recupero todos los usuarios del campo Users del proyecto dado
	r.HandleFunc("/proyect/users/{id}", handlers.GetUsersProyect).Methods("GET")
	//Añado un usuario al campo Users del proyecto
	r.HandleFunc("/proyect/users/{id}/{email}", handlers.AddUserProyect).Methods("POST")
	//Borro un usuario del campo Users del proyecto
	r.HandleFunc("/proyect/users/{id}/{email}", handlers.DeleteUserProyect).Methods("DELETE")

	//Sin usar
	//Recupero todos los proyectos
	r.HandleFunc("/proyects", handlers.GetProyects).Methods("GET")

}
