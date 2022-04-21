package handlers

import (
	"encoding/json"
	"net/http"
	"servidor/models"
	"servidor/utils"

	"github.com/gorilla/mux"
)

type UserCertificate struct {
	Certificate   []byte `json:"certificate"`
	PublicKeyAC   []byte `json:"publicKeyAC"`
	PublicKeyUser []byte `json:"publicKeyUser"`
}

//Aqui tengo los controladores que responden a las peticiones a las diferentes rutas

//Registro del usuario
func Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	//Compruebo que no exista ya
	existeUser := models.GetUser(user.Email)
	if existeUser.Email == "" {
		resultado := models.SignUp(user)
		if resultado == "" {
			w.WriteHeader(400)
			respuesta := "No se registro el usuario"
			json.NewEncoder(w).Encode(respuesta)
		} else {
			utils.CreateUserCertificate(user.Email, user.PublicKey)
			json.NewEncoder(w).Encode(resultado)
		}
	} else {
		w.WriteHeader(409)
		respuesta := "Usuario Duplicado"
		json.NewEncoder(w).Encode(respuesta)
	}

}

//Login del Usuario
func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	usuario := models.Login(user)
	if usuario.Empty() {
		w.WriteHeader(400)
		respuesta := "No existe el usuario"
		json.NewEncoder(w).Encode(respuesta)
	} else {
		//Genero el token para el usuario
		jwtToken := utils.GenerateJWT(usuario.Email)
		usuario.Token = jwtToken
		w.Header().Set("token", jwtToken)
		json.NewEncoder(w).Encode(usuario)
	}
}

//Por revisar y modificar

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	params := mux.Vars(r)
	id := params["id"]
	resultado := models.UpdateUser(id, user)
	if !resultado {
		w.WriteHeader(400)
		w.Write([]byte("No se actualizo el usuario"))
	} else {
		w.Write([]byte("Se actualizo el usuario"))
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	email := params["userEmail"]
	borrado := models.DeleteUser(email)
	if !borrado {
		w.WriteHeader(400)
		w.Write([]byte("No se borro el usuario"))
	} else {
		w.Write([]byte("Usuario borrado"))
	}
}

//SIN USAR

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
	email := params["email"]
	usuario := models.GetUser(email)
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
		respuesta := "No se creo el usuario"
		json.NewEncoder(w).Encode(respuesta)
		return
	}

	json.NewEncoder(w).Encode(usuarioID)
}

func GetUserCertificate(w http.ResponseWriter, r *http.Request) {
	//Obtengo el id de los parametros de la petición
	params := mux.Vars(r)
	user := params["userEmail"]

	certificate := utils.GetUserCertificate(user)
	publicKeyAC := utils.GetACpublicKey()
	userPublicKey := models.GetUser(user).PublicKey
	userCertificate := UserCertificate{
		Certificate:   certificate,
		PublicKeyAC:   publicKeyAC,
		PublicKeyUser: userPublicKey,
	}
	if len(certificate) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("No existe el certificado"))
	} else {
		json.NewEncoder(w).Encode(userCertificate)
	}
}
