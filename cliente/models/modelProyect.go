package models

import (
	"bytes"
	"cliente/config"
	"cliente/utils"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Proyect struct {
	ID          string     `bson:"_id,omitempty"`
	Name        string     `bson:"name,omitempty"`
	Description string     `bson:"description,omitempty"`
	Users       []UserRole `bson:"users,omitempty"`
	Check       string     `bson:"check"`
	Rol         string     `bson:"rol"`
}

type ProyectCipher struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Cipherdata  []byte             `bson:"cipherdata,omitempty"`
	Users       []UserRole         `bson:"users,omitempty"`
	Check       string             `bson:"check"`
	UpdateCheck string             `bson:"updateCheck"`
}

type UserRole struct {
	User string `bson:"user,omitempty"`
	Rol  string `bson:"rol,omitempty"`
}

//Recupero un proyecto dado su ID
func GetProyect(proyectID string) ProyectCipher {
	url := config.URLbase + "proyects/" + proyectID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req = AddTokenHeader(req)
	client := utils.GetClientHTTPS()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	UserSesion.Token = resp.Header.Get("refreshToken")
	if resp.StatusCode == 404 {
		fmt.Println("Proyecto no encontrado")
		return ProyectCipher{}
	} else {
		var proyect ProyectCipher
		json.NewDecoder(resp.Body).Decode(&proyect)
		return proyect
	}
}

//Creo un proyecto
func CreateProyect(newProyect Proyect) (bool, bool) {
	//Añadimos el email del usuario que esta creando el proyecto

	UserRole := UserRole{
		User: UserSesion.Email,
		Rol:  "Admin",
	}

	//newProyect.Users = append(newProyect.Users, UserSesion.Email)

	newProyect.Users = append(newProyect.Users, UserRole)

	//Generamos la clave aleatoria que se utilizara en el cifrado AES
	Krandom, IVrandom := utils.GenerateKeyIV()
	//Ciframos el proyecto
	proyectCipher := CifrarProyecto(newProyect, Krandom, IVrandom)
	h := sha1.New()
	h.Write(proyectCipher.Cipherdata)
	proyectCipher.Check = hex.EncodeToString(h.Sum(nil))
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
	req = AddTokenHeader(req)
	client := utils.GetClientHTTPS()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	UserSesion.Token = resp.Header.Get("refreshToken")
	switch resp.StatusCode {
	case 400:
		fmt.Println("El proyecto no pudo ser creado")
		return false, false
	case 401:
		fmt.Println("Token Expirado")
		return false, true
	default:
		var proyectID string
		json.NewDecoder(resp.Body).Decode(&proyectID)
		//Creamos la relacion para el usuario que crea el proyecto y para cada uno de los usuarios del campo user
		CreateProyectRelations(proyectID, Krandom, newProyect.Users)
		return true, false
	}
}

//Le paso el ID del proyecto junto a su clave de cifrado y creo relaciones Usuario-Proyecto para cada usuario pasado
func CreateProyectRelations(proyectID string, Krandom []byte, users []UserRole) {
	for i := 0; i < len(users); i++ {
		publicKeyUser, _ := GetPublicKey(users[i].User)
		KrandomCipher := utils.EncryptKeyWithPublicKey(publicKeyUser, Krandom)
		CreateRelation(users[i].User, proyectID, KrandomCipher)
	}
}

//Eliminar un proyecto
func DeleteProyect(proyectID string) (bool, bool) {
	url := config.URLbase + "proyects/" + proyectID
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req = AddTokenHeader(req)
	client := utils.GetClientHTTPS()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	UserSesion.Token = resp.Header.Get("refreshToken")
	switch resp.StatusCode {
	case 400:
		fmt.Println("El proyecto no pudo ser borrado")
		return false, false
	case 401:
		fmt.Println("Token Expirado")
		return false, true
	default:
		return true, false
	}
}

//Recuperar los usuarios de un proyecto
func GetUsersProyect(proyectID string) ([]UserRole, bool) {
	url := config.URLbase + "proyect/users/" + proyectID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req = AddTokenHeader(req)
	client := utils.GetClientHTTPS()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	UserSesion.Token = resp.Header.Get("refreshToken")
	var responseObject []UserRole
	switch resp.StatusCode {
	case 400:
		fmt.Println("El usuario no pudo ser eliminado del proyecto")
		return responseObject, false
	case 401:
		fmt.Println("Token Expirado")
		return responseObject, true
	default:
		json.NewDecoder(resp.Body).Decode(&responseObject)
		return responseObject, false
	}
}

