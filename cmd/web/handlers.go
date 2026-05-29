package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gesatessa/sbox/internal/models"
)

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

	title := r.FormValue("title")
	content := r.FormValue("content")
	expires, err := strconv.Atoi(r.FormValue("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// initialize a map to hold any validation errors that occur during form processing.
	fieldErrors := make(map[string]string)
	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "This field cannot be longer than 100 characters"
	}

	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "This field cannot be blank"
	}

	// make sure expires is one of the permitted values (1, 7, or 30).
	if expires != 1 && expires != 7 && expires != 30 {
		fieldErrors["expires"] = "This field must equal 1, 7, or 30"
	}

	if len(fieldErrors) > 0 {
		fmt.Fprint(w, fieldErrors)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
