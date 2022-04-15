package utils

import (
	"crypto/tls"
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
	now := time.Now().Add(time.Minute * 2)
	result := now.Before(tm)
	return result
}
