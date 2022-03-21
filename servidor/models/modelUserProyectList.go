package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserProyect struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"userID"`
	ProyectID  primitive.ObjectID `bson:"proyectID"`
	ListID     primitive.ObjectID `bson:"listID"`
	ProyectKey string             `bson:"proyectKey"`
	ListKey    string             `bson:"listKey"`
}
