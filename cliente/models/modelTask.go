package models

import (
	"bytes"
	"cliente/config"
	"cliente/utils"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID            string      `bson:"_id,omitempty"`
	Name          string      `bson:"name"`
	Description   string      `bson:"description"`
	Date          string      `bson:"date"`
	State         string      `bson:"state"`
	Progress      string      `bson:"progress"`
	Files         []TaskFiles `bson:"files"`
	Links         []TaskLinks `bson:"links"`
	Users         []string    `bson:"users"`
	Check         string      `bson:"check"`
	Creator       string      `bson:"creator"`
	SignCreator   SignStructEvents
	SignsReceived []SignStructEvents
	SignsClose    []SignStructEvents
}

type TaskCipher struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Cipherdata  []byte             `bson:"cipherdata,omitempty"`
	ListID      primitive.ObjectID `bson:"listID,omitempty"`
	Check       string             `bson:"check"`
	UpdateCheck string             `bson:"updateCheck"`
}

type TaskFiles struct {
	FileName string
	FileData template.URL
	UserFile string
	SignData SignStruct
}

type SignStruct struct {
	Sign     string
	UserSign string
}

type SignStructEvents struct {
	Sign     string
	UserSign string
	Message  string
}

type TaskLinks struct {
	LinkName string
	LinkUrl  string
	SignData SignStruct
	UserLink string
}

var TasksLocal []TaskCipher

var CurrentTask Task

//Recupero una tarea por su ID
func GetTask(taskID string, listID string) TaskCipher {
	url := config.URLbase + "tasks/" + taskID
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
	var taskCipher TaskCipher
	if resp.StatusCode == 400 {
		fmt.Println("No existe la tarea")
		return taskCipher
	} else {
		var taskCipher TaskCipher
		json.NewDecoder(resp.Body).Decode(&taskCipher)
		return taskCipher
	}
}

//Recupero las tareas por su ListID
func GetTasksByList(listID string) ([]Task, []TaskCipher) {
	url := config.URLbase + "tasks/list/" + listID
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
	var tasks []Task
	var tasksCipher []TaskCipher
	if resp.StatusCode == 400 {
		return tasks, tasksCipher
	} else {
		json.NewDecoder(resp.Body).Decode(&tasksCipher)
		var tasks []Task
		listKey := GetListKey(listID)
		for i := 0; i < len(tasksCipher); i++ {
			tasks = append(tasks, DescifrarTarea(tasksCipher[i], listKey))
		}
		return tasks, tasksCipher
	}
}

//Para añadir el usuario que a firmado los documentos de la tarea
func AddUserSign(task Task) Task {
	for i := 0; i < len(task.Files); i++ {
		if len(task.Files[i].SignData.Sign) > 0 {
			task.Files[i].SignData.UserSign = UserSesion.Email
		}
	}
	for i := 0; i < len(task.Links); i++ {
		if len(task.Links[i].SignData.Sign) > 0 {
			task.Links[i].SignData.UserSign = UserSesion.Email
		}
	}
	return task
}

//Creo una nueva tarea en el servidor para la lista dada
func CreateTask(stringListID string, task Task) bool {
	task = AddUserSign(task)
	listID, _ := primitive.ObjectIDFromHex(stringListID)
	//Recupero la clave de cifrado de la lista correspondiente
	listKey := GetListKey(stringListID)
	//Cifro la tarea
	taskCipher := CifrarTarea(task, listKey)
	taskCipher.ListID = listID

	h := sha1.New()
	h.Write(taskCipher.Cipherdata)
	taskCipher.Check = hex.EncodeToString(h.Sum(nil))

	//Pasamos el tipo Relation a JSON
	taskJSON, err := json.Marshal(taskCipher)
	if err != nil {
		fmt.Println(err)
	}
	url := config.URLbase + "task"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(taskJSON))
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
	if resp.StatusCode == 400 {
		fmt.Println("La tarea no pudo ser creada")
		return false
	} else {
		var newTaskID string
		json.NewDecoder(resp.Body).Decode(&newTaskID)
		fmt.Println(newTaskID)
		return true
	}
}

//Actualizar una tarea dado su ID y el de la lista
func UpdateTask(listIDstring string, task Task) string {
	listKey := GetListKey(listIDstring)
	taskCipher := CifrarTarea(task, listKey)
	taskID, _ := primitive.ObjectIDFromHex(task.ID)
	taskCipher.ID = taskID
	listID, _ := primitive.ObjectIDFromHex(listIDstring)
	taskCipher.ListID = listID

	//En updateCheck pongo el hash de los datos anteriores
	taskCipher.UpdateCheck = task.Check
	h := sha1.New()
	h.Write(taskCipher.Cipherdata)
	taskCipher.Check = hex.EncodeToString(h.Sum(nil))

	taskJSON, err := json.Marshal(taskCipher)
	if err != nil {
		fmt.Println(err)
	}
	url := config.URLbase + "tasks/" + task.ID
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(taskJSON))
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
		fmt.Println("La tarea no pudo ser borrada")
		return "Error"
	case 470:
		return "Ya actualizada"
	default:
		return "OK"
	}

}

