package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (app *Application) ReadJSON(r *http.Request, input any) error {
	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	err := dec.Decode(input)
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) createToken(mail string) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": mail,                                 // Subject (user identifier)
		"iss": "crowdfun-app",                       // Issuer
		"exp": time.Now().Add(3 * time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                    // Issued at
	})

	tokenString, err := claims.SignedString(secretKeyJwt)
	if err != nil {
		return "", err
	}

	// Print information about the created token
	return tokenString, nil
}

func (app *Application) WriteJSON(w http.ResponseWriter, data any) error {
	type envelope map[string]any

	res := envelope{
		"Response:": data,
	}

	js, err := json.Marshal(res)

	if err != nil {
		return err
	}

	w.Write(js)
	return nil
}