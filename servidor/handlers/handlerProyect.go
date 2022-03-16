package handlers

import (
	"encoding/json"
	"net/http"
	"servidor/models"
)

func GetProyects(w http.ResponseWriter, r *http.Request) {
	proyectos := models.GetProyects()
	if len(proyectos) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("No existen proyectos"))
	} else {
		json.NewEncoder(w).Encode(proyectos)
	}
}

func CreateProyect(w http.ResponseWriter, r *http.Request) {
	var proyecto models.Proyect
	json.NewDecoder(r.Body).Decode(&proyecto)
	models.CreateProyect(proyecto)

	w.Write([]byte("Proyectos"))
}
