package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// func home(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Add("Server", "Go")

// 	// template set: a collection of one or more templates that are parsed from a set of files.
// 	// The template set is used to execute a template by name.
// 	ts, err := template.ParseFiles("./ui/html/pages/home.tpl.html")
// 	if err != nil {
// 		log.Print(err.Error())
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	err = ts.Execute(w, nil)
// 	if err != nil {
// 		log.Print(err.Error())
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}
// }

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	files := []string{
		"./ui/html/base.tpl.html",
		"./ui/html/partials/nav.tpl.html",
		"./ui/html/pages/home.tpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template named "base" from the parsed template set.
	// This is the name of the template defined in the base.tpl.html file.
	// now, the template set contains named templates: base, title, nav, main. instead of containing HTML directly.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	msg := fmt.Sprintf("snippet id %d ...", id)
	w.Write([]byte(msg))
}
