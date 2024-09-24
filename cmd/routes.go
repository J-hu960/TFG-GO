package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) Routes() http.Handler {
	router := httprouter.New()

	router.GET("/", HealthCheck)
	router.HandlerFunc(http.MethodPost, "/v1/sign-up", app.registerUser)
	router.HandlerFunc(http.MethodPost, "/v1/sign-in", app.logInUser)
	router.Handle(http.MethodGet, "/v1/users/:id", app.GetUserInfo)
	router.Handle(http.MethodDelete, "/v1/users/:id", app.deleteUser)
	router.Handle(http.MethodPut, "/v1/users/:id", app.updateUser)

	return router
}
