package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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
	//Chequeo que el usuario admin esta registrado ya
	admin := models.GetUser("admin")
	if admin.Email != "" || user.Email == "admin" {
		//Compruebo que no exista ya
		existeUser := models.GetUser(user.Email)
		if existeUser.Email == "" {
			user.Status = "Activo"
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
	} else {
		respuesta := "Admin no encontrado"
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
		if usuario.Status != "Activo" {
			json.NewEncoder(w).Encode(models.User{})
		} else {
			//Genero el token para el usuario
			jwtToken := utils.GenerateJWT(usuario.Email)
			w.Header().Set("token", jwtToken)
			json.NewEncoder(w).Encode(usuario)
		}
	}
}

//Por revisar y modificar
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w = utils.SetRefreshToken(w, r)
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	params := mux.Vars(r)
	email := params["email"]
	resultado := models.UpdateUser(email, user.Status)
	if !resultado {
		w.WriteHeader(400)
		w.Write([]byte("No se actualizo el usuario"))
	} else {
		w.Write([]byte("Se actualizo el usuario"))
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w = utils.SetRefreshToken(w, r)
	params := mux.Vars(r)
	email := params["userEmail"]
	borrado := models.DeleteUser(email)
	if !borrado {
		w.WriteHeader(400)
		w.Write([]byte("No se borro el usuario"))
	} else {
		userCertFilename := "certs/users/" + email + "_cert.pem"
		e := os.Remove(userCertFilename)
		if e != nil {
			log.Fatal(e)
		}
		w.Write([]byte("Usuario borrado"))
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w = utils.SetRefreshToken(w, r)
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
	w = utils.SetRefreshToken(w, r)
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
	w = utils.SetRefreshToken(w, r)
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
	w = utils.SetRefreshToken(w, r)
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

func GetRefreshToken(w http.ResponseWriter, r *http.Request) {
	userToken := r.Header.Values("UserToken")[0]
	if userToken == "" {
		w.WriteHeader(400)
		w.Write([]byte("Error de identificación"))
	}
	refreshToken := utils.GenerateJWT(userToken)
	w.Header().Set("refreshToken", refreshToken)
	json.NewEncoder(w).Encode("Nuevo Token")
}
