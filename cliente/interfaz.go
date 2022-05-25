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
	User   models.User
	Datos  []models.DataUser
	Emails []string
}

type DataList struct {
	User         models.User
	ProyectID    string
	ProyectName  string
	UsersProyect []models.UserRole
}

type DataTask struct {
	Tasks       []models.Task
	ListID      string
	ListName    string
	ListUsers   []models.UserRole
	User        string
	Pendientes  int
	Progreso    int
	Finalizadas int
	TypeTasks   string
}

type DataConfigProyect struct {
	Proyect models.Proyect
	Emails  []string
	User    string
}

type DataConfigList struct {
	List   models.List
	Emails []models.UserRole
	User   string
}

type DataConfigTask struct {
	Task                  models.Task
	List                  models.List
	User                  string
	HasSignReceivedByUser func(models.Task) bool
	HasSignCloseByUser    func(models.Task) bool
	CheckCreator          func(models.Task) bool
	IsUserAssigned        func(models.Task) bool
}

type AdminView struct {
	Users      []models.User
	UserActual string
}

func InitUI() {
	//Inicializo
	UI, _ = lorca.New("", "", 1250, 800, "--allow-insecure-localhost")
	//Cargo la primera vista
	ChangeView(config.PreView + "user/login.html")
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

//Carga la pagina de Home con los proyectos y listas del usuario actual y a la pantalla de añadir un proyecto
func ChangeViewWithValues(nombreVista string, emails []string) {
	tmpl, err := template.ParseFiles(nombreVista)
	if err != nil {
		log.Fatal(err)
	}
	dataStruct := Data{
		User:   models.UserSesion,
		Datos:  models.DatosUsuario,
		Emails: emails,
	}
	buff := bytes.Buffer{}
	tmpl.Execute(&buff, dataStruct)
	loadableContents := "data:text/html," + url.PathEscape(buff.String())
	UI.Load(loadableContents)
}

//Cargo la vista de añadir una lista
func ChangeViewAddList(nombreVista string, proyectID string, proyectName string, usersProyect []models.UserRole) {
	tmpl, err := template.ParseFiles(nombreVista)
	if err != nil {
		log.Fatal(err)
	}
	dataStruct := DataList{
		User:         models.UserSesion,
		ProyectID:    proyectID,
		ProyectName:  proyectName,
		UsersProyect: usersProyect,
	}
	buff := bytes.Buffer{}
	tmpl.Execute(&buff, dataStruct)
	loadableContents := "data:text/html," + url.PathEscape(buff.String())
	UI.Load(loadableContents)
}

func ChangeViewTasks(nombreVista string, tasks []models.Task, listID string, listUsers []models.UserRole, listName string, typeTasks string, pendientes int, progreso int, finalizadas int) {
	tmpl, err := template.ParseFiles(nombreVista)
	if err != nil {
		log.Fatal(err)
	}

	dataStruct := DataTask{
		Tasks:       tasks,
		ListID:      listID,
		ListName:    listName,
		ListUsers:   listUsers,
		User:        models.UserSesion.Email,
		Pendientes:  pendientes,
		Progreso:    progreso,
		Finalizadas: finalizadas,
		TypeTasks:   typeTasks,
	}
	buff := bytes.Buffer{}
	tmpl.Execute(&buff, dataStruct)
	loadableContents := "data:text/html," + url.PathEscape(buff.String())
	UI.Load(loadableContents)
}

func ChangeViewConfigProyect(nombreVista string, proyect models.Proyect, emails []string, userEmail string) {
	tmpl, err := template.ParseFiles(nombreVista)
	if err != nil {
		log.Fatal(err)
	}
	dataStruct := DataConfigProyect{
		Proyect: proyect,
		Emails:  emails,
		User:    userEmail,
	}
	buff := bytes.Buffer{}
	tmpl.Execute(&buff, dataStruct)
	loadableContents := "data:text/html," + url.PathEscape(buff.String())
	UI.Load(loadableContents)
}

func ChangeViewConfigList(nombreVista string, list models.List, emails []models.UserRole, userEmail string) {
	tmpl, err := template.ParseFiles(nombreVista)
	if err != nil {
		log.Fatal(err)
	}
	dataStruct := DataConfigList{
		List:   list,
		Emails: emails,
		User:   userEmail,
	}
	buff := bytes.Buffer{}
	tmpl.Execute(&buff, dataStruct)
	loadableContents := "data:text/html," + url.PathEscape(buff.String())
	UI.Load(loadableContents)
}

func ChangeViewConfigTask(nombreVista string, task models.Task, list models.List) {

	tmpl, err := template.ParseFiles(nombreVista)
	if err != nil {
		log.Fatal(err)
	}
	dataStruct := DataConfigTask{
		Task: task,
		List: list,
		User: models.UserSesion.Email,
		HasSignReceivedByUser: func(task models.Task) bool {
			for i := 0; i < len(task.SignsReceived); i++ {
				if task.SignsReceived[i].UserSign == models.UserSesion.Email {
					return true
				}
			}
			return false
		},
		HasSignCloseByUser: func(task models.Task) bool {
			for i := 0; i < len(task.SignsClose); i++ {
				if task.SignsClose[i].UserSign == models.UserSesion.Email {
					return true
				}
			}
			return false
		},
		CheckCreator: func(task models.Task) bool {
			return task.Creator == models.UserSesion.Email
		},
		IsUserAssigned: func(task models.Task) bool {
			for i := 0; i < len(task.Users); i++ {
				if models.UserSesion.Email == task.Users[i] {
					return true
				}
			}
			return false
		},
	}

	buff := bytes.Buffer{}
	tmpl.Execute(&buff, dataStruct)
	loadableContents := "data:text/html," + url.PathEscape(buff.String())
	UI.Load(loadableContents)
}

func LoadTask(taskID string, listID string) (models.Task, models.List) {
	taskCipher := models.GetTask(taskID, listID)
	if !taskCipher.ID.IsZero() {
		listCipher := models.GetList(listID)
		listKey := models.GetListKey(listID)
		list := models.DescifrarLista(listCipher, listKey)
		task := models.DescifrarTarea(taskCipher, listKey)
		//Limpio de los usuarios de la lista aquellos que estan en la tarea
		for i := 0; i < len(task.Users); i++ {
			list.Users = models.FindAndDeleteUsers(list.Users, task.Users[i])
		}
		return task, list
	} else {
		return models.Task{}, models.List{}
	}
}

func ChangeViewAdminPanel(nombreVista string, users []models.User) {
	tmpl, err := template.ParseFiles(nombreVista)
	if err != nil {
		log.Fatal(err)
	}

	dataStruct := AdminView{
		Users:      users,
		UserActual: models.UserSesion.Email,
	}
	buff := bytes.Buffer{}
	tmpl.Execute(&buff, dataStruct)
	loadableContents := "data:text/html," + url.PathEscape(buff.String())
	UI.Load(loadableContents)
}
