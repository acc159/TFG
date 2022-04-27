package handlers

import (
	"encoding/json"
	"net/http"
	"servidor/models"
	"servidor/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateList(w http.ResponseWriter, r *http.Request) {
	w = utils.SetRefreshToken(w, r)
	var list models.List
	json.NewDecoder(r.Body).Decode(&list)
	resultado := models.CreateList(list)

	if resultado != "" {
		json.NewEncoder(w).Encode(resultado)
		return
	} else {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(resultado)
	}
}

func DeleteList(w http.ResponseWriter, r *http.Request) {
	w = utils.SetRefreshToken(w, r)
	params := mux.Vars(r)
	id := params["id"]

	borrada := models.DeleteList(id)
	if borrada {
		w.Write([]byte("Lista borrada"))
		return
	} else {
		w.WriteHeader(400)
		w.Write([]byte("No pudo ser borrada la lista"))
	}
}

func UpdateList(w http.ResponseWriter, r *http.Request) {
	w = utils.SetRefreshToken(w, r)
	params := mux.Vars(r)
	id := params["id"]
	var list models.List
	json.NewDecoder(r.Body).Decode(&list)
	modificado := models.UpdateList(list, id)

	switch modificado {
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

func GetListsByIDs(w http.ResponseWriter, r *http.Request) {
	w = utils.SetRefreshToken(w, r)
	var IDs []string
	json.NewDecoder(r.Body).Decode(&IDs)
	if len(IDs) > 0 {
		lists := models.GetListsByIDs(IDs)
		if len(lists) > 0 {
			json.NewEncoder(w).Encode(lists)
		} else {
			w.WriteHeader(400)
			w.Write([]byte("No hay ninguna lista para esos ids"))
		}
	} else {
		w.WriteHeader(400)
		w.Write([]byte("No has enviado ningun id"))
	}
}

func GetUsersList(w http.ResponseWriter, r *http.Request) {
	w = utils.SetRefreshToken(w, r)
	params := mux.Vars(r)
	id := params["id"]
	resultado := models.GetUsersList(id)
	if len(resultado) == 0 {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(resultado)
	} else {
		json.NewEncoder(w).Encode(resultado)
	}
}

func AddUserList(w http.ResponseWriter, r *http.Request) {
	w = utils.SetRefreshToken(w, r)
	params := mux.Vars(r)
	id := params["id"]
	email := params["email"]
	resultado := models.AddUserList(id, email)
	if !resultado {
		w.WriteHeader(400)
		w.Write([]byte("No se pudo añadir el usuario a la lista"))
	} else {
		w.Write([]byte("Usuario añadido a la lista"))
	}
}

func DeleteUserList(w http.ResponseWriter, r *http.Request) {
	w = utils.SetRefreshToken(w, r)
	params := mux.Vars(r)
	id := params["id"]
	email := params["email"]
	resultado := models.DeleteUserList(id, email)
	if !resultado {
		w.WriteHeader(400)
		w.Write([]byte("No se pudo borrar el usuario de la lista"))
	} else {
		w.Write([]byte("Usuario borrado de la lista"))
	}
}

func GetList(w http.ResponseWriter, r *http.Request) {
	w = utils.SetRefreshToken(w, r)
	params := mux.Vars(r)
	id := params["id"]
	list := models.GetList(id)
	var mongoId interface{}
	mongoId = list.ID
	stringID := mongoId.(primitive.ObjectID).Hex()
	if stringID == "000000000000000000000000" {
		w.WriteHeader(400)
		w.Write([]byte("No existe la lista"))
	} else {
		json.NewEncoder(w).Encode(list)
	}
}
