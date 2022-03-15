package main

import (
	"fmt"
	"net/http"
	"servidor/config"
	"servidor/middlewares"
	"servidor/routes"

	"github.com/gorilla/mux"
)

func main() {

	//Base de datos
	config.ConnectDB()
	defer config.DisconectDB()

	//Creo el enrutador
	r := mux.NewRouter()

	//Middlewares
	r.Use(middlewares.MiddlewareLog)
	r.Use(middlewares.MiddlewareAddJsonHeader)

	//Defino las rutas
	//Usuarios
	routes.User(r)

	//Lanzo el servidor
	http.ListenAndServe("localhost:8080", r)
	fmt.Println("Servidor corriendo en el puerto 8080")

}
