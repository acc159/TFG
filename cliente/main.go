package main

import (
	"cliente/config"
	"cliente/models"
	"os"
	"os/signal"
)

func main() {
	//Inicio la interfaz visual
	InitUI()
	defer UI.Close()

	models.GetTasksByList("623372b607ec94da39b16b2e")

	//Binds Cambiar Pantallas

	//Cambiar a la pantalla de registro
	UI.Bind("changeToRegister", func() {
		ChangeView("www/register.html")
	})

	//Cambiar a la pantalla de añadir un proyecto
	UI.Bind("changeToAddProyect", func() {
		ChangeViewWithValues(config.PreView + "addProyect.html")
	})

	//Cambiar a la pantalla de Home
	UI.Bind("changeHome", func() {
		models.GetUserProyectsLists()
		ChangeViewWithValues(config.PreView + "index.html")
	})

	//Cambiar a la pantalla de inicio de sesion
	UI.Bind("changeToLogin", func() {
		ChangeView("www/login.html")
	})

	//Cambiar a la pantalla de añadir una lista
	UI.Bind("changeToAddListGO", func(proyectID string) {
		//Recupero los usuarios del proyecto
		usersProyect := models.GetUsersProyect(proyectID)
		ChangeViewAddList(config.PreView+"addList.html", proyectID, usersProyect)
	})

	UI.Bind("changeToTasksGO", func(listID string) {
		tasks := models.GetTasksByList(listID)
		ChangeViewTasks(config.PreView+"tasks.html", tasks, listID)
	})

	UI.Bind("changeToProyectConfigGO", func(proyectID string) {
		proyectCipher := models.GetProyect(proyectID)
		//Descifro
		proyect := models.DescifrarProyecto(proyectCipher)
		ChangeViewConfig(config.PreView+"configProyect.html", proyect, models.List{})
	})

	UI.Bind("changeToListConfigGO", func(listID string, proyectID string) {
		listCipher := models.GetList(listID)
		proyectCipher := models.GetProyect(proyectID)
		//Descifro
		list := models.DescifrarLista(listCipher)
		proyect := models.DescifrarProyecto(proyectCipher)
		ChangeViewConfig(config.PreView+"configList.html", proyect, list)
	})

	//Binds Funcionalidades

	//Registro del usuario
	UI.Bind("registerGO", func(user_pass []string) bool {
		result := models.Register(user_pass[0], user_pass[1])
		if result {
			ChangeViewWithValues(config.PreView + "index.html")
		}
		return false
	})

	//Login del usuario
	UI.Bind("loginGO", func(user_pass []string) bool {
		result := models.LogIn(user_pass[0], user_pass[1])
		if result {
			models.GetUserProyectsLists()
			ChangeViewWithValues(config.PreView + "index.html")
		}
		return false
	})

	//Añadir un Proyecto
	UI.Bind("addProyectGO", func(newProyect models.Proyect) bool {
		return models.CreateProyect(newProyect)
	})

	//Añadir una lista
	UI.Bind("addListGO", func(list models.List, proyectID string) bool {
		listID := models.CreateList(list, proyectID)
		//Añado la lista a la relacion del proyecto para cada usuario miembro de la lista
		for i := 0; i < len(list.Users); i++ {
			models.AddListToRelation(proyectID, listID, list.Users[i])
		}
		return true
	})

	//Recuperar los emails de los usuarios registrados
	UI.Bind("getEmailsGO", func() []string {
		users := models.GetUsers()
		var usersEmails []string
		for i := 0; i < len(users); i++ {
			usersEmails = append(usersEmails, users[i].Email)
		}
		return usersEmails
	})

	//Eliminar un proyecto
	UI.Bind("deleteProyectGO", func(proyectID string) bool {
		return models.DeleteProyect(proyectID)
	})

	//Eliminar una lista
	UI.Bind("deleteListGO", func(listID string, proyectID string) bool {
		//Consigo los usuarios de la lista
		usersList := models.GetUsersList(listID)
		//Elimino todas la lista de todas las relaciones
		for i := 0; i < len(usersList); i++ {
			models.DeleteRelationList(proyectID, listID, usersList[i])
		}
		//Elimino la lista
		return models.DeleteList(listID)
	})

	//Cerrar la sesion
	UI.Bind("exitSesion", func() {
		models.LogOut()
		ChangeView("www/login.html")
	})

	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-UI.Done():
	}

}
