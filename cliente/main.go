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

	//Cambiar a la pantalla de inicio de sesion
	UI.Bind("changeToLogin", func() {
		ChangeView(config.PreView + "login.html")
	})

	//Cambiar a la pantalla de Home
	UI.Bind("changeHomeGO", func() bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return true
		} else {
			//Traigo del servidor todos los proyectos y listas del usuario
			models.GetUserProyectsLists()
			ChangeViewWithValues(config.PreView+"index.html", nil)
			return false
		}
	})

	//Cambiar a la pantalla de registro
	UI.Bind("changeToAdminGO", func() bool {
		users, tokenExpire := models.GetUsers()
		if tokenExpire {
			return tokenExpire
		}
		ChangeViewAdminPanel(config.PreView+"admin.html", users)
		return false
	})

	//Cambiar a la pantalla de añadir un proyecto
	UI.Bind("changeToAddProyectGO", func() bool {
		emails, tokenExpire := models.GetEmails()
		if tokenExpire {
			return tokenExpire
		}
		emails = utils.FindAndDelete(emails, models.UserSesion.Email)
		ChangeViewWithValues(config.PreView+"addProyect.html", emails)
		return false
	})

	//Cambiar a la pantalla de añadir una lista
	UI.Bind("changeToAddListGO", func(proyectID string, proyectName string) bool {
		//Recupero los usuarios del proyecto
		usersProyect, tokenExpire := models.GetUsersProyect(proyectID)
		if tokenExpire {
			return tokenExpire
		}
		ChangeViewAddList(config.PreView+"addList.html", proyectID, proyectName, usersProyect)
		return tokenExpire
	})

	//Cambiar a la pestaña de ver las tareas de una lista dado el ID de la lista
	UI.Bind("changeToTasksGO", func(listID string, listName string) bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return true
		} else {
			tasks := models.GetTasksByList(listID)
			ChangeViewTasks(config.PreView+"tasks.html", tasks, listID, nil, listName)
			return false
		}
	})

	//REVISAR SI QUIERO TRAERME DEL SERVIDOR O DE LOCAL
	//Cambiar a la pestaña de ver la configuracion del proyecto dado su ID
	UI.Bind("changeToProyectConfigGO", func(proyectID string) bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return true
		} else {
			proyectCipher := models.GetProyect(proyectID)
			proyectKey := models.GetProyectKey(proyectID, models.UserSesion.Email)
			proyect := models.DescifrarProyecto(proyectCipher, proyectKey)

			// var proyect models.Proyect
			// for i := 0; i < len(models.DatosUsuario); i++ {
			// 	if models.DatosUsuario[i].Proyecto.ID == proyectID {
			// 		proyect = models.DatosUsuario[i].Proyecto
			// 	}
			// }
			emailsNotInProyect := models.GetEmailsNotInProyect(proyect)
			ChangeViewConfigProyect(config.PreView+"configProyect.html", proyect, emailsNotInProyect, models.UserSesion.Email)
			return false
		}
	})

	//REVISAR SI QUIERO TRAERME DEL SERVIDOR O DE LOCAL
	//Cambiar a la pestaña de ver la configuracion de la lista dado su ID
	UI.Bind("changeToListConfigGO", func(listID string, proyectID string) bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return true
		} else {
			// var list models.List
			// var proyect models.Proyect
			// for i := 0; i < len(models.DatosUsuario); i++ {
			// 	if models.DatosUsuario[i].Proyecto.ID == proyectID {
			// 		proyect = models.DatosUsuario[i].Proyecto
			// 		for j := 0; j < len(models.DatosUsuario[i].Listas); j++ {
			// 			if models.DatosUsuario[i].Listas[j].ID == listID {
			// 				list = models.DatosUsuario[i].Listas[j]
			// 			}
			// 		}
			// 	}
			// }
			listCipher := models.GetList(listID)
			listKey := models.GetListKey(listID)
			list := models.DescifrarLista(listCipher, listKey)

			//Elimino los usuarios que ya pertenecen a la lista
			emailsProyect := models.GetProyect(proyectID).Users
			emailsList := list.Users
			for i := 0; i < len(emailsList); i++ {
				emailsProyect = utils.FindAndDelete(emailsProyect, emailsList[i])
			}
			ChangeViewConfigList(config.PreView+"configList.html", list, emailsProyect, models.UserSesion.Email)
			return false
		}
	})

	//Cambiar a la pestaña de añadir Tareas
	UI.Bind("changeToAddTaskGO", func(listID string) bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return true
		} else {
			listCipher := models.GetList(listID)
			list := models.DescifrarLista(listCipher, models.GetListKey(listID))
			listUsers := list.Users
			ChangeViewTasks(config.PreView+"addTask.html", nil, listID, listUsers, list.Name)
			return false
		}
	})

	//Cambiar a la pestaña de añadir Tareas
	UI.Bind("changeToTaskConfigGO", func(taskID string, listID string) bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return true
		} else {
			LoadTask(taskID, listID)
			return false
		}
	})

	/*------------------------------------------------------------------------- Binds Funcionalidades  ------------------------------------------------------------------*/

	//Registro del usuario
	UI.Bind("registerGO", func(user_pass []string) string {
		result := models.Register(user_pass[0], user_pass[1])
		return result
	})

	//Login del usuario
	UI.Bind("loginGO", func(user_pass []string) bool {
		result := models.LogIn(user_pass[0], user_pass[1])
		if result {
			utils.CheckExpirationTimeToken(models.UserSesion.Token)
			models.GetUserProyectsLists()
			ChangeViewWithValues(config.PreView+"index.html", nil)
		}
		return result
	})

	//Recuperar los emails de los usuarios registrados
	// UI.Bind("getEmailsGO", func() []string {
	// 	users := models.GetUsers()
	// 	var usersEmails []string
	// 	for i := 0; i < len(users); i++ {
	// 		usersEmails = append(usersEmails, users[i].Email)
	// 	}
	// 	return usersEmails
	// })

	//Cerrar la sesion
	UI.Bind("exitSesion", func() {
		models.LogOut()
		ChangeView(config.PreView + "login.html")
	})

	//Elimino a un usuario de todo el sistema
	UI.Bind("deleteUserGO", func(userEmail string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
		} else {
			deleteOk := models.DeleteUser(userEmail)
			return []bool{deleteOk, false}
		}
	})

	//Añadir un Proyecto
	UI.Bind("addProyectGO", func(newProyect models.Proyect) []bool {
		proyectOk, tokenExpire := models.CreateProyect(newProyect)
		return []bool{proyectOk, tokenExpire}
	})

	//Eliminar un proyecto
	UI.Bind("deleteProyectGO", func(proyectID string) []bool {
		deleteOk, tokenExpire := models.DeleteProyect(proyectID)
		return []bool{deleteOk, tokenExpire}
	})

	//Borrar un usuario de un Proyecto
	UI.Bind("deleteUserProyectGO", func(userEmail string, proyectID string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
		} else {
			deleteOk, tokenExpire := models.DeleteUserProyect(proyectID, userEmail)
			return []bool{deleteOk, tokenExpire}
		}
	})

	//Añadir un usuario al proyecto
	UI.Bind("addUserProyectGO", func(userEmail string, proyectID string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
		} else {
			addOk, tokenExpire := models.AddUserProyect(proyectID, userEmail)
			return []bool{addOk, tokenExpire}
		}
	})

	//Actualizar un proyecto
	UI.Bind("updateProyectGO", func(newProyect models.Proyect) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
		} else {
			updateOk, tokenExpire := models.UpdateProyect(newProyect)
			return []bool{updateOk, tokenExpire}
		}
	})

	//Añadir una lista
	UI.Bind("addListGO", func(list models.List, proyectID string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
		} else {
			addOk, tokenExpire := models.CreateList(list, proyectID)
			return []bool{addOk, tokenExpire}
		}
	})

	//Eliminar una lista
	UI.Bind("deleteListGO", func(listID string, proyectID string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
		} else {
			//Consigo los usuarios de la lista
			usersList := models.GetUsersList(listID)
			//Elimino la lista de todas las relaciones
			for i := 0; i < len(usersList); i++ {
				models.DeleteRelationList(proyectID, listID, usersList[i])
			}
			//Elimino la lista
			deleteOk, tokenExpire := models.DeleteList(listID)
			return []bool{deleteOk, tokenExpire}
		}
	})

	//Borrar un usuario de una Lista
	UI.Bind("deleteUserListGO", func(userEmail string, listID string, proyectID string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
		} else {
			//Borro al usuario de la lista
			models.DeleteUserList(listID, userEmail)
			//Borro la lista en la relacion donde aparece
			deleteOk := models.DeleteRelationList(proyectID, listID, userEmail)
			return []bool{deleteOk, false}
		}
	})

	//Actualizar una lista
	UI.Bind("updateListGO", func(newList models.List) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
		} else {
			deleteOk, tokenExpire := models.UpdateList(newList)
			return []bool{deleteOk, tokenExpire}
		}
	})

	//Añadir un usuario a una lista
	UI.Bind("addUserListGO", func(userEmail string, proyectID string, listID string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
		} else {
			addOK := models.AddUserList(userEmail, proyectID, listID)
			return []bool{addOK, false}
		}
	})

	//Añadir una tarea a una lista
	UI.Bind("addTaskGO", func(listID string, task models.Task, dateString string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
		} else {
			addOK := models.CreateTask(listID, task)
			return []bool{addOK, false}
		}
	})

	//Eliminar una tarea
	UI.Bind("deleteTaskGO", func(taskID string) []bool {
		addOK, tokenExpire := models.DeleteTask(taskID)
		return []bool{addOK, tokenExpire}
	})

	//Actualizar una tarea
	UI.Bind("updateTaskGO", func(newTask models.Task, listID string) []bool {

		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
		} else {
			updateOK := models.UpdateTask(listID, newTask)
			return []bool{updateOK, false}
		}
	})

	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-UI.Done():
	}

}
