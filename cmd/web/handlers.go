package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gesatessa/sbox/internal/models"
	"github.com/gesatessa/sbox/internal/validator"
)

// type snippetCreateForm struct {
// 	Title       string
// 	Content     string
// 	Expires     int
// 	FieldErrors map[string]string
// }

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

// func (app *application) home(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Add("Server", "Go")

// 	files := []string{
// 		"./ui/html/base.tpl.html",
// 		"./ui/html/partials/nav.tpl.html",
// 		"./ui/html/pages/home.tpl.html",
// 	}

// 	ts, err := template.ParseFiles(files...)
// 	if err != nil {
// 		app.serverError(w, r, err)
// 		return
// 	}

// 	// err = ts.Execute(w, nil)
// 	// Execute the template named "base" from the parsed template set.
// 	// This is the name of the template defined in the base.tpl.html file.
// 	// now, the template set contains named templates: base, title, nav, main. instead of containing HTML directly.
// 	err = ts.ExecuteTemplate(w, "base", nil)
// 	if err != nil {
// 		app.serverError(w, r, err)
// 		return
// 	}
// }

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// panic("oooooooooops")
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.tpl.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		// if no record is found, return a 404 Not Found response to the client.
		// otherwise, if there is an error (e.g., database connection issue),
		// log the error and return a 500 Internal Server Error response.
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tpl.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 7,
	}

	app.render(w, r, http.StatusOK, "create.tpl.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// limit the size of the request body to prevent malicious clients
	// from sending large requests that could consume server resources.
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.FormValue("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "title cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "title cannot be longer than 100 characters")
	form.CheckField(validator.NotBlank(form.Content), "content", "content field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 30), "expires", "expires must be eitehr: 1, 7 or 30")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
