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

type List struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Cipherdata  []byte             `bson:"cipherdata,omitempty"`
	Users       []string           `bson:"users,omitempty"`
	ProyectID   primitive.ObjectID `bson:"proyectID,omitempty"`
	Check       string             `bson:"check,omitempty"`
	UpdateCheck string             `bson:"updateCheck,omitempty"`
}

//Crear una lista
func CreateList(list List) string {
	//Creo un contexto
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//Obtengo la coleccion
	coleccion := config.InstanceDB.DB.Collection("lists")
	resultID, err := coleccion.InsertOne(ctx, list)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	stringID := resultID.InsertedID.(primitive.ObjectID).Hex()
	return stringID
}

//Borra una lista dado
func DeleteList(idString string) bool {
	//1.Borrar todas las tareas cuyo listID sea el id pasado por parametro
	DeleteTasksByListID(idString)
	//2.Borrar la lista en si
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coleccion := config.InstanceDB.DB.Collection("lists")
	id, _ := primitive.ObjectIDFromHex(idString)
	filter := bson.D{{Key: "_id", Value: id}}
	err := coleccion.FindOneAndDelete(ctx, filter)
	return err.Err() == nil
}

//Recupero todas las listas con los ids pasados como parametro
func GetListsByIDs(stringsIDs []string) []Proyect {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coleccion := config.InstanceDB.DB.Collection("lists")
	//Paso todos los ids de string a ObjectID
	var ids []primitive.ObjectID
	for i := 0; i < len(stringsIDs); i++ {
		id, _ := primitive.ObjectIDFromHex(stringsIDs[i])
		ids = append(ids, id)
	}
	var lists []Proyect
	filter := bson.M{"_id": bson.M{"$in": ids}}
	result, err := coleccion.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	err = result.All(ctx, &lists)
	if err != nil {
		log.Fatal(err)
	}
	return lists
}

//Recupero los usuarios de una lista
func GetUsersList(idString string) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coleccion := config.InstanceDB.DB.Collection("lists")
	id, _ := primitive.ObjectIDFromHex(idString)
	var proyecto Proyect
	err := coleccion.FindOne(ctx, bson.M{"_id": id}).Decode(&proyecto)
	if err != nil {
		log.Println(err)
	}
	return proyecto.Users
}

//Recupero los ids de las listas cuyo proyectID es el pasado
func GetListsByProyect(proyectID string) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coleccion := config.InstanceDB.DB.Collection("lists")
	id, _ := primitive.ObjectIDFromHex(proyectID)
	result, err := coleccion.Find(ctx, bson.M{"proyectID": id})
	if err != nil {
		log.Println(err)
	}
	var listas []List
	err = result.All(ctx, &listas)
	if err != nil {
		log.Fatal(err)
	}
	//Lo paso a una array de strings
	var listaStrings []string
	for i := 0; i < len(listas); i++ {
		listaStrings = append(listaStrings, listas[i].ID.Hex())
	}
	return listaStrings
}

func CheckModificationList(listID string, Updatecheck string) bool {
	list := GetList(listID)
	if list.ID.IsZero() || list.Check != Updatecheck {
		return true
	}
	return false
}

//Actualizar una lista
func UpdateList(list List, idString string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coleccion := config.InstanceDB.DB.Collection("lists")

	if CheckModificationList(idString, list.UpdateCheck) {
		return "Ya modificada"
	}
	list.UpdateCheck = ""
	id, _ := primitive.ObjectIDFromHex(idString)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.M{"$set": list}

	var updatedDoc bson.D
	err := coleccion.FindOneAndUpdate(ctx, filter, update).Decode(&updatedDoc)
	if err != nil {
		log.Println(err)
	}
	if len(updatedDoc) == 0 {
		return "Error"
	}
	return "OK"
}

//Añado un usuario al array Users de la lista
func AddUserList(stringID string, user string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coleccion := config.InstanceDB.DB.Collection("lists")
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

//Elimino un usuario del array Users de la lista
func DeleteUserList(listStringID string, user string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coleccion := config.InstanceDB.DB.Collection("lists")
	id, _ := primitive.ObjectIDFromHex(listStringID)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"users": user}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var listUpdated List
	err := coleccion.FindOneAndUpdate(ctx, filter, update, opts).Decode(&listUpdated)
	if err != nil {
		log.Println(err)
	}
	//Si ya no quedan mas usuarios en la lista la borro
	if len(listUpdated.Users) == 0 {
		DeleteListRelation("admin", listUpdated.ProyectID.Hex(), listUpdated.ID.Hex())
		DeleteList(listStringID)
	}
	return true
}

//Recupero una lista por su ID
func GetList(idString string) List {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coleccion := config.InstanceDB.DB.Collection("lists")
	id, _ := primitive.ObjectIDFromHex(idString)
	var list List
	err := coleccion.FindOne(ctx, bson.M{"_id": id}).Decode(&list)
	if err != nil {
		log.Println(err)
	}
	return list
}

/*
//Borro dado el proyectID todas
func DeleteListByProyectID(proyectID string) bool {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	coleccion := config.InstanceDB.DB.Collection("lists")
	id, _ := primitive.ObjectIDFromHex(proyectID)
	filter := bson.D{{Key: "proyectID", Value: id}}
	_, err := coleccion.DeleteMany(ctx, filter)
	return err == nil
}
*/
