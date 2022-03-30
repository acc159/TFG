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

	dataStruct := Data{
		User:  models.UserSesion,
		Datos: models.DatosUsuario,
	}
	buff := bytes.Buffer{}

	tmpl.Execute(&buff, dataStruct)

	loadableContents := "data:text/html," + url.PathEscape(buff.String())
	UI.Load(loadableContents)
}

type DataList struct {
	ProyectID    string
	UsersProyect []string
}

func ChangeViewAddList(nombreVista string, proyectID string, usersProyect []string) {

	//Recupero los usuarios del proyecto

	tmpl, err := template.ParseFiles(nombreVista)
	if err != nil {
		log.Fatal(err)
	}

	dataStruct := DataList{
		ProyectID:    proyectID,
		UsersProyect: usersProyect,
	}
	buff := bytes.Buffer{}

	tmpl.Execute(&buff, dataStruct)

	loadableContents := "data:text/html," + url.PathEscape(buff.String())
	UI.Load(loadableContents)
}
