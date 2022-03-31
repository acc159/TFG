package models

import (
	"bytes"
	"cliente/config"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Nombre           string             `bson:"nombre"`
	Descripcion      string             `bson:"apellidos"`
	Fecha            primitive.DateTime `bson:"fecha"`
	Estado           string             `bson:"estado"`
	ArchivosAdjuntos string             `bson:"archivos_adjuntos"`
	EnlacesAdjuntos  string             `bson:"enlaces_adjuntos"`
}

type TaskCipher struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Cipherdata string             `bson:"cipherdata,omitempty"`
	ListID     primitive.ObjectID `bson:"listID,omitempty"`
}

//Recupero las tareas por su ListID
func GetTasksByList(listID string) []Task {
	resp, err := http.Get(config.URLbase + "tasks/" + listID)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var tasks []Task
	if resp.StatusCode == 400 {
		fmt.Println("Ningun tarea para dicha lista")
		return tasks
	} else {
		var tasksCipher []TaskCipher
		json.NewDecoder(resp.Body).Decode(&tasksCipher)
		var tasks []Task
		for i := 0; i < len(tasksCipher); i++ {
			tasks = append(tasks, DescifrarTarea(tasksCipher[i]))
		}
		return tasks
	}
}

func CreateTask() {

	listID, _ := primitive.ObjectIDFromHex("6239fb356f2ad453296c5807")

	//Creo la relacion a enviar
	task := TaskCipher{
		Cipherdata: "CREADA DESDE EL CLIENTE",
		ListID:     listID,
	}

	//Pasamos el tipo Relation a JSON
	taskJSON, err := json.Marshal(task)
	if err != nil {
		fmt.Println(err)
	}

	//Peticion POST
	url := config.URLbase + "task"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(taskJSON))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		fmt.Println("La tarea no pudo ser creada")
	} else {
		var newTaskID string
		json.NewDecoder(resp.Body).Decode(&newTaskID)
		fmt.Println(newTaskID)
	}

}

func UpdateTask() {

	listID, _ := primitive.ObjectIDFromHex("6239fb356f2ad453296c5807")
	task := TaskCipher{
		Cipherdata: "ACTUALIZADA",
		ListID:     listID,
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		fmt.Println(err)
	}

	url := config.URLbase + "tasks/" + "623b64ef2945dc21f09354d9"

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(taskJSON))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		fmt.Println("La tarea no pudo ser actualizada")
	} else {
		var newTaskID string
		json.NewDecoder(resp.Body).Decode(&newTaskID)
		fmt.Println(newTaskID)
	}
}

func DeleteTask() {

	taskID := "6239fb826f2ad453296c580a"

	url := config.URLbase + "tasks/" + taskID
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		fmt.Println("La tarea no pudo ser borrada")
	} else {
		var resultado string
		json.NewDecoder(resp.Body).Decode(&resultado)
		fmt.Println(resultado)
	}
}

func DescifrarTarea(taskCipher TaskCipher) Task {
	return Task{
		ID:          taskCipher.ID,
		Nombre:      "Nombre de la tarea",
		Descripcion: "Descripcion de la tarea",
		Estado:      "En progreso",
	}
}

func CifrarTarea(task Task) TaskCipher {
	return TaskCipher{
		Cipherdata: "DASFSDFASDFSDF",
	}
}
