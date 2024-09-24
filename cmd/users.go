package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"jordi.tfg.rewrite/internal/data"
	"jordi.tfg.rewrite/internal/validator"
)

func (app *Application) registerUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Phone    string `json:"Phone"`
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
	v.Check(len(input.Phone) == 9, "phone", "Provide a valid phone number please.")
	v.Check(len(input.Password) >= 6 && len(input.Password) <= 25, "password", "Must be between 6 and 25 characters.")
	v.Check(strings.Contains(input.Email, "@"), "email", "please provide a valid email.")

	if !v.Valid() {
		errors, err := json.Marshal(v.Errors)
		if err != nil {
			w.Write([]byte("Error with json errors"))
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
		w.Write([]byte("Error generating hashed password"))
	}

	newUser := data.User{
		Pk_User:  1,
		Email:    input.Email,
		Phone:    input.Phone,
		Password: password,
	}

	//TODO: Create hashedPassword && insert to BBDD

	err = app.Models.Users.Insert(&newUser)
	if err != nil {
		w.Write([]byte("Error inserting user into db..."))
		fmt.Print(err)
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

	err = app.WriteJSON(w, data)
	if err != nil {
		w.Write([]byte("Error sending back JSON response"))
	}
}
