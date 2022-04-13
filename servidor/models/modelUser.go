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
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Email      string             `bson:"email,omitempty"`
	ServerKey  []byte             `bson:"server_key,omitempty"`
	PublicKey  []byte             `bson:"public_key,omitempty"`
	PrivateKey []byte             `bson:"private_key,omitempty"`
	Token      string             `bson:"token,omitempty"`
}

//Metodo para comprobar si el usuario esta vacio o tiene datos
func (u User) Empty() bool {
	return u.ID.Hex() == "000000000000000000000000"
}

//Registro del usuario devolviendo el ID del documento creado en la base de datos
func SignUp(user User) string {
	//Bcrypt
	KservidorHash, err := bcrypt.GenerateFromPassword(user.ServerKey, 12)
	if err != nil {
		fmt.Println("Error al hashear")
	}
	user.ServerKey = KservidorHash
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coleccion := config.InstanceDB.DB.Collection("users")
	result, err := coleccion.InsertOne(ctx, user)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	stringObjectID := result.InsertedID.(primitive.ObjectID).Hex()
	return stringObjectID
}

//Login del Usuario -> En este caso solo recuperamos el usuario la comprobacion sobre este la realizo en el cliente, devuelvo el usuario
func Login(userLogin User) User {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coleccion := config.InstanceDB.DB.Collection("users")
	var usuario User
	err := coleccion.FindOne(ctx, bson.M{"email": userLogin.Email}).Decode(&usuario)
	if err != nil {
		return usuario
	}
	//Comparo
	err = bcrypt.CompareHashAndPassword(usuario.ServerKey, userLogin.ServerKey)
	if err != nil {
		return User{}
	} else {
		return usuario
	}
}

//Revisar
func UpdateUser(idString string, usuario User) bool {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coleccion := config.InstanceDB.DB.Collection("users")
	id, _ := primitive.ObjectIDFromHex(idString)

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: usuario}}

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

//Borrar el usuario cuyo email es pasado por parametro
func DeleteUser(userEmail string) bool {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coleccion := config.InstanceDB.DB.Collection("users")
	filter := bson.D{{Key: "email", Value: userEmail}}
	err := coleccion.FindOneAndDelete(ctx, filter)
	return err.Err() == nil
}

//SIN USAR
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

func GetUser(email string) User {
	//Creo un contexto
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("users")
	var usuario User
	//Consulto a la base de datos
	err := coleccion.FindOne(ctx, bson.M{"email": email}).Decode(&usuario)
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
