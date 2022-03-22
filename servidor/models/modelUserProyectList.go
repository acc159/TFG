package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Relation struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"userID,omitempty"`
	ProyectID  primitive.ObjectID `bson:"proyectID,omitempty"`
	ListID     primitive.ObjectID `bson:"listID,omitempty"`
	ProyectKey string             `bson:"proyectKey,omitempty"`
	ListKey    string             `bson:"listKey,omitempty"`
}

func CreateRelation() {

}

func GetProyectListbyUser(stringID string) {

}

func DeleteRelation() {

}
