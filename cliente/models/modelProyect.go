package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProyectCipher struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Cipherdata string             `bson:"cipherdata,omitempty"`
	Users      []string           `bson:"users,omitempty"`
}

func GetProyect() {

}

func DeleteProyect() {

}