//Borrar una tarea
func DeleteTask(taskID string) (bool, bool) {
	url := config.URLbase + "tasks/" + taskID
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
		fmt.Println("La tarea no pudo ser borrada")
		return false, false
	case 401:
		fmt.Println("Token Expirado")
		return false, true
	default:
		var resultado string
		json.NewDecoder(resp.Body).Decode(&resultado)
		return true, false
	}
}

//Borrar todas las tareas de una lista
func DeleteTaskByListID(listID string) (bool, bool) {
	url := config.URLbase + "tasks/list/" + listID
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
		fmt.Println("Las tareas no pueden ser borradas")
		return false, false
	case 401:
		fmt.Println("Token Expirado")
		return false, true
	default:
		var resultado string
		json.NewDecoder(resp.Body).Decode(&resultado)
		return true, false
	}
}

//Descifro la tarea con la clave de la lista de la tarea
func DescifrarTarea(taskCipher TaskCipher, key []byte) Task {
	descifradoBytes := utils.DescifrarAES(key, taskCipher.Cipherdata)
	task := BytesToTask(descifradoBytes)
	task.ID = taskCipher.ID.Hex()
	task.Check = taskCipher.Check
	return task
}

//Cifro una tarea con la clave de cifrado de la lista generando un nuevo IV para la tarea
func CifrarTarea(task Task, listKey []byte) TaskCipher {
	//Genero un nuevo IV
	_, IV := utils.GenerateKeyIV()
	//Cifro la tarea
	taskCipherBytes := utils.CifrarAES(listKey, IV, TaskToBytes(task))
	//Asigno la parte cifrada a la tarea
	taskCipher := TaskCipher{
		Cipherdata: taskCipherBytes,
	}
	return taskCipher
}

func DeleteUserTask(task *Task, user string) {
	newUsers := utils.FindAndDelete(task.Users, user)
	task.Users = newUsers

}

//Paso de Tarea a []Bytes
func TaskToBytes(task Task) []byte {
	taskBytes, _ := json.Marshal(task)
	return taskBytes
}

//Paso de []Bytes a Tarea
func BytesToTask(datos []byte) Task {
	var task Task
	err := json.Unmarshal(datos, &task)
	if err != nil {
		fmt.Println("error:", err)
	}
	return task
}

func CheckTaskChanges(listID string) bool {
	_, tasksCipher := GetTasksByList(listID)

	if len(tasksCipher) != len(TasksLocal) {
		return true
	} else {
		for i := 0; i < len(tasksCipher); i++ {
			if !bytes.Equal(tasksCipher[i].Cipherdata, TasksLocal[i].Cipherdata) {
				return true
			}
		}
	}
	return false
}

func GetSignFile(taskID string, listID string, filename string) string {
	taskCipher := GetTask(taskID, listID)
	task := DescifrarTarea(taskCipher, GetListKey(listID))
	for i := 0; i < len(task.Files); i++ {
		if task.Files[i].FileName == filename {
			return task.Files[i].SignData.Sign
		}
	}
	return ""
}

func GetLinkFile(taskID string, listID string, linkName string) string {
	taskCipher := GetTask(taskID, listID)
	task := DescifrarTarea(taskCipher, GetListKey(listID))
	for i := 0; i < len(task.Links); i++ {
		if task.Links[i].LinkName == linkName {
			return task.Links[i].SignData.Sign
		}
	}
	return ""
}

func SignFile(fileData string) []byte {
	signature := utils.Sign([]byte(fileData), GetPrivateKeyUser())
	return signature
}

func GetEventCreation() string {
	return CurrentTask.SignCreator.Sign
}

func GetEventReceived(userSign string) string {
	var sign string
	for i := 0; i < len(CurrentTask.SignsReceived); i++ {
		if CurrentTask.SignsReceived[i].UserSign == userSign {
			sign = CurrentTask.SignsReceived[i].Sign
		}
	}
	return sign
}

func GetEventClosed(userSign string) string {
	var sign string
	for i := 0; i < len(CurrentTask.SignsClose); i++ {
		if CurrentTask.SignsClose[i].UserSign == userSign {
			sign = CurrentTask.SignsClose[i].Sign
		}
	}
	return sign
}

func GetNumbersOfStates(tasks []Task) (int, int, int) {
	var numberPendiente int
	var numberProgreso int
	var numberFinalizada int
	for i := 0; i < len(tasks); i++ {
		if tasks[i].State == "Pendiente" {
			numberPendiente++
		} else if tasks[i].State == "En Proceso" {
			numberProgreso++
		} else {
			numberFinalizada++
		}
	}
	return numberPendiente, numberProgreso, numberFinalizada
}
