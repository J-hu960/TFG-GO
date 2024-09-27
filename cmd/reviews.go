package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"jordi.tfg.rewrite/internal/data"
)

func (app *Application) CreateReviews(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	projectId, err := app.ReadUrlId(ps)
	if err != nil {
		app.serveResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	var input struct {
		Content string `json:"content"`
	}

	err = app.ReadJSON(r, &input)
	if err != nil {
		app.serveResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	userMailContext := r.Context().Value("mail")
	userMail, ok := userMailContext.(string)
	if !ok {
		app.serveResponseError(w, http.StatusBadRequest, "Malformed JWT or missing user email in context")
		return
	}
	user, err := app.Models.Users.GetByMail(userMail)
	if err != nil {
		app.serveResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	newReview := data.Review{
		Content:    input.Content,
		Id_user:    user.Pk_User,
		Id_project: projectId,
	}

	err = app.Models.Reviews.Insert(newReview)
	if err != nil {
		app.serveResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = app.WriteJSON(w, newReview, "Review Created: ")
	if err != nil {
		app.serveResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
}

func (app *Application) DeleteReview(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := app.ReadUrlId(ps)
	if err != nil {
		app.serveResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = app.Models.Reviews.Delete(id)
	if err != nil {
		app.serveResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusAccepted)

}
