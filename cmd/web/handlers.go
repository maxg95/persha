package main

import (
	"fmt"
	"net/http"
	"os"

	"persha.maxg95/internal/response"
	"persha.maxg95/internal/validator"
)

var recipient = os.Getenv("RECIPIENT")

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/home.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) insertRequest(w http.ResponseWriter, r *http.Request) {
	_, err := app.db.Exec(`
        CREATE TABLE IF NOT EXISTS requests (
            number TEXT
        )
    `)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	request := r.PostForm.Get("request")

	if !validator.IsValidNumber(request) {
		app.badRequest(w, r, fmt.Errorf("Invalid input."))
		return
	}

	_, err = app.db.Exec("INSERT INTO requests (number) VALUES ($1)", request)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := map[string]interface{}{"Number": request}
	err = app.sendEmail(w, r, recipient, data, "mail.tmpl")
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) pomynky(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/pomynky.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) pomynalni_obidy(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/pomynalni_obidy.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) pomynalni_obidy_lutsk(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/pomynalni_obidy_lutsk.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) vesillia(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/vesillia.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) vesillia_cafe(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/vesillia_cafe.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) vesillia_lutsk(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/vesillia_lutsk.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) keiterinh(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/keiterinh.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) keiterinh_cafe(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/keiterinh_cafe.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) keiterinh_lutsk(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/keiterinh_lutsk.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) sendEmail(w http.ResponseWriter, r *http.Request, recipient string, data map[string]interface{}, templateName string) error {
	err := app.mailer.Send(recipient, data, templateName)
	if err != nil {
		return err
	}
	return nil
}
