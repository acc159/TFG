package models

import (
	"context"
	"log"
	"servidor/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Relation struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"userID,omitempty"`
	ProyectID  primitive.ObjectID `bson:"proyectID,omitempty"`
	ListID     primitive.ObjectID `bson:"listID,omitempty"`
	ProyectKey string             `bson:"proyectKey,omitempty"`
	ListKey    string             `bson:"listKey,omitempty"`
}

func CreateRelation(relation Relation) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coleccion := config.InstanceDB.DB.Collection("users_proyects_lists")
	_, err := coleccion.InsertOne(ctx, relation)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func GetRelations(idString string) []Relation {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coleccion := config.InstanceDB.DB.Collection("users_proyects_lists")

	userID, _ := primitive.ObjectIDFromHex(idString)
	filter := bson.D{{Key: "userID", Value: userID}}

	results, err := coleccion.Find(ctx, filter)
	if err != nil {
		log.Println(err)
	}

	var relations []Relation
	results.All(ctx, &relations)
	if err != nil {
		log.Println(err)
	}
	return relations
}

func DeleteRelation(userID primitive.ObjectID, proyectID primitive.ObjectID, listID primitive.ObjectID) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coleccion := config.InstanceDB.DB.Collection("users_proyects_lists")
	filter := bson.D{{Key: "userID", Value: userID}, {Key: "proyectID", Value: proyectID}, {Key: "listID", Value: listID}}
	err := coleccion.FindOneAndDelete(ctx, filter)
	if err.Err() != nil {
		return false
	}
	return true
}
