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
	Email      string             `bson:"email",omitempty`
	ServerKey  []byte             `bson:"server_key,omitempty"`
	PublicKey  string             `bson:"public_key,omitempty"`
	PrivateKey string             `bson:"private_key,omitempty"`
}

//Metodo para comprobar si el usuario esta vacio o tiene datos
func (u User) Empty() bool {
	var mongoId interface{}
	mongoId = u.ID
	stringObjectID := mongoId.(primitive.ObjectID).Hex()
	fmt.Println(stringObjectID)
	return (stringObjectID == "000000000000000000000000")
}

func SignUp(user User) string {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coleccion := config.InstanceDB.DB.Collection("users")

	result, err := coleccion.InsertOne(ctx, user)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	stringObjectID := result.InsertedID.(primitive.ObjectID).Hex()
	fmt.Println(stringObjectID)
	return stringObjectID

}

func Login(user User) User {
	//Creo un contexto
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("users")

	var usuario User

	//Consulto a la base de datos
	err := coleccion.FindOne(ctx, bson.M{"email": user.Email}).Decode(&usuario)
	if err != nil {
		return usuario
	}

	return usuario
}

func GetUsers() []User {
	//Creo un contexto
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("users")

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
	coleccion := config.InstanceDB.DB.Collection("users")
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
	coleccion := config.InstanceDB.DB.Collection("users")
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

func UpdateUser(idString string, usuario User) bool {
	//Creo un contexto
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("users")
	//Paso el string a un primitive.objectID
	id, _ := primitive.ObjectIDFromHex(idString)

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: usuario}}

	//Consulto a la base de datos
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

func DeleteUser(idString string) bool {
	//Creo un contexto
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("users")
	//Paso el string a un primitive.objectID
	id, _ := primitive.ObjectIDFromHex(idString)

	filter := bson.D{{"_id", id}}

	err := coleccion.FindOneAndDelete(ctx, filter)
	if err.Err() != nil {
		return false
	}
	return true
}

/*
func AddProyect(proyecto ProyectKey, idString string) {
	//Creo un contexto
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("users")

	id, _ := primitive.ObjectIDFromHex(idString)

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$push", bson.D{{"proyectos", proyecto}}}}

	coleccion.UpdateOne(ctx, filter, update)

}
*/
