package utils

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

func FindAndDelete(data []string, delete string) []string {
	var respuesta []string
	for i := 0; i < len(data); i++ {
		if data[i] != delete {
			respuesta = append(respuesta, data[i])
		}
	}
	return respuesta
}

func GetClientHTTPS() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return client
}

func CheckExpirationTimeToken(tokenString string) bool {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		log.Fatal(err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println(ok)
	}
	var tm time.Time
	switch iat := claims["exp"].(type) {
	case float64:
		tm = time.Unix(int64(iat), 0)
	}
	now := time.Now().Add(time.Second * 15)
	result := now.Before(tm)
	return result
}

func ToBase64FromByte(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func ToByteFromBase64(dataString string) []byte {
	data, _ := base64.StdEncoding.DecodeString(dataString)
	return data
}
