package models

import (
	"bytes"
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
	descifrado.Users = []string{"a", "juanito@gmail.com"}
	return descifrado
}

func CifrarProyecto(proyect Proyect) ProyectCipher {
	proyectCipher := ProyectCipher{
		Cipherdata: "Proyecto Cifrado 2",
	}
	return proyectCipher
}

func CreateProyect(newProyect Proyect) bool {
	//Ciframos el proyecto
	proyectCipher := CifrarProyecto(newProyect)

	//Enviamos el proyecto cifrado al servidor
	proyectJSON, err := json.Marshal(proyectCipher)
	if err != nil {
		fmt.Println(err)
	}

	url := config.URLbase + "proyect"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(proyectJSON))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		fmt.Println("El proyecto no pudo ser creado")
		return false
	} else {
		//Si el proyecto se crea con exito nos devuelve el ID del proyecto creado
		var proyectID string
		json.NewDecoder(resp.Body).Decode(&proyectID)

		//Creamos la relacion para el usuario que crea el proyecto y para cada uno de los usuarios del campo user
		for i := 0; i < len(newProyect.Users); i++ {
			CreateRelation(newProyect.Users[i], proyectID, "")
		}
		return CreateRelation(UserSesion.Email, proyectID, "")
	}
}

func DeleteProyect(proyectID string, listsIDs []string) bool {

	url := config.URLbase + "proyects/" + proyectID

	listsJSON, err := json.Marshal(listsIDs)
	if err != nil {
		fmt.Println(err)
	}

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(listsJSON))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		fmt.Println("El proyecto no pudo ser borrado")
		return false
	} else {
		return true
	}
}

func GetProyectsByIDs() {

}
