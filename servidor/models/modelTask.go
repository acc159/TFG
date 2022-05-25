package models

import (
	"context"
	"fmt"
	"log"
	"servidor/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Cipherdata  []byte             `bson:"cipherdata,omitempty"`
	ListID      primitive.ObjectID `bson:"listID,omitempty"`
	Check       string             `bson:"check,omitempty"`
	UpdateCheck string             `bson:"updateCheck,omitempty"`
}

//Recupero las tareas que pertenecen a una lista
func GetTasksByList(idString string) []Task {
	//Creo un contexto
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("tasks")
	listID, _ := primitive.ObjectIDFromHex(idString)
	filter := bson.M{"listID": listID}
	result, err := coleccion.Find(ctx, filter)
	var tasks []Task
	if err != nil {
		log.Println(err)
		return tasks
	}
	err = result.All(ctx, &tasks)
	if err != nil {
		log.Println(err)
		return tasks
	}
	return tasks
}

func GetTaskByID(taskIDstring string) Task {
	//Creo un contexto
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("tasks")
	taskID, _ := primitive.ObjectIDFromHex(taskIDstring)
	filter := bson.M{"_id": taskID}
	var task Task
	err := coleccion.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		log.Println(err)
	}
	return task
}

//Crear una tarea
func CreateTask(task Task) string {
	//Creo un contexto
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("tasks")
	result, err := coleccion.InsertOne(ctx, task)
	if err != nil {
		fmt.Println(err)
	}
	stringObjectID := result.InsertedID.(primitive.ObjectID).Hex()
	return stringObjectID
}

//Para borrar una tarea
func DeleteTask(idString string) bool {
	//Creo un contexto
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("tasks")
	id, _ := primitive.ObjectIDFromHex(idString)
	filter := bson.D{{Key: "_id", Value: id}}
	err := coleccion.FindOneAndDelete(ctx, filter)
	if err.Err() != nil {
		return false
	}
	return true
}

func CheckModificationTask(taskID string, Updatecheck string) bool {
	task := GetTaskByID(taskID)
	if task.ID.IsZero() || task.Check != Updatecheck {
		return true
	}
	return false
}

//Para actualizar una tarea
func UpdateTask(newTask Task, idString string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("tasks")
	//Primero compruebo si se realizo alguna actualizacion o borrado antes
	if CheckModificationTask(idString, newTask.UpdateCheck) {
		return "Ya modificada"
	}
	newTask.UpdateCheck = ""
	id, _ := primitive.ObjectIDFromHex(idString)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.M{"$set": newTask}
	var updatedDoc bson.D
	err := coleccion.FindOneAndUpdate(ctx, filter, update).Decode(&updatedDoc)
	if err != nil {
		log.Println(err)
	}
	if len(updatedDoc) == 0 {
		return "Error"
	}
	return "OK"
}

//Para eliminar todas las tareas que pertenezcan a la misma Lista
func DeleteTasksByListID(idString string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("tasks")
	id, _ := primitive.ObjectIDFromHex(idString)
	filter := bson.D{{Key: "listID", Value: id}}
	results, err := coleccion.DeleteMany(ctx, filter)
	if err != nil {
		fmt.Println(err)
	}
	if results.DeletedCount == 0 {
		return false
	}
	return true
}
