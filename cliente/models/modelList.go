package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ListCipher struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	Cipherdata string               `bson:"cipherdata,omitempty"`
	Users      []primitive.ObjectID `bson:"users,omitempty"`
}
