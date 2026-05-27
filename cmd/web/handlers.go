package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	// w.Header().Add("Server", "Go")

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

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// dummy data to test the handler
	// title := "task 1"
	// content := "this is the content of task 1"
	// expires := 1

	title := "Why learn Go?"
	content := "Go is designed to be simple & efficient. Especially a great choice for building web applications & microservices."
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
