package main

import (
	"cliente/config"
	"cliente/models"
	"cliente/utils"
	"os"
	"os/signal"
)

func main() {
	//Inicio la interfaz visual
	InitUI()
	defer UI.Close()

	/*------------------------------------------------------------------------- Binds Cambiar Pantallas  ------------------------------------------------------------------*/

	//Cambiar a la pantalla de registro
	UI.Bind("changeToRegister", func() {
		ChangeView(config.PreView + "user/register.html")
	})

	//Cambiar a la pantalla de inicio de sesion
	UI.Bind("changeToLogin", func() {
		ChangeView(config.PreView + "user/login.html")
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
		ChangeViewAdminPanel(config.PreView+"user/admin.html", users)
		return false
	})

	//Cambiar a la pantalla de añadir un proyecto
	UI.Bind("changeToAddProyectGO", func() bool {
		emails, tokenExpire := models.GetEmails()
		if tokenExpire {
			return tokenExpire
		}
		emails = utils.FindAndDelete(emails, models.UserSesion.Email)
		emails = utils.FindAndDelete(emails, "admin")
		ChangeViewWithValues(config.PreView+"proyect/addProyect.html", emails)
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
				ChangeViewAddList(config.PreView+"list/addList.html", proyectID, proyectName, usersProyect)
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

				ChangeViewTasks(config.PreView+"task/tasks.html", tasks, listID, nil, listName, "")
			}
			return []bool{false, false}
		}
	})

	UI.Bind("changeToTasksCustomGO", func(listID string, listName string, typeTasks string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, true}
		} else {
			//Compruebo que la lista existe
			if models.ExistList(listID) {
				tasks, tasksCipher := models.GetTasksByList(listID)

				var tasksFilter []models.Task
				var tasksCipherFilter []models.TaskCipher

				for i := 0; i < len(tasks); i++ {
					if tasks[i].State == typeTasks {
						tasksFilter = append(tasksFilter, tasks[i])
						tasksCipherFilter = append(tasksCipherFilter, tasksCipher[i])
					}
				}
				models.TasksLocal = tasksCipherFilter
				ChangeViewTasks(config.PreView+"task/tasks.html", tasksFilter, listID, nil, listName, typeTasks)
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
				emailsNotInProyect = utils.FindAndDelete(emailsNotInProyect, "admin")
				ChangeViewConfigProyect(config.PreView+"proyect/configProyect.html", proyect, emailsNotInProyect, models.UserSesion.Email)
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
					emailsProyect = models.FindAndDeleteUsers(emailsProyect, emailsList[i].User)
				}
				ChangeViewConfigList(config.PreView+"list/configList.html", list, emailsProyect, models.UserSesion.Email)
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
				ChangeViewTasks(config.PreView+"task/addTask.html", nil, listID, listUsers, list.Name, "")
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
				ChangeViewConfigTask(config.PreView+"task/configTask.html", task, list)
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
	UI.Bind("loginGO", func(user_pass []string) string {
		result := models.LogIn(user_pass[0], user_pass[1])

		if result == "OK" {
			utils.CheckExpirationTimeToken(models.UserSesion.Token)
			models.GetUserProyectsLists()
			ChangeViewWithValues(config.PreView+"index.html", nil)
		}
		return result
	})

	//Cerrar la sesion
	UI.Bind("exitSesion", func() {
		models.LogOut()
		ChangeView(config.PreView + "user/login.html")
	})

	//Elimino a un usuario de todo el sistema
	UI.Bind("UpdateStatusGO", func(userEmail string, status string) []bool {
		blockOK, tokenExpire := models.UpdateStatus(userEmail, status)
		return []bool{blockOK, tokenExpire}
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
			return []bool{false, false, false, true}
		} else if models.ExistProyect(proyectID) {
			if !models.CheckUserOnProyect(proyectID, userEmail) {
				deleteOk, tokenExpire := models.DeleteUserProyect(proyectID, userEmail)
				return []bool{true, deleteOk, true, tokenExpire}
			} else {
				return []bool{true, true, false, false}
			}
		}
		return []bool{false, false, false, false}
	})

	//Añadir un usuario al proyecto
	UI.Bind("addUserProyectGO", func(userEmail string, proyectID string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, false, false, true}
		} else if models.ExistProyect(proyectID) {
			if models.CheckUserOnProyect(proyectID, userEmail) {
				addOk, tokenExpire := models.AddUserProyect(proyectID, userEmail)
				return []bool{true, addOk, true, tokenExpire}
			} else {
				return []bool{true, true, false, false}
			}
		}
		return []bool{false, false, false, false}
	})

	//Actualizar un proyecto
	UI.Bind("updateProyectGO", func(newProyect models.Proyect) []string {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []string{"false", "false", "true"}
		} else {
			var proyectOK string
			if models.ExistProyect(newProyect.ID) {
				proyectOK = "True"
				updateStatus := models.UpdateProyect(newProyect)
				return []string{proyectOK, updateStatus, "false"}
			} else {
				proyectOK = "false"
				return []string{proyectOK, "false", "false"}
			}

		}
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

	//Actualizar una lista
	UI.Bind("updateListGO", func(newList models.List) []string {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []string{"false", "false", "true"}
		} else {
			listOK, updateStatus, tokenExpire := models.UpdateList(newList)
			return []string{listOK, updateStatus, tokenExpire}
		}
	})

	//Borrar un usuario de una Lista
	UI.Bind("deleteUserListGO", func(userEmail string, listID string, proyectID string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, false, false, true}
		} else if models.ExistProyect(proyectID) {
			if !models.CheckUserOnList(listID, userEmail) {
				//Borro al usuario de la lista
				models.DeleteUserList(listID, userEmail)
				//Borro la lista en la relacion donde aparece
				deleteOk := models.DeleteRelationList(proyectID, listID, userEmail)
				return []bool{true, deleteOk, true, false}
			} else {
				return []bool{true, true, false, false}
			}
		}
		return []bool{false, false, false, false}
	})

	//Añadir un usuario a una lista
	UI.Bind("addUserListGO", func(userEmail string, proyectID string, listID string) []bool {
		if !utils.CheckExpirationTimeToken(models.UserSesion.Token) {
			return []bool{false, false, false, true}
		} else if models.ExistProyect(proyectID) {
			if models.CheckUserOnList(listID, userEmail) {
				addOk := models.AddUserList(userEmail, proyectID, listID)
				return []bool{true, addOk, true, false}
			} else {
				return []bool{true, true, false, false}
			}
		}
		return []bool{false, false, false, false}
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

	UI.Bind("verifySignDataTaskGO", func(data string, linkName string, taskID string, listID string, userSign string, dataType string) []bool {
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
				certificate := models.GetCertificateUser(userSign)
				publicKey := utils.PemToPublicKey(certificate.PublicKeyUser)
				return []bool{utils.CheckSign(sign, []byte(data), publicKey), false}
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
				certificate := models.GetCertificateUser(userSign)
				publicKey := utils.PemToPublicKey(certificate.PublicKeyUser)
				if models.VerifyCertificateSign(certificate) && models.VerifyPublicKeyWithCertificate(certificate) {
					return []bool{utils.CheckSign(sign, []byte(data), publicKey), false}
				} else {
					return []bool{false, false}
				}
			} else {
				return []bool{false, false}
			}
		}
	})

	UI.Bind("refreshTokenGO", func() bool {
		refreshOK := models.RefreshTokenUser()
		return refreshOK
	})

	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-UI.Done():
	}

}
