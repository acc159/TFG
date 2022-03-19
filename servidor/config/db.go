package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var nombreDB = "pruebas"

type MongoConection struct {
	DB     *mongo.Database
	Client *mongo.Client
}

var InstanceDB MongoConection

func ConnectDB() {

	cadena_conexion := "mongodb://127.0.0.1:27017"

	client, err := mongo.NewClient(options.Client().ApplyURI(cadena_conexion))
	if err != nil {
		log.Fatal(err)
	}

	//Establecemos un context antes de conectarnos
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	InstanceDB = MongoConection{
		DB:     client.Database(nombreDB),
		Client: client,
	}

	CreateIndexUniqueUsers()
	CreateIndexCompose()
	CreateIndexListIDinTask()
}

func CreateIndexUniqueUsers() {
	coleccion := InstanceDB.DB.Collection("usuarios")
	_, err := coleccion.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)

	if err != nil {
		fmt.Println(err)
	}

}

func CreateIndexListIDinTask() {
	coleccion := InstanceDB.DB.Collection("tasks")
	_, err := coleccion.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{{Key: "listID", Value: 1}},
		},
	)

	if err != nil {
		fmt.Println(err)
	}

}

func CreateIndexCompose() {
	coleccion := InstanceDB.DB.Collection("Usuarios_Proyectos")
	_, err := coleccion.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{"userID", 1}, {"proyectID", 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		fmt.Println(err)
	}
}

func DisconectDB() {
	InstanceDB.Client.Disconnect(context.TODO())
}
