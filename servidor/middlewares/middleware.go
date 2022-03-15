package middlewares

import (
	"log"
	"net/http"
)

func MiddlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Petición -> IP: %s Método: %s Ruta: %s\n", r.RemoteAddr, r.Method, r.URL)
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
