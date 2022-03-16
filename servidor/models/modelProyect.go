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
	Cipherdata string             `bson:"cipherdata"`
	Users      []User             `bson:"users"`
}

func GetProyects() []Proyect {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	coleccion := config.InstanceDB.DB.Collection("proyectos")

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

func CreateProyect(proyecto Proyect) {
	//Creo un contexto
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("proyectos")
	//Inserto el usuario pasado por parametro
	result, err := coleccion.InsertOne(ctx, proyecto)
	if err != nil {
		log.Fatal(err)
	}

	//Paso el primitive.ObjectID a un string
	stringObjectID := result.InsertedID.(primitive.ObjectID).Hex()
	fmt.Println(stringObjectID)
}
