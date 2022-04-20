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
	UI.Bind("changeToAddListGO", func(proyectID string, proyectName string) []bool {

		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
		} else {
			if models.ExistProyect(proyectID) {
				//Recupero los usuarios del proyecto
				usersProyect, _ := models.GetUsersProyect(proyectID)
				ChangeViewAddList(config.PreView+"addList.html", proyectID, proyectName, usersProyect)
			}
			return []bool{false, false}
		}
	})

	//Cambiar a la pestaña de ver las tareas de una lista dado el ID de la lista
	UI.Bind("changeToTasksGO", func(listID string, listName string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
		} else {
			//Compruebo que la lista existe
			if models.ExistList(listID) {
				tasks, tasksCipher := models.GetTasksByList(listID)
				models.TasksLocal = tasksCipher

				ChangeViewTasks(config.PreView+"tasks.html", tasks, listID, nil, listName)
			}
			return []bool{false, false}
		}
	})

	//REVISAR SI QUIERO TRAERME DEL SERVIDOR O DE LOCAL
	//Cambiar a la pestaña de ver la configuracion del proyecto dado su ID
	UI.Bind("changeToProyectConfigGO", func(proyectID string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
		} else {
			if models.ExistProyect(proyectID) {
				proyectKey := models.GetProyectKey(proyectID, models.UserSesion.Email)
				proyectCipher := models.GetProyect(proyectID)
				proyect := models.DescifrarProyecto(proyectCipher, proyectKey)
				emailsNotInProyect := models.GetEmailsNotInProyect(proyect)
				ChangeViewConfigProyect(config.PreView+"configProyect.html", proyect, emailsNotInProyect, models.UserSesion.Email)
				return []bool{true, false}
			}
			return []bool{false, false}
		}
	})

	//REVISAR SI QUIERO TRAERME DEL SERVIDOR O DE LOCAL
	//Cambiar a la pestaña de ver la configuracion de la lista dado su ID
	UI.Bind("changeToListConfigGO", func(listID string, proyectID string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
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
			listKey := models.GetListKey(listID)
			if len(listKey) > 0 {
				listCipher := models.GetList(listID)
				list := models.DescifrarLista(listCipher, listKey)
				//Elimino los usuarios que ya pertenecen a la lista
				emailsProyect := models.GetProyect(proyectID).Users
				emailsList := list.Users
				for i := 0; i < len(emailsList); i++ {
					emailsProyect = utils.FindAndDelete(emailsProyect, emailsList[i])
				}
				ChangeViewConfigList(config.PreView+"configList.html", list, emailsProyect, models.UserSesion.Email)
				return []bool{true, false}
			}
			return []bool{false, false}
		}
	})

	//Cambiar a la pestaña de añadir Tareas
	UI.Bind("changeToAddTaskGO", func(listID string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
		} else {
			if models.ExistList(listID) {
				listCipher := models.GetList(listID)
				list := models.DescifrarLista(listCipher, models.GetListKey(listID))
				listUsers := list.Users
				ChangeViewTasks(config.PreView+"addTask.html", nil, listID, listUsers, list.Name)
				return []bool{true, false}
			}
			return []bool{false, false}
		}
	})

	//Cambiar a la pestaña de añadir Tareas
	UI.Bind("changeToTaskConfigGO", func(taskID string, listID string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, false, true}
		} else if models.ExistList(listID) {
			task, list := LoadTask(taskID, listID)
			if task.ID != "" && list.ID != "" {
				models.CurrentTask = task
				ChangeViewConfigTask(config.PreView+"configTask.html", task, list)
				return []bool{true, true, false}
			}
			return []bool{true, false, false}
		}
		return []bool{false, false, false}
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
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, false, true}
		} else if models.ExistProyect(proyectID) {
			deleteOk, tokenExpire := models.DeleteProyect(proyectID)
			return []bool{true, deleteOk, tokenExpire}
		}
		return []bool{false, false, false}
	})

	//Borrar un usuario de un Proyecto
	UI.Bind("deleteUserProyectGO", func(userEmail string, proyectID string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, false, true}
		} else if models.ExistProyect(proyectID) {
			deleteOk, tokenExpire := models.DeleteUserProyect(proyectID, userEmail)
			return []bool{true, deleteOk, tokenExpire}
		}
		return []bool{false, false, false}
	})

	//Añadir un usuario al proyecto
	UI.Bind("addUserProyectGO", func(userEmail string, proyectID string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, false, true}
		} else if models.ExistProyect(proyectID) {
			addOk, tokenExpire := models.AddUserProyect(proyectID, userEmail)
			return []bool{true, addOk, tokenExpire}
		}
		return []bool{false, false, false}
	})

	//Actualizar un proyecto
	UI.Bind("updateProyectGO", func(newProyect models.Proyect) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, false, true}
		} else if models.ExistProyect(newProyect.ID) {
			updateOk, tokenExpire := models.UpdateProyect(newProyect)
			return []bool{true, updateOk, tokenExpire}
		}
		return []bool{false, false, false}
	})

	//Añadir una lista
	UI.Bind("addListGO", func(list models.List, proyectID string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, false, true}
		} else if models.ExistProyect(proyectID) {
			addOk, tokenExpire := models.CreateList(list, proyectID)
			return []bool{true, addOk, tokenExpire}
		}
		return []bool{false, false, false}
	})

	//Eliminar una lista
	UI.Bind("deleteListGO", func(listID string, proyectID string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
		} else {
			if models.ExistList(listID) {
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
			return []bool{false, false}
		}
	})

	//Borrar un usuario de una Lista
	UI.Bind("deleteUserListGO", func(userEmail string, listID string, proyectID string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, false, true}
		} else if models.ExistList(listID) {

			//Borro al usuario de la lista
			models.DeleteUserList(listID, userEmail)
			//Borro la lista en la relacion donde aparece
			deleteOk := models.DeleteRelationList(proyectID, listID, userEmail)
			return []bool{true, deleteOk, false}
		} else {
			return []bool{false, false, false}
		}
	})

	//Actualizar una lista
	UI.Bind("updateListGO", func(newList models.List) []string {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []string{"false", "true"}
		} else {

			listOK, updateStatus, tokenExpire := models.UpdateList(newList)
			return []string{listOK, updateStatus, tokenExpire}
		}
	})

	//Añadir un usuario a una lista
	UI.Bind("addUserListGO", func(userEmail string, proyectID string, listID string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, false, true}
		} else if models.ExistList(listID) {
			addOK := models.AddUserList(userEmail, proyectID, listID)
			return []bool{true, addOK, false}
		} else {
			return []bool{false, false, false}
		}
	})

	//Añadir una tarea a una lista
	UI.Bind("addTaskGO", func(listID string, task models.Task, dateString string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, false, true}
		} else {
			if models.ExistList(listID) {
				addOK := models.CreateTask(listID, task)
				return []bool{true, addOK, false}
			}
			return []bool{false, false, false}
		}
	})

	//Eliminar una tarea
	UI.Bind("deleteTaskGO", func(taskID string, listID string) []bool {

		if models.ExistList(listID) {
			addOK, tokenExpire := models.DeleteTask(taskID)
			return []bool{true, addOK, tokenExpire}
		}
		return []bool{false, false, false}
	})

	//Eliminar una tarea
	UI.Bind("deleteTaskByListGO", func(listID string) []bool {
		deleteOk, tokenExpire := models.DeleteTaskByListID(listID)
		return []bool{deleteOk, tokenExpire}
	})

	//Actualizar una tarea
	UI.Bind("updateTaskGO", func(newTask models.Task, listID string) []string {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []string{"false", "false", "true"}
		} else {
			if models.ExistList(listID) {
				updateStatus := models.UpdateTask(listID, newTask)
				return []string{"true", updateStatus, "false"}
			}
			return []string{"false", "false", "false"}
		}
	})

	UI.Bind("checkChangesGO", func() bool {
		return models.CheckChanges()
	})

	UI.Bind("checkTaskChangesGO", func(listID string) bool {
		return models.CheckTaskChanges(listID)
	})

	UI.Bind("signGO", func(fileData string) string {
		sign := models.SignFile(fileData)
		signBase64 := utils.ToBase64FromByte(sign)
		return signBase64
	})

	// UI.Bind("verifyFileSignGO", func(data string, fileName string, taskID string, listID string, signUser string) []bool {
	// 	if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
	// 		return []bool{false, true}
	// 	} else {
	// 		signBase64 := models.GetSignFile(taskID, listID, fileName)
	// 		if signBase64 != "" {
	// 			sign := utils.ToByteFromBase64(signBase64)
	// 			publicKey, tokenExpire := models.GetPublicKey(signUser)
	// 			if tokenExpire {
	// 				return []bool{false, tokenExpire}
	// 			}
	// 			return []bool{utils.CheckSign(sign, []byte(data), publicKey), tokenExpire}
	// 		} else {
	// 			return []bool{false, false}
	// 		}
	// 	}
	// })

	UI.Bind("verifySignDataTaskGO", func(data string, linkName string, taskID string, listID string, signUser string, dataType string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
		} else {
			var signBase64 string
			switch dataType {
			case "link":
				signBase64 = models.GetLinkFile(taskID, listID, linkName)
			case "file":
				signBase64 = models.GetSignFile(taskID, listID, linkName)
			}
			if signBase64 != "" {
				sign := utils.ToByteFromBase64(signBase64)
				publicKey, tokenExpire := models.GetPublicKey(signUser)
				if tokenExpire {
					return []bool{false, tokenExpire}
				}
				return []bool{utils.CheckSign(sign, []byte(data), publicKey), tokenExpire}
			} else {
				return []bool{false, false}
			}
		}
	})

	UI.Bind("verifySignEventGO", func(evenType string, data string, userSign string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
		} else {
			var signBase64 string
			switch evenType {
			case "Creacion":
				signBase64 = models.GetEventCreation()
			case "Recepcion":
				signBase64 = models.GetEventReceived(userSign)
			default:
				signBase64 = models.GetEventClosed(userSign)
			}
			if signBase64 != "" {
				sign := utils.ToByteFromBase64(signBase64)
				publicKey, tokenExpire := models.GetPublicKey(userSign)
				if tokenExpire {
					return []bool{false, tokenExpire}
				}
				return []bool{utils.CheckSign(sign, []byte(data), publicKey), tokenExpire}
			} else {
				return []bool{false, false}
			}
		}
	})

	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-UI.Done():
	}

}
