package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"jordi.tfg.rewrite/internal/data"
	"jordi.tfg.rewrite/internal/validator"
)

func (app *Application) registerUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"Email"`
		Password string `json:"Password"`
	}
	err := app.ReadJSON(r, &input)

	if err != nil {
		w.Write([]byte("Error decoding json"))
		fmt.Print(err)
		return
	}

	v := validator.New()
	// v.Check(len(input.Phone) == 9, "phone", "Provide a valid phone number please.")
	v.Check(len(input.Password) >= 6 && len(input.Password) <= 25, "password", "Must be between 6 and 25 characters.")
	v.Check(validator.ValidateMail(input.Email), "email", "please provide a valid email.")

	if !v.Valid() {
		errors, err := json.Marshal(v.Errors)
		if err != nil {
			app.serveResponseError(w, http.StatusInternalServerError, v.Errors)
			return
		}
		w.Write(errors)
		return
	}
	password := data.PasswordS{
		PlainText: input.Password,
	}
	err = password.CreateHashedPassword()
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	newUser := data.User{
		Pk_User: 1,
		Email:   input.Email,
		// Phone:    input.Phone,
		Password: password,
	}

	err = app.Models.Users.Insert(&newUser)
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jwtToken, err := app.createToken(input.Email)
	if err != nil {
		w.Write([]byte("Error creating JWT"))
		return
	}

	data := map[string]any{
		"New User": newUser,
		"Token":    jwtToken,
	}

	err = app.WriteJSON(w, data, "New user")
	if err != nil {
		w.Write([]byte("Error sending back JSON response"))
	}
}

func (app *Application) logInUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"Email"`
		Password string `json:"Password"`
	}

	err := app.ReadJSON(r, &input)

	if err != nil {
		app.serveResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	v := validator.New()

	v.Check(validator.ValidateMail(input.Email), "email", "provide a valid email please")
	v.Check(len(input.Password) >= 7 && len(input.Password) <= 25, "password", "password must be between 7 and 25 chr")

	if !v.Valid() {
		app.serveResponseError(w, http.StatusBadRequest, v.Errors)
		return
	}

	//Recuperar usuari que fa solicitud
	user, err := app.Models.Users.GetByMail(input.Email)
	if err != nil {
		app.serveResponseError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	//Veure si el hash de la contraseña == a la que tenim a bbdd per l'usuari

	if !data.VerifyPassword(user.Password.HashedPasswd, input.Password) {
		app.serveResponseError(w, http.StatusBadRequest, "Bad credentials")
		return
	}

	//Si tot OK, crear token y retornarlo

	jwtToken, err := app.createToken(input.Email)
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = app.WriteJSON(w, jwtToken, "token")
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
func (app *Application) GetUserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := app.ReadUrlId(ps)
	if err != nil {
		app.serveResponseError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user, err := app.Models.Users.GetById(int64(id))
	fmt.Print(user)

	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = app.WriteJSON(w, user, "User")
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

}

func (app *Application) deleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, err := app.ReadUrlId(ps)
	if err != nil {
		app.serveResponseError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = app.Models.Users.DeleteById(id)
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
func (app *Application) updateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var input struct {
		Description string `json:"description"`
		Role        string `json:"role"`
	}

	v := validator.New()

	v.Check(len(input.Description) <= 250, "description", "cannot be not more than 250 ch")
	v.Check(input.Role == "user" || input.Role == "admin" || input.Role == "moderator", "role", "not a real role, sorry")

	if !v.Valid() {
		app.serveResponseError(w, http.StatusBadRequest, v.Errors)
		return

	}
	err := app.ReadJSON(r, &input)
	if err != nil {
		app.serveResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := app.ReadUrlId(ps)
	if err != nil {
		app.serveResponseError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user, err := app.Models.Users.GetById(id)
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if input.Description != "" {
		user.Description = input.Description
	}
	if input.Role != "" {
		user.Role = input.Role
	}

	err = app.Models.Users.Update(user)

	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return

	}

	err = app.WriteJSON(w, user, "Updated User: ")
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
