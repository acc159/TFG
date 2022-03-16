package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID               primitive.ObjectID `json:"_id,omitempty"`
	Nombre           string             `json:"nombre"`
	Descripcion      string             `json:"apellidos"`
	FechaLimite      string             `json:"fecha_limite"`
	ArchivosAdjuntos string             `json:"archivos_adjuntos"`
	EnlacesAdjuntos  string             `json:"enlaces_adjuntos"`
}
