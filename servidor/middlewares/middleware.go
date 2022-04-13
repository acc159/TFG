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

// func MiddlewareJWT(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {

// 			if r.URL.String() == "/login" || r.URL.String() == "/signup" {
// 				next.ServeHTTP(w, r)
// 			} else {
// 				authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
// 				if len(authHeader) != 2 {
// 					w.WriteHeader(http.StatusUnauthorized)
// 					w.Write([]byte("Token mal formado"))
// 				} else {
// 					jwtToken := authHeader[1]
// 					token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
// 						if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 							fmt.Println("ALLI")
// 							return nil, fmt.Errorf("Inesperado metodo de firma: %v", token.Header["alg"])
// 						}
// 						return []byte(utils.SecretKey), nil
// 					})
// 					if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 						ctx := context.WithValue(r.Context(), "props", claims)
// 						// Access context values in handlers like this
// 						// props, _ := r.Context().Value("props").(jwt.MapClaims)
// 						next.ServeHTTP(w, r.WithContext(ctx))
// 					} else {
// 						fmt.Println(err)
// 						w.WriteHeader(http.StatusUnauthorized)
// 						w.Write([]byte("No autorizado"))
// 					}
// 				}
// 			}
// 		})
// }

func ValidateTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		if bearerToken != "" {
			if utils.ValidateToken(bearerToken) {
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
}
