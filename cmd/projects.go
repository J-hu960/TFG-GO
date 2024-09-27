package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"jordi.tfg.rewrite/internal/data"
	"jordi.tfg.rewrite/internal/validator"
)

func (app *Application) ListProjects(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string
		Category []string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Category = app.readCSV(qs, "category", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, *v)
	input.Filters.PageSize = app.readInt(qs, "pageSize", 10, *v)

	input.Filters.Sort = app.readString(qs, "sort", "created_at")
	input.Filters.SortSafeList = []string{"pk_project", "name", "created_at", "-pk_project", "-name", "-created_at"}

	fmt.Print("Input: ", input)
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.serveResponseError(w, http.StatusBadRequest, v.Errors)
		return
	}

	projects, err := app.Models.Projects.GetAll(input.Name, input.Category, input.Filters)
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = app.WriteJSON(w, projects, "Projects: ")
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return

	}

}

func (app *Application) CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name           string   `json:"name"`
		Photos         []string `json:"photos"`
		Link_web       string   `json:"link_web"`
		Description    string   `json:"description"`
		FoundsExpected int64    `json:"founds_expected"`
		Category       []string `json:"category"`
	}

	err := app.ReadJSON(r, &input)
	if err != nil {
		app.serveResponseError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	//validar input..

	//

	newProject := data.Project{
		Name:           input.Name,
		Photos:         input.Photos,
		Link_web:       input.Link_web,
		Description:    input.Description,
		FoundsExpected: input.FoundsExpected,
		Category:       input.Category,
	}

	userMailContext := r.Context().Value("mail")
	userMail, ok := userMailContext.(string)

	if !ok {
		app.serveResponseError(w, http.StatusBadRequest, "Mal formed JWT")
		return
	}
	user, err := app.Models.Users.GetByMail(userMail)

	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	newProject.IdCreator = user.Pk_User

	err = app.Models.Projects.Insert(newProject)

	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = app.WriteJSON(w, newProject, "New Project")

	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *Application) ListSingleProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, err := app.ReadUrlId(ps)
	if err != nil {
		app.serveResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	project, err := app.Models.Projects.GetById(id)

	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = app.WriteJSON(w, project, "Project")
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

}

