package utils

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt"
)

var SecretKey = []byte{}

func PrepareJWT() {
	var secret string = os.Getenv("SECRET_JWT")

	if secret == "" {
		fmt.Println("Inserta la clave para firmar los JWT")
		fmt.Scanf("%v\n", &secret)
	}

	dataHash := sha256.New()
	_, err := dataHash.Write([]byte(secret))
	if err != nil {
		panic(err)
	}
	dataHashsum := dataHash.Sum(nil)
	SecretKey = dataHashsum
}

func GenerateJWT(email string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = email
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		fmt.Println(err)
	}
	return tokenString
}

func ValidateToken(bearerToken string) (bool, string) {
	tokenString := strings.Split(bearerToken, "Bearer ")[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Inesperado metodo de firma: %v", token.Header["alg"])
		}
		return []byte(SecretKey), nil
	})
	if err != nil {
		fmt.Println(err)
	}

	var user string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user = claims["user"].(string)
	}

	return token.Valid, user
}
