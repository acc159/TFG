package models

import (
	"bytes"
	"cliente/config"
	"cliente/utils"
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
	Cipherdata []byte             `bson:"cipherdata,omitempty"`
	Users      []string           `bson:"users,omitempty"`
}

//Recupero un proyecto dado su ID
func GetProyect(proyectID string) ProyectCipher {

	url := config.URLbase + "proyects/" + proyectID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := utils.GetClientHTTPS()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	//Compruebo si no hay ningun usuario
	if resp.StatusCode == 404 {
		fmt.Println("Proyecto no encontrado")
		return ProyectCipher{}
	} else {
		var proyect ProyectCipher
		json.NewDecoder(resp.Body).Decode(&proyect)
		//Descifro
		return proyect
	}
}

//Creo un proyecto
func CreateProyect(newProyect Proyect) bool {
	//Añadimos el email del usuario que esta creando el proyecto
	newProyect.Users = append(newProyect.Users, UserSesion.Email)
	//Generamos la clave aleatoria que se utilizara en el cifrado AES
	Krandom, IVrandom := utils.GenerateKeyIV()
	//Ciframos el proyecto
	proyectCipher := CifrarProyecto(newProyect, Krandom, IVrandom)
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
	client := utils.GetClientHTTPS()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		fmt.Println("El proyecto no pudo ser creado")
		return false
	}
	//Si el proyecto se crea con exito nos devuelve el ID del proyecto creado
	var proyectID string
	json.NewDecoder(resp.Body).Decode(&proyectID)
	//Creamos la relacion para el usuario que crea el proyecto y para cada uno de los usuarios del campo user
	CreateProyectRelations(proyectID, Krandom, newProyect.Users)
	return true
}

//Le paso el ID del proyecto junto a su clave de cifrado y creo relaciones Usuario-Proyecto para cada usuario pasado
func CreateProyectRelations(proyectID string, Krandom []byte, users []string) {
	for i := 0; i < len(users); i++ {
		publicKeyUser := GetPublicKey(users[i])
		KrandomCipher := utils.EncryptKeyWithPublicKey(publicKeyUser, Krandom)
		CreateRelation(users[i], proyectID, KrandomCipher)
	}
}

//Eliminar un proyecto
func DeleteProyect(proyectID string) bool {
	url := config.URLbase + "proyects/" + proyectID
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := utils.GetClientHTTPS()
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

//Recuperar los usuarios de un proyecto
func GetUsersProyect(proyectID string) []string {
	url := config.URLbase + "proyect/users/" + proyectID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := utils.GetClientHTTPS()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var responseObject []string
	if resp.StatusCode == 404 {
		return responseObject
	} else {
		json.NewDecoder(resp.Body).Decode(&responseObject)
		return responseObject
	}
}

//Elimino al usuario del array Users del proyecto
func DeleteUserProyect(proyectID string, userEmail string) bool {
	//Recupero la relacion para quitar tambien al usuario de las listas del proyecto donde este
	relation := GetRelationUserProyect(userEmail, proyectID)
	//Para cada lista que tiene el proyecto elimino al usuario de dicha lista
	for i := 0; i < len(relation.Lists); i++ {
		DeleteUserList(relation.Lists[i].ListID, userEmail)
	}
	url := config.URLbase + "proyect/users/" + proyectID + "/" + userEmail
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := utils.GetClientHTTPS()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		fmt.Println("El usuario no pudo ser eliminado del proyecto")
		return false
	} else {
		//Actualizo el proyecto en local
		for i := 0; i < len(DatosUsuario); i++ {
			if DatosUsuario[i].Proyecto.ID == proyectID {
				DatosUsuario[i].Proyecto.Users = utils.FindAndDelete(DatosUsuario[i].Proyecto.Users, userEmail)
			}
		}
		return true
	}
}

//Devuelve la clave del proyecto descifrada
func GetProyectKey(proyectIDstring string, userEmail string) []byte {
	relation := GetRelationUserProyect(UserSesion.Email, proyectIDstring)
	ProyectKeyCipher := relation.ProyectKey
	privateKey := GetPrivateKeyUser()
	proyectKey := utils.DescifrarRSA(privateKey, ProyectKeyCipher)
	return proyectKey
}