func (app *Application) EditProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Name           string   `json:"name"`
		Photos         []string `json:"photos"`
		Link_web       string   `json:"link_web"`
		Description    string   `json:"description"`
		FoundsExpected int64    `json:"founds_expected"`
		Category       []string `json:"category"`
	}

	err := app.ReadJSON(r, &input)
	if err != nil {
		app.serveResponseError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	id, err := app.ReadUrlId(ps)
	if err != nil {
		app.serveResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	project, err := app.Models.Projects.GetById(id)
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//en comptes de mirar cada null value, podriem fer tots un puntero i fer un loop i mirar != nil
	if input.Name != "" {
		project.Name = input.Name
	}
	if input.Link_web != "" {
		project.Link_web = input.Link_web

	}

	if input.Description != "" {
		project.Description = input.Description
	}

	if input.FoundsExpected != 0 {
		project.FoundsExpected = input.FoundsExpected
	}

	if len(input.Category) > 0 {
		project.Category = input.Category
	}
	if len(input.Photos) > 0 {
		project.Photos = input.Photos
	}

	err = app.Models.Projects.UpdateProject(project)
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = app.WriteJSON(w, project, "Updated Project: ")

	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

}

func (app *Application) DeleteProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, err := app.ReadUrlId(ps)
	if err != nil {
		app.serveResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = app.Models.Projects.DeleteProject(id)
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

}

func (app *Application) NewLikeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	projectId, err := app.ReadUrlId(ps)
	if err != nil {
		app.serveResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	userMailContext := r.Context().Value("mail")

	userMail, ok := userMailContext.(string)
	if !ok {
		app.serveResponseError(w, http.StatusInternalServerError, "JWT mal formated")
		return
	}

	user, err := app.Models.Users.GetByMail(userMail)
	if err != nil {
		app.serveResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	hasAlreadyLiked, err := app.Models.Users.HasAlreadyLikedProject(user.Pk_User, projectId)
	print(hasAlreadyLiked)

	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !hasAlreadyLiked {

		hasAlreadyDisliked, err := app.Models.Users.HasAlreadyDislikedProject(user.Pk_User, projectId)
		if err != nil {
			app.serveResponseError(w, http.StatusInternalServerError, err.Error())
			return

		}

		if hasAlreadyDisliked {
			err = app.Models.Projects.SubDislike(projectId)
			if err != nil {
				app.serveResponseError(w, http.StatusInternalServerError, err.Error())
				return

			}

			err = app.Models.Users.DeleteUserProjectDislikeRelation(user.Pk_User, projectId)
			if err != nil {
				app.serveResponseError(w, http.StatusInternalServerError, err.Error())
				return

			}
		}

		err = app.Models.Projects.AddLike(projectId)
		if err != nil {
			app.serveResponseError(w, http.StatusInternalServerError, err.Error())
			return

		}

		err = app.Models.Users.CreateUserProjectLikeRelation(user.Pk_User, projectId)
		if err != nil {
			app.serveResponseError(w, http.StatusInternalServerError, err.Error())

		}

		data := map[string]any{
			"user":     user,
			"project":  projectId,
			"Relation": "Like added and created user-projectliked-relation",
		}

		err = app.WriteJSON(w, data, "Response: ")
		if err != nil {
			app.serveResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}

		return

	}

	err = app.Models.Projects.SubLike(projectId)
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())

	}

	err = app.Models.Users.DeleteUserProjectLikeRelation(user.Pk_User, projectId)
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())

	}

	data := map[string]any{
		"user":     user,
		"project":  projectId,
		"Relation": "Like substringed and removed user-projectliked-relation",
	}

	err = app.WriteJSON(w, data, "Response: ")
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
func (app *Application) NewDislikeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	projectId, err := app.ReadUrlId(ps)
	if err != nil {
		app.serveResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	userMailContext := r.Context().Value("mail")

	userMail, ok := userMailContext.(string)
	if !ok {
		app.serveResponseError(w, http.StatusInternalServerError, "JWT mal formated")
		return
	}

	user, err := app.Models.Users.GetByMail(userMail)
	if err != nil {
		app.serveResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	hasAlreadyDisliked, err := app.Models.Users.HasAlreadyDislikedProject(user.Pk_User, projectId)

	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !hasAlreadyDisliked {

		hasAlreadyLiked, err := app.Models.Users.HasAlreadyLikedProject(user.Pk_User, projectId)
		if err != nil {
			app.serveResponseError(w, http.StatusInternalServerError, err.Error())
			return

		}

		if hasAlreadyLiked {
			err = app.Models.Projects.SubLike(projectId)
			if err != nil {
				app.serveResponseError(w, http.StatusInternalServerError, err.Error())
				return

			}

			err = app.Models.Users.DeleteUserProjectLikeRelation(user.Pk_User, projectId)
			if err != nil {
				app.serveResponseError(w, http.StatusInternalServerError, err.Error())
				return

			}

		}

		err = app.Models.Projects.AddDislike(projectId)
		if err != nil {
			app.serveResponseError(w, http.StatusInternalServerError, err.Error())
			return

		}

		err = app.Models.Users.CreateUserProjectDisLikeRelation(user.Pk_User, projectId)
		if err != nil {
			app.serveResponseError(w, http.StatusInternalServerError, err.Error())

		}

		data := map[string]any{
			"user":     user,
			"project":  projectId,
			"Relation": "Dislike added and created user-projectdisliked-relation",
		}

		err = app.WriteJSON(w, data, "Response: ")
		if err != nil {
			app.serveResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}

		return

	}

	err = app.Models.Projects.SubDislike(projectId)
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())

	}

	err = app.Models.Users.DeleteUserProjectDislikeRelation(user.Pk_User, projectId)
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())

	}

	data := map[string]any{
		"user":     user,
		"project":  projectId,
		"Relation": "Disliked substringed and removed user-projectdisliked-relation",
	}

	err = app.WriteJSON(w, data, "Response: ")
	if err != nil {
		app.serveResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
