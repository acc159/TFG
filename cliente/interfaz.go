package main

import (
	"bytes"
	"cliente/models"

	"html/template"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/zserge/lorca"
)

var UI lorca.UI

type Data struct {
	User  models.User
	Datos []models.DatosUser
}

func InitUI() {
	//Inicializo
	UI, _ = lorca.New("", "", 800, 700)
	//Cargo la primera vista
	ChangeView("www/login.html")
}

func ChangeView(nombreVista string) {
	content, err := ioutil.ReadFile(nombreVista)
	if err != nil {
		log.Fatal(err)
	}
	myFileContents := string(content)

	loadableContents := "data:text/html," + url.PathEscape(myFileContents)
	UI.Load(loadableContents)
}

func ChangeViewWithValues(nombreVista string) {

	tmpl, err := template.ParseFiles(nombreVista)
	if err != nil {
		log.Fatal(err)
	}
	// var users = []string{"<script>alert('a')</script>", "Canada", "Japan"}

	// proyect1 := models.ProyectCipher{
	// 	Cipherdata: "A",
	// 	Users:      users,
	// }
	// proyect2 := models.ProyectCipher{
	// 	Cipherdata: "dsafdsfasdfasdfsdfds",
	// }

	// var proyects []models.DatosUser
	// proyects = append(proyects, proyect1)
	// proyects = append(proyects, proyect2)

	dataStruct := Data{
		User:  models.UserSesion,
		Datos: models.DatosUsuario,
	}
	data := dataStruct
	buff := bytes.Buffer{}

	tmpl.Execute(&buff, data)

	loadableContents := "data:text/html," + url.PathEscape(buff.String())
	UI.Load(loadableContents)
}