//Añadir un usuario a un proyecto
func AddUserProyect(proyectIDstring string, userEmail string) bool {
	//Recupero la clave publica que usare para cifrar la clave del proyecto
	publicKey := GetPublicKey(userEmail)
	//Recupero la clave del proyecto descifrada
	proyectKey := GetProyectKey(proyectIDstring, UserSesion.Email)
	//Cifro la clave del proyecto con la clave publica del usuario nuevo añadido
	newProyectKey := utils.CifrarRSA(publicKey, proyectKey)
	CreateRelation(userEmail, proyectIDstring, newProyectKey)
	//Añado el usuario al array de usuarios del proyecto
	url := config.URLbase + "proyect/users/" + proyectIDstring + "/" + userEmail
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := utils.GetClientHTTPS()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		fmt.Println("El usuario no pudo ser añadido al proyecto")
		return false
	} else {
		//Actualizo el proyecto en local
		for i := 0; i < len(DatosUsuario); i++ {
			if DatosUsuario[i].Proyecto.ID == proyectIDstring {
				DatosUsuario[i].Proyecto.Users = append(DatosUsuario[i].Proyecto.Users, userEmail)
			}
		}
		return true
	}
}

func GetEmailsNotInProyect(proyect Proyect) []string {
	usersAll := GetUsers()
	var emails []string
	for i := 0; i < len(usersAll); i++ {
		emails = append(emails, usersAll[i].Email)
	}
	//Elimino a los usuarios que ya pertenecen al proyecto
	emailsProyect := proyect.Users
	for i := 0; i < len(emailsProyect); i++ {
		emails = utils.FindAndDelete(emails, emailsProyect[i])
	}
	return emails
}

//Cifrado y Descifrado
func DescifrarProyecto(proyectCipher ProyectCipher, key []byte) Proyect {
	descifradoBytes := utils.DescifrarAES(key, proyectCipher.Cipherdata)
	proyect := BytesToProyect(descifradoBytes)
	proyect.ID = proyectCipher.ID.Hex()
	proyect.Users = proyectCipher.Users
	return proyect
}

func CifrarProyecto(proyect Proyect, key []byte, IV []byte) ProyectCipher {
	//Paso el proyecto a []byte
	proyectBytes := ProyectToBytes(proyect)
	//Cifro
	proyectCipherBytes := utils.CifrarAES(key, IV, proyectBytes)
	proyectCipher := ProyectCipher{
		Cipherdata: proyectCipherBytes,
		Users:      proyect.Users,
	}
	return proyectCipher
}

func ProyectToBytes(proyect Proyect) []byte {
	proyectBytes, _ := json.Marshal(proyect)
	return proyectBytes
}

func BytesToProyect(datos []byte) Proyect {
	var proyect Proyect
	err := json.Unmarshal(datos, &proyect)
	if err != nil {
		fmt.Println("error:", err)
	}
	return proyect
}

func UpdateProyect(newProyect Proyect) bool {
	//Recupero la relacion del proyecto para obtener la Key de cifrado
	relation := GetRelationUserProyect(UserSesion.Email, newProyect.ID)
	ProyectKeyCipher := relation.ProyectKey
	//Descifro la clave con mi clave privada
	privateKey := GetPrivateKeyUser()
	proyectKey := utils.DescifrarRSA(privateKey, ProyectKeyCipher)

	// //Recupero el IV del proyecto cifrado
	// oldProyect := GetProyect(newProyect.ID)
	// IV := utils.GetIV(oldProyect.Cipherdata)

	//Genero un nuevo IV
	_, IV := utils.GenerateKeyIV()
	//Cifro el nuevo proyecto y me quedo con la parte de los datos cifrados
	proyectCipher := CifrarProyecto(newProyect, proyectKey, IV)
	proyectCipher.ID, _ = primitive.ObjectIDFromHex(newProyect.ID)
	//Actualizo el proyecto en el servidor
	proyectJSON, err := json.Marshal(proyectCipher)
	if err != nil {
		fmt.Println(err)
	}
	url := config.URLbase + "proyects/" + newProyect.ID
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(proyectJSON))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := utils.GetClientHTTPS()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		fmt.Println("El proyecto no pudo ser actualizado")
		return false
	} else {
		return true
	}
}
