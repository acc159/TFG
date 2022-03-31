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
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Cipherdata string             `bson:"cipherdata,omitempty"`
	ListID     primitive.ObjectID `bson:"listID,omitempty"`
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
	if err != nil {
		log.Println(err)
	}
	var tasks []Task
	err = result.All(ctx, &tasks)
	if err != nil {
		log.Println(err)
	}

	return tasks
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
	fmt.Println(stringObjectID)
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

//Para actualizar una tarea
func UpdateTask(newTask Task, idString string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("tasks")

	id, _ := primitive.ObjectIDFromHex(idString)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.M{"$set": newTask}
	var updatedDoc bson.D
	err := coleccion.FindOneAndUpdate(ctx, filter, update).Decode(&updatedDoc)
	if err != nil {
		log.Println(err)
	}
	if len(updatedDoc) == 0 {
		return false
	}
	return true
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
