package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	files := []string{
		"./ui/html/base.tpl.html",
		"./ui/html/partials/nav.tpl.html",
		"./ui/html/pages/home.tpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// err = ts.Execute(w, nil)
	// Execute the template named "base" from the parsed template set.
	// This is the name of the template defined in the base.tpl.html file.
	// now, the template set contains named templates: base, title, nav, main. instead of containing HTML directly.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	msg := fmt.Sprintf("snippet id %d ...", id)
	w.Write([]byte(msg))
}
