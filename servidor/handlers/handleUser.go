package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"servidor/models"

	"github.com/gorilla/mux"
)

//Aqui tengo los controladores que responden a las peticiones a las diferentes rutas

func Signup(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	usuarios := models.GetUsers()
	//Si no existe ningun usuario devuelve un error indicandolo
	if len(usuarios) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("No existen usuarios"))
		return
	}
	json.NewEncoder(w).Encode(usuarios)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	//Obtengo el id de los parametros de la petición
	params := mux.Vars(r)
	id := params["id"]
	usuario := models.GetUser(id)
	if usuario.Empty() {
		w.WriteHeader(400)
		w.Write([]byte("No existe el usuario"))
	} else {
		json.NewEncoder(w).Encode(usuario)
	}

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var usuario models.User
	json.NewDecoder(r.Body).Decode(&usuario)
	usuarioID := models.CreateUser(usuario)

	if usuarioID == "" {
		w.WriteHeader(400)
		w.Write([]byte("No se creo el usuario"))
		return
	}

	json.NewEncoder(w).Encode(usuarioID)
}

type ModifiedStructure struct {
	Campo string `json:"campo"`
	Valor string `json:"valor"`
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	var estructuraNueva ModifiedStructure
	json.NewDecoder(r.Body).Decode(&estructuraNueva)

	//Obtengo el id de los parametros de la petición
	params := mux.Vars(r)
	id := params["id"]
	usuario := models.UpdateUser(id, estructuraNueva.Campo, estructuraNueva.Valor)
	if usuario.Empty() {
		w.WriteHeader(400)
		w.Write([]byte("No se actualizo el usuario"))
	} else {
		json.NewEncoder(w).Encode(usuario)
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	borrado := models.DeleteUser(id)
	if borrado == false {
		w.WriteHeader(400)
		w.Write([]byte("No se borro el usuario"))
	} else {
		w.Write([]byte("Usuario borrado"))
	}

}