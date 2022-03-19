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
	Users      []string           `bson:"users"`
}

type ProyectKey struct {
	Proyecto primitive.ObjectID `bson:"proyecto,omitempty"`
	Key      string             `bson:"key"`
}

/*
func CreateIndexUnique() {
	coleccion := config.InstanceDB.DB.Collection("proyectos")

	indexName, err := coleccion.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(indexName)
	}

}
*/

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
		fmt.Println(err)
	} else {
		//Paso el primitive.ObjectID a un string
		stringObjectID := result.InsertedID.(primitive.ObjectID).Hex()
		fmt.Println(stringObjectID)
	}
}
