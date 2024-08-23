package model

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type ContextKey string

const tokenContextKey ContextKey = "jwtToken"

func CreateJWT(secret string, userId int) string {
	var jwtSecret = []byte(secret)
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["user_id"] = userId

	// Firmar el token con la clave secreta
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Fatal(err)
	}

	return tokenString
}

func AuthMiddleware(next http.Handler, JWTSecret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var jwtSecret = []byte(JWTSecret)
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Token no proporcionado", http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Método de firma inesperado: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Println("Usuario autenticado")
			ctx := context.WithValue(r.Context(), tokenContextKey, token)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
		}
	})
}

func ObtenerIdJWT(ctx context.Context) (float64, bool) {
	token, ok := ctx.Value(tokenContextKey).(*jwt.Token)
	if !ok {
		return 0, false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false
	}

	// Aquí puedes acceder a los datos del JWT
	userID := claims["user_id"].(float64)

	return userID, true
}
