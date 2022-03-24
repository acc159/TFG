package main

import (
	"cliente/models"
	"cliente/view"
	"fmt"
	"os"
	"os/signal"
)

func main() {

	view.InitUI()
	defer view.UI.Close()

	view.UI.Bind("llamarGO", func(arrayUserPass []string) bool {
		// username := view.UI.Eval("document.querySelector('#username').value").String()
		// password := view.UI.Eval("document.querySelector('#password').value").String()
		fmt.Println(arrayUserPass[0])
		fmt.Println(arrayUserPass[1])
		return models.LogIn(arrayUserPass[0], arrayUserPass[1])
	})

	view.UI.Bind("cambiarVistaenGO", func(nombreVista string) {
		view.ChangeView(nombreVista)
	})

	//Otra alternativa
	// view.UI.Bind("recibirProyectosyListasenGO", func(id string) []models.Relation {
	// 	var relations []models.Relation
	// 	models.GetProyectsListsByUser(id, &relations)
	// 	return relations
	// })

	view.UI.Bind("recibirProyectosyListasenGO", func(id string) interface{} {
		models.GetProyectsListsByUser(id)
		return models.Relations
	})

	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-view.UI.Done():
	}

	// ui.Bind("prueba", func() []models.Task {
	// ui.Eval("alert('Estoy aqui')")

	// 	task := models.Task{
	// 		Nombre:      "<script>alert('dsafasd')</script>",
	// 		Descripcion: "asdfdf",
	// 		FechaLimite: "ASDFADSF",
	// 	}

	// 	task2 := models.Task{
	// 		Nombre:      "OTRO",
	// 		Descripcion: "OTRO",
	// 		FechaLimite: "OITERWSARI",
	// 	}

	// 	tasks := []models.Task{
	// 		task, task2,
	// 	}

	// 	return tasks
	// })

	// //Bindeo las funciones
	// ui.Bind("llamarGO", func() int {

	// 	log.Println("Me han llamado desde Javascript")
	// 	//Aqui obtengo un numero que hay en un span en el html desde GOLANG y lo imprimo aqui
	// 	numero := ui.Eval("document.querySelector('#numero').textContent").String()
	// 	fmt.Println(numero)
	// 	valor, _ := strconv.Atoi(numero)
	// 	return valor

	// })

	// //Cargo la primera pantalla

	// content, err := ioutil.ReadFile("www/login.html")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// myFileContents := string(content)
	// loadableContents := "data:text/html," + url.PathEscape(myFileContents)
	// ui.Load(loadableContents)
	// <-ui.Done()

	//models.GetUsers()
	// models.GetProyectsListsByUser()
	// models.CreateRelation("6231fae79f2ad453236d5804", "6239fac76f2ad453276c5804", "9239fac76f2ad453296c5804")
	// models.GetTasksByList()

	//models.CreateTask()
	//models.DeleteTask()
	//models.DeleteRelation("6239fac76f2ad453296c5809", "6239faf56f2ad453296c5807", "6239fb356f2ad453296c5808")
	//models.UpdateTask()
}
