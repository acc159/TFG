package models

import (
	"context"
	"fmt"
	"log"
	"servidor/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Proyect struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Cipherdata []byte             `bson:"cipherdata,omitempty"`
	Users      []string           `bson:"users,omitempty"`
}

//Creo un proyecto
func CreateProyect(proyecto *Proyect) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
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

//Elimino un proyecto, ademas de las listas asociadas a este y las relaciones de este
func DeleteProyect(proyectIDstring string) bool {

	//1. Borrar las listas y sus tareas
	listsIDStrings := GetListsByProyect(proyectIDstring)
	for i := 0; i < len(listsIDStrings); i++ {
		DeleteList(listsIDStrings[i])
	}

	//2. Borrar las relaciones en las que aparezca el proyecto
	DeleteRelationByProyectID(proyectIDstring)

	//3. Borrar el proyecto en si
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coleccion := config.InstanceDB.DB.Collection("proyects")
	id, _ := primitive.ObjectIDFromHex(proyectIDstring)
	filter := bson.D{{Key: "_id", Value: id}}

	err := coleccion.FindOneAndDelete(ctx, filter)
	return err.Err() == nil
}

//Recupero todos los proyectos que tengan un id de los que haya en el array de ids que paso como parametro
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

//Campo Users:

//Recupero los usuarios que tiene el proyecto cuyo id paso como parametro
func GetUsersProyect(idString string) []string {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coleccion := config.InstanceDB.DB.Collection("proyects")
	id, _ := primitive.ObjectIDFromHex(idString)
	var proyecto Proyect
	err := coleccion.FindOne(ctx, bson.M{"_id": id}).Decode(&proyecto)
	if err != nil {
		log.Println(err)
	}
	return proyecto.Users
}

//Añado un usuario al array Users del proyecto
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

//Elimino un usuario del array Users del proyecto
func DeleteUserProyect(proyectStringID string, user string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coleccion := config.InstanceDB.DB.Collection("proyects")
	id, _ := primitive.ObjectIDFromHex(proyectStringID)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"users": user}}
	//Para que me devuelva el documento actualizado
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updateProyect Proyect
	err := coleccion.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updateProyect)
	if err != nil {
		log.Println(err)
	}
	if len(updateProyect.Users) == 0 {
		DeleteProyect(proyectStringID)
	} else {
		DeleteRelation(user, proyectStringID)
	}
	return true
}

//Sin Usar

//Actualizo el proyecto
func UpdateProyect(proyecto Proyect, stringID string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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

//Recupero todos los proyectos
func GetProyects() []Proyect {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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

//Recupero un proyecto por su ID
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
