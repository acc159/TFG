package config

import (
	"context"
	"fmt"
	"log"
	"os"
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
	cadena := os.Getenv("CADENA_CONEXION")

	if cadena == "" {
		fmt.Println("Inserta la cadena de conexion de la base de datos")
		fmt.Scanf("%v\n", &cadena)
	}

	//cadena_conexion := "mongodb://127.0.0.1:27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(cadena))
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
	//Creo los indices
	CreateIndexUniqueUsers()
	CreateIndexComposeUserProyectList()
	CreateIndexListIDinTask()
}

func CreateIndexUniqueUsers() {
	coleccion := InstanceDB.DB.Collection("users")
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

func CreateIndexComposeUserProyectList() {
	coleccion := InstanceDB.DB.Collection("users_proyects_lists")
	_, err := coleccion.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "userEmail", Value: 1}, {Key: "proyectID", Value: 1}},
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
