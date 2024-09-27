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
	router.Handler(http.MethodPost, "/v1/projects", app.requireTokenAuth(app.CreateProjectHandler))
	router.Handler(http.MethodGet, "/v1/projects", app.requireTokenAuth(app.ListProjects))
	router.Handler(http.MethodGet, "/v1/projects/:id", app.requireTokenAuth(app.adaptHttprouterHandle(app.ListSingleProject)))
	router.Handler(http.MethodPut, "/v1/projects/:id", app.requireTokenAuth(app.adaptHttprouterHandle(app.EditProject)))
	router.Handler(http.MethodDelete, "/v1/projects/:id", app.requireTokenAuth(app.adaptHttprouterHandle(app.DeleteProject)))
	router.Handler(http.MethodPost, "/v1/reviews/:id", app.requireTokenAuth(app.adaptHttprouterHandle(app.CreateReviews)))
	router.Handler(http.MethodDelete, "/v1/reviews/:id", app.requireTokenAuth(app.adaptHttprouterHandle(app.DeleteReview)))
	router.Handler(http.MethodPut, "/v1/projects-like/:id", app.requireTokenAuth(app.adaptHttprouterHandle(app.NewLikeHandler)))
	router.Handler(http.MethodPut, "/v1/projects-dislike/:id", app.requireTokenAuth(app.adaptHttprouterHandle(app.NewDislikeHandler)))

	return router
}
