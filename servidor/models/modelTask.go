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
	Cipherdata string             `bson:"cipherdata"`
	ListID     primitive.ObjectID `bson:"listID,omitempty"`
}

func GetTasksByList(idString string) []Task {
	//Creo un contexto
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("tasks")

	listID, _ := primitive.ObjectIDFromHex(idString)

	filter := bson.M{"listID": listID}

	result, err := coleccion.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	var tasks []Task
	err = result.All(ctx, &tasks)
	if err != nil {
		log.Fatal(err)
	}

	return tasks
}

func CreateTask(task Task) string {
	//Creo un contexto
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
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

func DeleteTask(idString string) bool {
	//Creo un contexto
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("tasks")

	id, _ := primitive.ObjectIDFromHex(idString)
	filter := bson.D{{"_id", id}}

	result, err := coleccion.DeleteOne(ctx, filter)

	if err != nil {
		fmt.Println(err)
	}
	if result.DeletedCount == 0 {
		return false
	}
	return true
}

func UpdateTask(newTask Task, idString string) bool {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("tasks")

	id, _ := primitive.ObjectIDFromHex(idString)
	filter := bson.D{{"_id", id}}
	update := bson.M{"$set": newTask}

	result, err := coleccion.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
	}

	if result.ModifiedCount == 0 {
		return false
	} else {
		return true
	}
}

func DeleteByListID(idString string) bool {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("tasks")

	id, _ := primitive.ObjectIDFromHex(idString)
	filter := bson.D{{"listID", id}}
	results, err := coleccion.DeleteMany(ctx, filter)
	if err != nil {
		fmt.Println(err)
	}
	if results.DeletedCount == 0 {
		return false
	}
	return true
}
