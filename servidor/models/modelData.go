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

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Cipherdata string             `bson:"cipherdata"`
}

func GetUsers() []User {
	//Creo un contexto
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("usuarios")

	//Consulto a la base de datos
	result, err := coleccion.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("ASDFDSF")
		log.Fatal(err)
	}

	//Lo paso a struct de GO la consulta
	var usuarios []User
	err = result.All(ctx, &usuarios)
	if err != nil {
		fmt.Println("ASDFDSF")
		log.Fatal(err)
	}
	return usuarios
}

func GetUser(idString string) User {
	//Creo un contexto
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("usuarios")
	//Paso el string a un primitive.objectID
	id, _ := primitive.ObjectIDFromHex(idString)

	var usuario User

	//Consulto a la base de datos
	err := coleccion.FindOne(ctx, bson.M{"_id": id}).Decode(&usuario)
	if err != nil {
		return usuario
	}
	return usuario
}

func CreateUser(usuario User) string {
	//Creo un contexto
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("usuarios")
	//Inserto el usuario pasado por parametro
	result, err := coleccion.InsertOne(ctx, usuario)
	if err != nil {
		log.Fatal(err)
		return "0"
	}

	//Paso el primitive.ObjectID a un string
	stringObjectID := result.InsertedID.(primitive.ObjectID).Hex()
	fmt.Println(stringObjectID)
	return stringObjectID
}
