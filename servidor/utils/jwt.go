package utils

import (
	"fmt"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt"
)

var SecretKey = []byte("mysuperSecret")

func GenerateJWT(email string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = email
	claims["exp"] = time.Now().Add(time.Minute * 4).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		fmt.Println(err)
	}
	return tokenString
}

func ValidateToken(bearerToken string) bool {
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
	return token.Valid
}
