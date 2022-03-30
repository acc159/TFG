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

type List struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	Cipherdata string               `bson:"cipherdata,omitempty"`
	Users      []primitive.ObjectID `bson:"users,omitempty"`
	proyectID  primitive.ObjectID   `bson:"proyectID,omitempty"`
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

func DeleteList(idString string) bool {

	//1.Borrar todas las tareas cuyo listID sea el id pasado por parametro
	DeleteByListID(idString)
	//2.Borrar la lista en si
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coleccion := config.InstanceDB.DB.Collection("lists")

	id, _ := primitive.ObjectIDFromHex(idString)
	filter := bson.D{{Key: "_id", Value: id}}
	err := coleccion.FindOneAndDelete(ctx, filter)
	if err.Err() != nil {
		return false
	}
	return true
}

func GetListsByIDs(stringsIDs []string) []Proyect {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coleccion := config.InstanceDB.DB.Collection("lists")

	var ids []primitive.ObjectID

	for i := 0; i < len(stringsIDs); i++ {
		id, _ := primitive.ObjectIDFromHex(stringsIDs[i])
		ids = append(ids, id)
	}

	var lists []Proyect
	filter := bson.M{"_id": bson.M{"$in": ids}}

	result, err := coleccion.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	err = result.All(ctx, &lists)
	if err != nil {
		log.Fatal(err)
	}

	return lists

}

//Posibles
func AddUserToList() {

}
