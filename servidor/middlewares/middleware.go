package middlewares

import (
	"log"
	"net/http"
	"servidor/utils"
)

func MiddlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Petición -> IP: %s Método: %s Ruta: %s\n", r.RemoteAddr, r.Method, r.URL)
			//Guardo el log en la base de datos
			entrada := "Petición -> IP: " + r.RemoteAddr + " Método: " + r.Method + " Ruta: " + r.URL.String()
			go utils.SaveLogs(entrada)
			next.ServeHTTP(w, r)
		})
}

func MiddlewareAddJsonHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
}

func ValidateTokenMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.String() == "/login" || r.URL.String() == "/signup" {
				next.ServeHTTP(w, r)
			} else {
				bearerToken := r.Header.Get("Authorization")
				if bearerToken != "" {
					valido, userToken := utils.ValidateToken(bearerToken)
					if valido {
						r.Header.Set("UserToken", userToken)
						next.ServeHTTP(w, r)
					} else {
						w.WriteHeader(http.StatusUnauthorized)
						w.Write([]byte("No autorizado"))
					}
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Falta de Token"))
				}
			}
		})
}
