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
	FechaLimite      string             `bson:"fecha_limite"`
	ArchivosAdjuntos string             `bson:"archivos_adjuntos"`
	EnlacesAdjuntos  string             `bson:"enlaces_adjuntos"`
}

type TaskCipher struct {
	Cipherdata string             `bson:"cipherdata,omitempty"`
	ListID     primitive.ObjectID `bson:"listID,omitempty"`
}

func GetTasksByList() {
	listID := "6239fb356f2ad453296c5807"

	resp, err := http.Get(config.URLbase + "tasks/" + listID)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		fmt.Println("Ningun tarea para dicha lista")
	} else {
		var responseObject []Task
		json.NewDecoder(resp.Body).Decode(&responseObject)
		fmt.Println(responseObject)
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
	relationJSON, err := json.Marshal(task)
	if err != nil {
		fmt.Println(err)
	}

	//Peticion POST
	url := config.URLbase + "task"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(relationJSON))
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
