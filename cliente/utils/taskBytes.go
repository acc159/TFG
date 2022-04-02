package utils

import (
	"cliente/models"
	"encoding/json"
	"fmt"
)

func TaskToBytes() []byte {
	task := models.Task{
		Nombre: "ADSFSDF",
	}
	taskBytes, _ := json.Marshal(task)
	return taskBytes
}

func BytesToTask(datos []byte) models.Task {
	var task models.Task
	err := json.Unmarshal(datos, &task)
	if err != nil {
		fmt.Println("error:", err)
	}
	return task
}
