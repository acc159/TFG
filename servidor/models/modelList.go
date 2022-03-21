package models

import (
	"context"
	"fmt"
	"servidor/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type List struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	Cipherdata string               `bson:"cipherdata"`
	Users      []primitive.ObjectID `bson:"users,omitempty"`
}

func CreateList(list List) bool {
	//Creo un contexto
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("lists")
	_, err := coleccion.InsertOne(ctx, list)

	if err != nil {
		fmt.Println(err)
		return false
	}
	return true

}

func UpdateList(list List, idString string) bool {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coleccion := config.InstanceDB.DB.Collection("lists")

	id, _ := primitive.ObjectIDFromHex(idString)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.M{"$set": list}

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

func DeleteList(idString string) bool {

	//1.Borrar todas las tareas cuyo listID sea el id pasado por parametro
	DeleteByListID(idString)
	//2.Borrar la lista en si
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coleccion := config.InstanceDB.DB.Collection("lists")

	id, _ := primitive.ObjectIDFromHex(idString)
	filter := bson.D{{Key: "_id", Value: id}}
	result, err := coleccion.DeleteOne(ctx, filter)

	if err != nil {
		fmt.Println(err)
	}
	if result.DeletedCount == 0 {
		return false
	}
	return true
}

//Posibles
func AddUserToList() {

}

func GetList() {

}
