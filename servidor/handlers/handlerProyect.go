package handlers

import (
	"encoding/json"
	"net/http"
	"servidor/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JsonCustom struct {
	StringIDs []string `json:"stringsIDs"`
}

func GetProyects(w http.ResponseWriter, r *http.Request) {
	proyectos := models.GetProyects()
	if len(proyectos) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("No existen proyectos"))
	} else {
		json.NewEncoder(w).Encode(proyectos)
	}
}

func GetProyect(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	proyecto := models.GetProyect(id)
	var mongoId interface{}
	mongoId = proyecto.ID
	stringID := mongoId.(primitive.ObjectID).Hex()
	if stringID == "000000000000000000000000" {
		w.WriteHeader(400)
		w.Write([]byte("No existe el proyecto"))
	} else {
		json.NewEncoder(w).Encode(proyecto)
	}
}

func CreateProyect(w http.ResponseWriter, r *http.Request) {
	var proyecto *models.Proyect
	json.NewDecoder(r.Body).Decode(&proyecto)
	models.CreateProyect(proyecto)
	if proyecto.ID.Hex() == "" {
		w.WriteHeader(400)
		respuesta := "Proyecto no creado"
		json.NewEncoder(w).Encode(respuesta)
	} else {
		json.NewEncoder(w).Encode(proyecto.ID.Hex())
	}
}

func UpdateProyect(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var proyecto models.Proyect
	json.NewDecoder(r.Body).Decode(&proyecto)
	resultado := models.UpdateProyect(proyecto, id)
	switch resultado {
	case "Error":
		w.WriteHeader(400)
		w.Write([]byte("No se pudo modificar la lista"))
	case "Ya modificada":
		w.WriteHeader(470)
		respuesta := "Lista ya actualizada"
		json.NewEncoder(w).Encode(respuesta)
	default:
		w.Write([]byte("Lista actualizada"))
		return
	}
}

func DeleteProyect(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	resultado := models.DeleteProyect(id)
	if !resultado {
		w.WriteHeader(400)
		w.Write([]byte("No se pudo borrar el proyecto"))
	} else {
		w.Write([]byte("Proyecto borrado"))
	}
}

func AddUserProyect(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	email := params["email"]
	resultado := models.AddUserProyect(id, email)
	if !resultado {
		w.WriteHeader(400)
		w.Write([]byte("No se pudo añadir el usuario al proyecto"))
	} else {
		w.Write([]byte("Usuario añadido al proyecto"))
	}
}

func DeleteUserProyect(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	email := params["email"]
	resultado := models.DeleteUserProyect(id, email)
	if !resultado {
		w.WriteHeader(400)
		w.Write([]byte("No se pudo borrar el usuario del proyecto"))
	} else {
		w.Write([]byte("Usuario borrado del proyecto"))
	}
}

func GetProyectsByIDs(w http.ResponseWriter, r *http.Request) {
	var IDs JsonCustom
	json.NewDecoder(r.Body).Decode(&IDs)
	if len(IDs.StringIDs) > 0 {
		proyects := models.GetProyectsByIDs(IDs.StringIDs)
		if len(proyects) > 0 {
			json.NewEncoder(w).Encode(proyects)
		} else {
			w.WriteHeader(400)
			w.Write([]byte("No hay ningun proyecto para esos ids"))
		}
	} else {
		w.WriteHeader(400)
		w.Write([]byte("No has enviado ningun id"))
	}

}

func GetUsersProyect(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	resultado := models.GetUsersProyect(id)
	if len(resultado) == 0 {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(resultado)
	} else {
		json.NewEncoder(w).Encode(resultado)
	}
}
