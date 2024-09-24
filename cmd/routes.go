package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) Routes() http.Handler {
	router := httprouter.New()

	router.GET("/", HealthCheck)
	router.HandlerFunc(http.MethodPost, "/v1/sign-up", app.registerUser)

	return router
}