//Elimino al usuario del array Users del proyecto
func DeleteUserProyect(proyectID string, userEmail string) (bool, bool) {
	//Recupero la relacion para quitar tambien al usuario de las listas del proyecto donde este
	relation, tokenExpire := GetRelationUserProyect(userEmail, proyectID)
	if tokenExpire {
		return false, true
	} else {
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
		req = AddTokenHeader(req)
		client := utils.GetClientHTTPS()
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		UserSesion.Token = resp.Header.Get("refreshToken")
		switch resp.StatusCode {
		case 400:
			fmt.Println("El usuario no pudo ser eliminado del proyecto")
			return false, false
		case 401:
			fmt.Println("Token Expirado")
			return false, true
		default:
			//Actualizo el proyecto en local
			for i := 0; i < len(DatosUsuario); i++ {
				if DatosUsuario[i].Proyecto.ID == proyectID {
					DatosUsuario[i].Proyecto.Users = FindAndDeleteUsers(DatosUsuario[i].Proyecto.Users, userEmail)
				}
			}
			return true, false
		}
	}
}

//Devuelve la clave del proyecto descifrada
func GetProyectKey(proyectIDstring string, userEmail string) []byte {
	relation, _ := GetRelationUserProyect(UserSesion.Email, proyectIDstring)
	ProyectKeyCipher := relation.ProyectKey
	privateKey := GetPrivateKeyUser()
	proyectKey := utils.DescifrarRSA(privateKey, ProyectKeyCipher)
	return proyectKey
}

//Añadir un usuario a un proyecto
func AddUserProyect(proyectIDstring string, userEmail string) (bool, bool) {
	//Recupero la clave publica que usare para cifrar la clave del proyecto
	publicKey, _ := GetPublicKey(userEmail)
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
	req = AddTokenHeader(req)
	client := utils.GetClientHTTPS()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	UserSesion.Token = resp.Header.Get("refreshToken")
	switch resp.StatusCode {
	case 400:
		fmt.Println("El usuario no pudo ser añadido al proyecto")
		return false, false
	case 401:
		fmt.Println("Token Expirado")
		return false, true
	default:
		// //Actualizo el proyecto en local
		// for i := 0; i < len(DatosUsuario); i++ {
		// 	if DatosUsuario[i].Proyecto.ID == proyectIDstring {
		// 		DatosUsuario[i].Proyecto.Users = append(DatosUsuario[i].Proyecto.Users, userEmail)
		// 	}
		// }
		return true, false
	}
}

func GetEmailsNotInProyect(proyect Proyect) []string {
	usersAll, _ := GetUsers()
	var emails []string
	for i := 0; i < len(usersAll); i++ {
		emails = append(emails, usersAll[i].Email)
	}
	//Elimino a los usuarios que ya pertenecen al proyecto
	emailsProyect := proyect.Users
	for i := 0; i < len(emailsProyect); i++ {
		emails = utils.FindAndDelete(emails, emailsProyect[i].User)
	}
	return emails
}

//Cifrado y Descifrado
func DescifrarProyecto(proyectCipher ProyectCipher, key []byte) Proyect {
	descifradoBytes := utils.DescifrarAES(key, proyectCipher.Cipherdata)
	proyect := BytesToProyect(descifradoBytes)
	proyect.ID = proyectCipher.ID.Hex()
	proyect.Users = proyectCipher.Users
	proyect.Check = proyectCipher.Check
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

func UpdateProyect(newProyect Proyect) string {
	//Recupero la relacion del proyecto para obtener la Key de cifrado
	relation, _ := GetRelationUserProyect(UserSesion.Email, newProyect.ID)
	ProyectKeyCipher := relation.ProyectKey
	//Descifro la clave con mi clave privada
	privateKey := GetPrivateKeyUser()
	proyectKey := utils.DescifrarRSA(privateKey, ProyectKeyCipher)
	//Genero un nuevo IV
	_, IV := utils.GenerateKeyIV()
	//Cifro el nuevo proyecto y me quedo con la parte de los datos cifrados
	proyectCipher := CifrarProyecto(newProyect, proyectKey, IV)
	proyectCipher.ID, _ = primitive.ObjectIDFromHex(newProyect.ID)

	//En updateCheck pongo el hash de los datos anteriores
	proyectCipher.UpdateCheck = newProyect.Check
	h := sha1.New()
	h.Write(proyectCipher.Cipherdata)
	proyectCipher.Check = hex.EncodeToString(h.Sum(nil))

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
	req = AddTokenHeader(req)
	client := utils.GetClientHTTPS()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	UserSesion.Token = resp.Header.Get("refreshToken")
	switch resp.StatusCode {
	case 400:
		return "Error"
	case 470:
		return "Ya actualizada"
	default:
		return "OK"
	}
}

func ExistProyect(proyectID string) bool {
	relation, _ := GetRelationUserProyect(UserSesion.Email, proyectID)
	return !relation.ProyectID.IsZero()
}

func CheckUserOnProyect(proyectID string, userEmail string) bool {
	relation, _ := GetRelationUserProyect(userEmail, proyectID)
	return relation.ProyectID.IsZero()
}

func FindAndDeleteUsers(data []UserRole, delete string) []UserRole {
	var respuesta []UserRole
	for i := 0; i < len(data); i++ {
		if data[i].User != delete {
			respuesta = append(respuesta, data[i])
		}
	}
	return respuesta
}
