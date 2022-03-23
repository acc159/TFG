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

type Proyect struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Cipherdata string             `bson:"cipherdata,omitempty"`
	Users      []string           `bson:"users,omitempty"`
}

func GetProyects() []Proyect {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	coleccion := config.InstanceDB.DB.Collection("proyects")

	//Consulto a la base de datos
	result, err := coleccion.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	//Lo paso a struct de GO la consulta
	var proyectos []Proyect
	err = result.All(ctx, &proyectos)
	if err != nil {
		log.Fatal(err)
	}

	return proyectos
}

func GetProyect(idString string) Proyect {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coleccion := config.InstanceDB.DB.Collection("proyects")
	id, _ := primitive.ObjectIDFromHex(idString)
	var proyecto Proyect
	err := coleccion.FindOne(ctx, bson.M{"_id": id}).Decode(&proyecto)
	if err != nil {
		log.Println(err)
	}
	return proyecto
}

func CreateProyect(proyecto Proyect) {
	//Creo un contexto
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("proyects")
	//Inserto el usuario pasado por parametro
	result, err := coleccion.InsertOne(ctx, proyecto)
	if err != nil {
		fmt.Println(err)
	} else {
		//Paso el primitive.ObjectID a un string
		stringObjectID := result.InsertedID.(primitive.ObjectID).Hex()
		fmt.Println(stringObjectID)
	}
}

func UpdateProyect(proyecto Proyect, stringID string) bool {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coleccion := config.InstanceDB.DB.Collection("proyects")

	id, _ := primitive.ObjectIDFromHex(stringID)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: proyecto}}

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

func DeleteProyect(stringID string) bool {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coleccion := config.InstanceDB.DB.Collection("proyects")
	id, _ := primitive.ObjectIDFromHex(stringID)
	filter := bson.D{{Key: "_id", Value: id}}

	err := coleccion.FindOneAndDelete(ctx, filter)
	if err.Err() != nil {
		return false
	}
	return true

}

func addUsers() {

}

func deleteUsers() {
	//Usar $pullAll
}
