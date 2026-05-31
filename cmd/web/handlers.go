package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gesatessa/sbox/internal/models"
	"github.com/gesatessa/sbox/internal/validator"
)

// `stract tags` tell th decoder how to map HTML form values into the different struct fields.
// NOTE: type conversions are handled automatically. (expires from string to int)
type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

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

	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		// if there is a problem decoding the form, send 400 bad request response to the client.
		app.clientError(w, http.StatusBadRequest)
		return
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
	// add the flash message if the post created successfully.
	app.sessionManager.Put(r.Context(), "flash", "snippet created successfully.")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
