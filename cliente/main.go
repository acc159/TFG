package main

import (
	"cliente/models"
)

func main() {

	// ui, _ := lorca.New("", "", 480, 320)
	// defer ui.Close()

	// ui.Bind("prueba", func() []models.Task {
	// 	ui.Eval("alert('Estoy aqui')")

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
	models.DeleteRelation("6239fac76f2ad453296c5809", "6239faf56f2ad453296c5807", "6239fb356f2ad453296c5808")
}
