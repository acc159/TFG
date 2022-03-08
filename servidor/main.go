package main

import (
	"log"
	"net/http"
)

func main() {
	println("Hola mundo desde el servidor")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
