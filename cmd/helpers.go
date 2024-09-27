package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"jordi.tfg.rewrite/internal/validator"
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

func (app *Application) WriteJSON(w http.ResponseWriter, data any, prefix string) error {
	type envelope map[string]any

	res := envelope{
		prefix: data,
	}

	js, err := json.Marshal(res)

	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")

	w.Write(js)
	return nil
}

func (app *Application) ReadUrlId(ps httprouter.Params) (int64, error) {

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil || id <= 0 {
		return 0, err
	}

	return int64(id), nil

}

func (app *Application) adaptHttprouterHandle(h httprouter.Handle) http.HandlerFunc { //recibe una funcion con ps httprouter y retorna un handlerFunc normal(es lo que necesita nuestro middleware)
	//dentro de la nueva funcion, llamamos a la funcion original
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		h(w, r, params)
	}
}

func (app *Application) readString(qs url.Values, key string, defaultvalue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultvalue
	}
	return s
}

func (app *Application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	// Si la cadena CSV está vacía, devuelve el valor por defecto
	if csv == "" {
		return defaultValue
	}

	values := strings.Split(csv, ",")

	// Elimina elementos vacíos del slice (en caso de que existan)
	nonEmptyValues := []string{}
	for _, value := range values {
		if value != "" {
			nonEmptyValues = append(nonEmptyValues, value)
		}
	}

	// Si no hay valores no vacíos, devuelve el valor por defecto
	if len(nonEmptyValues) == 0 {
		return defaultValue
	}

	return nonEmptyValues
}

func (app *Application) readInt(qs url.Values, key string, defaultValue int, v validator.Validator) int {
	s := qs.Get(key)

	// Si no hay valor (cadena vacía), devuelve el valor por defecto
	if s == "" {
		return defaultValue
	}

	// Intentar convertir la cadena a un entero
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer")
		return defaultValue
	}

	return i
}
