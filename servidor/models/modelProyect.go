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

func CreateProyect(proyecto *Proyect) {
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
		id, _ := primitive.ObjectIDFromHex(stringObjectID)
		proyecto.ID = id
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

func DeleteProyect(proyectIDstring string, listsIDs []string) bool {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//1. Borrar las listas del proyecto con las tareas asociadas a esas listas
	for i := 0; i < len(listsIDs); i++ {
		DeleteList(listsIDs[i])
	}

	//2. Borrar las relaciones en las que aparezca el proyecto
	DeleteRelationByProyectID(proyectIDstring)

	//3. Borrar el proyecto en si
	coleccion := config.InstanceDB.DB.Collection("proyects")
	id, _ := primitive.ObjectIDFromHex(proyectIDstring)
	filter := bson.D{{Key: "_id", Value: id}}

	err := coleccion.FindOneAndDelete(ctx, filter)
	if err.Err() != nil {
		return false
	}
	return true

}

func AddUserProyect(stringID string, user string) bool {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coleccion := config.InstanceDB.DB.Collection("proyects")
	id, _ := primitive.ObjectIDFromHex(stringID)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.M{"$push": bson.M{"users": user}}

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

func DeleteUserProyect(stringID string, user string) bool {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coleccion := config.InstanceDB.DB.Collection("proyects")
	id, _ := primitive.ObjectIDFromHex(stringID)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"users": user}}

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

func GetProyectsByIDs(stringsIDs []string) []Proyect {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coleccion := config.InstanceDB.DB.Collection("proyects")

	var ids []primitive.ObjectID

	for i := 0; i < len(stringsIDs); i++ {
		id, _ := primitive.ObjectIDFromHex(stringsIDs[i])
		ids = append(ids, id)
	}

	var proyects []Proyect
	filter := bson.M{"_id": bson.M{"$in": ids}}

	result, err := coleccion.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	err = result.All(ctx, &proyects)
	if err != nil {
		log.Fatal(err)
	}

	return proyects

}
