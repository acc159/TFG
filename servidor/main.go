package main

import (
	"log"
	"net/http"
	"servidor/config"
	"servidor/middlewares"
	"servidor/routes"
	"servidor/utils"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
}

func main() {

	//Cargar variables de entorno
	LoadEnv()

	//Cargo la Autoridad Certificadora
	utils.LoadAC()

	utils.PrepareJWT()

	/*

		secret := os.Getenv("SECRET")
		if secret != "" {
			fmt.Println(secret)
		} else {
			fmt.Println("No asignada")
		}
	*/

	//Iniciar base de datos
	//config.InitMongoDB()

	//Base de datos
	config.ConnectDB()
	defer config.DisconectDB()

	//Creo el enrutador
	r := mux.NewRouter()

	//Middlewares
	r.Use(middlewares.MiddlewareLog)
	r.Use(middlewares.MiddlewareAddJsonHeader)
	r.Use(middlewares.ValidateTokenMiddleware)
	//Defino las rutas
	//Usuarios
	routes.User(r)
	//Proyectos
	routes.Proyect(r)
	//Tareas
	routes.Task(r)
	//Listas
	routes.List(r)
	//Relaciones
	routes.Relation(r)

	//Lanzo el servidor
	//http.ListenAndServe("localhost:8080", r)
	http.ListenAndServeTLS("localhost:443", "certs/server.crt", "certs/server.key", r)

	// sigChan := make(chan os.Signal)
	// signal.Notify(sigChan)
	// select {
	// case <-sigChan:
	// 	config.Hola()
	// }
}
