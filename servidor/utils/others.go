package utils

import (
	"context"
	"fmt"
	"net/http"
	"servidor/config"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func SaveLogs(entrada string) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coleccion := config.InstanceDB.DB.Collection("logs_server")

	log := bson.M{"entrada": entrada}
	_, err := coleccion.InsertOne(ctx, log)
	if err != nil {
		fmt.Println(err)
	}
}

func SetRefreshToken(w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	userToken := r.Header.Values("UserToken")[0]
	refreshToken := GenerateJWT(userToken)
	w.Header().Set("refreshToken", refreshToken)
	return w
}
