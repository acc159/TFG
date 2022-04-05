package main

import (
	"bytes"
	"cliente/config"
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
	Datos []models.DataUser
}

type DataList struct {
	ProyectID    string
	UsersProyect []string
}

type DataTask struct {
	Tasks  []models.Task
	ListID string
}

type DataConfig struct {
	Proyect models.Proyect
	List    models.List
	Emails  []string
	User    string
}

func InitUI() {
	//Inicializo
	UI, _ = lorca.New("", "", 800, 700)
	//Cargo la primera vista
	ChangeView(config.PreView + "login.html")
}

//Cambia las vistas entre login y register
func ChangeView(nombreVista string) {
	content, err := ioutil.ReadFile(nombreVista)
	if err != nil {
		log.Fatal(err)
	}
	myFileContents := string(content)
	loadableContents := "data:text/html," + url.PathEscape(myFileContents)
	UI.Load(loadableContents)
}

//Carga la pagina de Home con los proyectos y listas del usuario actual
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

//Cargo la vista de añadir una lista
func ChangeViewAddList(nombreVista string, proyectID string, usersProyect []string) {
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

func ChangeViewTasks(nombreVista string, tasks []models.Task, listID string) {
	tmpl, err := template.ParseFiles(nombreVista)
	if err != nil {
		log.Fatal(err)
	}
	dataStruct := DataTask{
		Tasks:  tasks,
		ListID: listID,
	}
	buff := bytes.Buffer{}
	tmpl.Execute(&buff, dataStruct)
	loadableContents := "data:text/html," + url.PathEscape(buff.String())
	UI.Load(loadableContents)
}

func ChangeViewConfig(nombreVista string, proyect models.Proyect, list models.List, emails []string, userEmail string) {
	tmpl, err := template.ParseFiles(nombreVista)
	if err != nil {
		log.Fatal(err)
	}
	dataStruct := DataConfig{
		Proyect: proyect,
		List:    list,
		Emails:  emails,
		User:    userEmail,
	}
	buff := bytes.Buffer{}
	tmpl.Execute(&buff, dataStruct)
	loadableContents := "data:text/html," + url.PathEscape(buff.String())
	UI.Load(loadableContents)
}
