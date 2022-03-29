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
	var proyecto models.Proyect
	json.NewDecoder(r.Body).Decode(&proyecto)
	models.CreateProyect(proyecto)

	w.Write([]byte("Creado proyecto"))
}

func UpdateProyect(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var proyecto models.Proyect
	json.NewDecoder(r.Body).Decode(&proyecto)
	resultado := models.UpdateProyect(proyecto, id)
	if !resultado {
		w.WriteHeader(400)
		w.Write([]byte("No se pudo actualizar el proyecto"))
	} else {
		w.Write([]byte("Proyecto actualizado"))
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

	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	resultado := models.AddUserProyect(id, user.Email)
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

	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	resultado := models.DeleteUserProyect(id, user.Email)
	if !resultado {
		w.WriteHeader(400)
		w.Write([]byte("No se pudo borrar el usuario del proyecto"))
	} else {
		w.Write([]byte("Usuario borrado"))
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
