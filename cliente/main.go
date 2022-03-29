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

	//Binds Pruebas

	//Binds Finales

	//Cambiar a la pantalla de registro
	UI.Bind("changeToRegister", func() {
		ChangeView("www/register.html")
	})

	//Cambiar a la pantalla de inicio de sesion
	UI.Bind("changeToLogin", func() {
		ChangeView("www/login.html")
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

	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-UI.Done():
	}

}
