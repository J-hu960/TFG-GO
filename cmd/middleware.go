package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (app *Application) requireTokenAuth(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeader) <= 0 || authHeader[0] != "Bearer" {
			app.serveResponseError(w, http.StatusNetworkAuthenticationRequired, "No token provided or bad token")
			return
		}

		tokenString := authHeader[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validar que se está utilizando el algoritmo HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secretKeyJwt, nil
		})

		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return

		}

		// Puedes acceder a los claims como "sub", "exp", etc.
		mail := claims["sub"].(string)

		// Añadir la información del usuario al contexto si es necesario
		// Para pasar esta info al siguiente handler
		r = r.WithContext(context.WithValue(r.Context(), "mail", mail))

		// Continuar con el siguiente handler
		next.ServeHTTP(w, r)

	})

}
