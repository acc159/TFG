package main

import (
	"cliente/config"
	"cliente/models"
	"cliente/utils"
	"fmt"
	"os"
	"os/signal"
)

func pruebaFirma() {
	signature := utils.Sign([]byte("Hola mundo"), models.GetPrivateKeyUser())
	resultado := utils.CheckSign(signature, []byte("Hola mundo"), utils.PemToPublicKey(models.UserSesion.PublicKey))
	fmt.Println(resultado)
}

func main() {
	//Inicio la interfaz visual
	InitUI()
	defer UI.Close()

	/*------------------------------------------------------------------------- Binds Cambiar Pantallas  ------------------------------------------------------------------*/

	//Cambiar a la pantalla de registro
	UI.Bind("changeToRegister", func() {
		ChangeView(config.PreView + "register.html")
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
		ChangeView(config.PreView + "login.html")
	})

	//Cambiar a la pantalla de añadir una lista
	UI.Bind("changeToAddListGO", func(proyectID string) {
		//Recupero los usuarios del proyecto
		usersProyect := models.GetUsersProyect(proyectID)
		ChangeViewAddList(config.PreView+"addList.html", proyectID, usersProyect)
	})

	//Cambiar a la pestaña de ver las tareas de una lista dado el ID de la lista
	UI.Bind("changeToTasksGO", func(listID string) {
		tasks := models.GetTasksByList(listID)
		ChangeViewTasks(config.PreView+"tasks.html", tasks, listID, nil)
	})

	//REVISAR SI QUIERO TRAERME DEL SERVIDOR O DE LOCAL
	//Cambiar a la pestaña de ver la configuracion del proyecto dado su ID
	UI.Bind("changeToProyectConfigGO", func(proyectID string) {
		var proyect models.Proyect
		for i := 0; i < len(models.DatosUsuario); i++ {
			if models.DatosUsuario[i].Proyecto.ID == proyectID {
				proyect = models.DatosUsuario[i].Proyecto
			}
		}
		emailsNotInProyect := models.GetEmailsNotInProyect(proyect)
		ChangeViewConfig(config.PreView+"configProyect.html", proyect, models.List{}, emailsNotInProyect, models.UserSesion.Email)
	})

	//REVISAR SI QUIERO TRAERME DEL SERVIDOR O DE LOCAL
	//Cambiar a la pestaña de ver la configuracion de la lista dado su ID
	UI.Bind("changeToListConfigGO", func(listID string, proyectID string) {
		var list models.List
		var proyect models.Proyect
		for i := 0; i < len(models.DatosUsuario); i++ {
			if models.DatosUsuario[i].Proyecto.ID == proyectID {
				proyect = models.DatosUsuario[i].Proyecto
				for j := 0; j < len(models.DatosUsuario[i].Listas); j++ {
					if models.DatosUsuario[i].Listas[j].ID == listID {
						list = models.DatosUsuario[i].Listas[j]
					}
				}
			}
		}
		//Elimino los usuarios que ya pertenecen a la lista
		emailsProyect := models.GetProyect(proyectID).Users
		emailsList := list.Users
		for i := 0; i < len(emailsList); i++ {
			emailsProyect = utils.FindAndDelete(emailsProyect, emailsList[i])
		}
		ChangeViewConfig(config.PreView+"configList.html", proyect, list, emailsProyect, models.UserSesion.Email)
	})

	//Cambiar a la pestaña de añadir Tareas
	UI.Bind("changeToAddTaskGO", func(listID string) {
		listCipher := models.GetList(listID)
		list := models.DescifrarLista(listCipher, models.GetListKey(listID))
		listUsers := list.Users
		ChangeViewTasks(config.PreView+"addTask.html", nil, listID, listUsers)
	})

	//Cambiar a la pestaña de añadir Tareas
	UI.Bind("changeToTaskConfigGO", func(taskID string, listID string) {
		LoadTask(taskID, listID)
	})

	/*------------------------------------------------------------------------- Binds Funcionalidades  ------------------------------------------------------------------*/

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

	//Recuperar los emails de los usuarios registrados
	UI.Bind("getEmailsGO", func() []string {
		users := models.GetUsers()
		var usersEmails []string
		for i := 0; i < len(users); i++ {
			usersEmails = append(usersEmails, users[i].Email)
		}
		return usersEmails
	})

	//Cerrar la sesion
	UI.Bind("exitSesion", func() {
		models.LogOut()
		ChangeView(config.PreView + "login.html")
	})
	//Elimino a un usuario de todo el sistema
	UI.Bind("deleteUserGO", func() {
		models.DeleteUser(models.UserSesion.Email)
		ChangeView(config.PreView + "login.html")
	})

	//Añadir un Proyecto
	UI.Bind("addProyectGO", func(newProyect models.Proyect) bool {
		return models.CreateProyect(newProyect)
	})

	//Eliminar un proyecto
	UI.Bind("deleteProyectGO", func(proyectID string) bool {
		return models.DeleteProyect(proyectID)
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

	//Actualizar un proyecto
	UI.Bind("updateProyectGO", func(newProyect models.Proyect) bool {
		return models.UpdateProyect(newProyect)
	})

	//Añadir una lista
	UI.Bind("addListGO", func(list models.List, proyectID string) bool {
		models.CreateList(list, proyectID)
		return true
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
			list := models.DescifrarLista(listCipher, models.GetListKey(listID))
			proyect := models.DescifrarProyecto(proyectCipher, models.GetProyectKey(proyectID, userEmail))
			emailsProyect := models.GetProyect(proyectID).Users
			ChangeViewConfig(config.PreView+"configList.html", proyect, list, emailsProyect, models.UserSesion.Email)
		}
	})

	//Actualizar una lista
	UI.Bind("updateListGO", func(newList models.List) bool {
		return models.UpdateList(newList)
	})

	//Añadir un usuario a una lista
	UI.Bind("addUserListGO", func(userEmail string, proyectID string, listID string) bool {
		return models.AddUserList(userEmail, proyectID, listID)
	})

	//Añadir una tarea a una lista
	UI.Bind("addTaskGO", func(listID string, task models.Task, dateString string) bool {
		return models.CreateTask(listID, task)
	})

	//Eliminar una tarea
	UI.Bind("deleteTaskGO", func(taskID string) bool {
		return models.DeleteTask(taskID)
	})

	//Actualizar una tarea
	UI.Bind("updateTaskGO", func(newTask models.Task, listID string) bool {
		//return models.UpdateList(newList)
		resultado := models.UpdateTask(listID, newTask)
		if resultado {
			return true
		} else {
			return false
		}
	})

	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-UI.Done():
	}

}
