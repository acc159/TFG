package models

import (
	"context"
	"fmt"
	"log"
	"servidor/config"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Email      string             `bson:"email"`
	Cipherdata string             `bson:"cipherdata"`
	Edad       int                `bson:"edad"`
	Proyectos  []ProyectKey       `bson:"proyectos"`
}

//Metodo para comprobar si el usuario esta vacio o tiene datos
func (u User) Empty() bool {
	var mongoId interface{}
	mongoId = u.ID
	stringObjectID := mongoId.(primitive.ObjectID).Hex()
	fmt.Println(stringObjectID)
	return (stringObjectID == "000000000000000000000000" && u.Cipherdata == "")
}

func GetUsers() []User {
	//Creo un contexto
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("usuarios")

	//Consulto a la base de datos
	result, err := coleccion.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	//Lo paso a struct de GO la consulta
	var usuarios []User
	err = result.All(ctx, &usuarios)
	if err != nil {
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
		//Compruebo si el error que se me da es por valor duplicado
		textoError := err.Error()
		if duplicado := strings.Index(textoError, "duplicate"); duplicado != -1 {
			return "Duplicado"
		}
		return "0"
	}

	//Paso el primitive.ObjectID a un string
	stringObjectID := result.InsertedID.(primitive.ObjectID).Hex()
	fmt.Println(stringObjectID)
	return stringObjectID
}

func UpdateUser(idString string, campo string, valorNuevo interface{}) User {
	//Creo un contexto
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("usuarios")
	//Paso el string a un primitive.objectID
	id, _ := primitive.ObjectIDFromHex(idString)

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{campo, valorNuevo}}}}

	var usuario User
	//Consulto a la base de datos
	result, err := coleccion.UpdateOne(ctx, filter, update)
	if err != nil {
		return usuario
	}
	if result.ModifiedCount > 0 {
		return GetUser(idString)
	}
	return usuario
}

func DeleteUser(idString string) bool {
	//Creo un contexto
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("usuarios")
	//Paso el string a un primitive.objectID
	id, _ := primitive.ObjectIDFromHex(idString)

	filter := bson.D{{"_id", id}}

	result, err := coleccion.DeleteOne(ctx, filter)
	if err != nil || result.DeletedCount == 0 {
		return false
	}
	return true
}

func AddProyect(proyecto ProyectKey, idString string) {
	//Creo un contexto
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("usuarios")

	id, _ := primitive.ObjectIDFromHex(idString)

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$push", bson.D{{"proyectos", proyecto}}}}

	coleccion.UpdateOne(ctx, filter, update)

}
