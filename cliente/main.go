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

	UI.Bind("deleteUserGO", func() {
		models.DeleteUser(models.UserSesion.Email)
		ChangeView("www/login.html")
	})

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
			models.AddListToRelation(proyectID, listID, list.Users[i], "Clave para la lista Cifrada")
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

	//Borrar un usuario de un Proyecto
	UI.Bind("deleteUserProyectGO", func(userEmail string, proyectID string) {
		//Borro al usuario del proyecto
		models.DeleteUserProyect(proyectID, userEmail)
		//Borro la relacion
		models.DeleteUserProyectRelation(userEmail, proyectID)
	})

	//Añadir un usuario al proyecto
	UI.Bind("addUserProyectGO", func(userEmail string, proyectID string) bool {
		return models.AddUserProyect(proyectID, userEmail)
	})

	//Borrar un usuario de una Lista
	UI.Bind("deleteUserListGO", func(userEmail string, listID string, proyectID string) {
		//Borro al usuario de la lista
		models.DeleteUserList(listID, userEmail)
		//Borro la lista en la relacion donde aparece
		models.DeleteRelationList(proyectID, listID, userEmail)
		if userEmail == models.UserSesion.Email {
			//Si el usuario borrado es el que es el actual de la sesion lo redirijo a la home
			models.GetUserProyectsLists()
			ChangeViewWithValues(config.PreView + "index.html")
		} else {
			listCipher := models.GetList(listID)
			proyectCipher := models.GetProyect(proyectID)
			//Descifro
			list := models.DescifrarLista(listCipher)
			proyect := models.DescifrarProyecto(proyectCipher)
			ChangeViewConfig(config.PreView+"configList.html", proyect, list)
		}
	})

	UI.Bind("addUserListGO", func(userEmail string, proyectID string, listID string) {
		models.AddUserList(userEmail, proyectID, listID)
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
