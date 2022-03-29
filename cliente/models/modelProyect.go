package models

import (
	"cliente/config"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Proyect struct {
	ID          string   `bson:"_id,omitempty"`
	Name        string   `bson:"name,omitempty"`
	Description string   `bson:"description,omitempty"`
	Users       []string `bson:"users,omitempty"`
}

type ProyectCipher struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Cipherdata string             `bson:"cipherdata,omitempty"`
}

func GetProyect(proyectID string) ProyectCipher {
	resp, err := http.Get(config.URLbase + "proyects/" + proyectID)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	//Compruebo si no hay ningun usuario
	if resp.StatusCode == 404 {
		fmt.Println("Proyecto no encontrado")
		return ProyectCipher{}
	} else {
		var responseObject ProyectCipher
		json.NewDecoder(resp.Body).Decode(&responseObject)
		return responseObject
	}
}

func DescifrarProyecto(proyecto ProyectCipher) Proyect {
	var descifrado Proyect
	descifrado.ID = proyecto.ID.Hex()
	descifrado.Name = "Nombre Proyecto"
	descifrado.Description = "Esto seria la descripcion del proyecto"
	descifrado.Users = []string{"pepito@gmail.com", "juanito@gmail.com"}
	return descifrado
}

func GetProyectsByIDs() {

}

func DeleteProyect() {

}
