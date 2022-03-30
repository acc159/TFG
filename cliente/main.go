package main

import (
	"cliente/config"
	"cliente/models"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	//Inicio la interfaz visual
	InitUI()
	defer UI.Close()

	//Binds Pruebas

	//Binds Finales

	//Cambiar a la pantalla de registro
	UI.Bind("changeToRegister", func() {
		ChangeView("www/register.html")
	})

	UI.Bind("changeToAddProyect", func() {
		ChangeViewWithValues(config.PreView + "addProyect.html")
	})

	UI.Bind("changeHome", func() {
		models.GetUserProyectsLists()
		ChangeViewWithValues(config.PreView + "index.html")
	})

	//Cambiar a la pantalla de inicio de sesion
	UI.Bind("changeToLogin", func() {
		ChangeView("www/login.html")
	})

	UI.Bind("changeToAddListGO", func(proyectID string, usersProyect []string) {
		ChangeViewAddList(config.PreView+"addList.html", proyectID, usersProyect)
	})

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

	UI.Bind("addProyectGO", func(newProyect models.Proyect) bool {
		return models.CreateProyect(newProyect)
	})

	UI.Bind("addListGO", func(list models.List, proyectID string) {
		listID := models.CreateList(list)

		//Añado la lista a la relacion del proyecto para cada usuario miembro de la lista
		for i := 0; i < len(list.Users); i++ {
			models.AddListToRelation(proyectID, listID)
		}

		fmt.Println(list)
		fmt.Println(proyectID)
	})

	UI.Bind("getEmailsGO", func() []string {
		users := models.GetUsers()
		var usersEmails []string
		for i := 0; i < len(users); i++ {
			usersEmails = append(usersEmails, users[i].Email)
		}
		return usersEmails
	})

	UI.Bind("deleteProyectGO", func(proyectID string, listsIDs []string) bool {
		return models.DeleteProyect(proyectID, listsIDs)
	})

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
