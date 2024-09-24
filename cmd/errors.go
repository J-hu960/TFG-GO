package main

import (
	"net/http"
)

func (app *Application) serveResponseError(w http.ResponseWriter, statusCode int, response any) {
	w.WriteHeader(statusCode)
	app.WriteJSON(w, response, "error")
}
