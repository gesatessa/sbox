package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/gesatessa/sbox/internal/models"
)

// templateData is a struct that holds the dynamic data
// that we want to pass to our HTML templates when rendering them.
type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	// convert the time to UTC before formatting it.
	// This ensures that the date and time are displayed in a consistent way regardless of the server's local tz

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var funcs = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// initialize a new map to act as the cache
	data := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// get the fileName from the filePath
		fileName := filepath.Base(page)

		files := []string{
			"./ui/html/base.tpl.html",
			"./ui/html/partials/nav.tpl.html",
			page,
		}

		// ts, err := template.ParseFiles(files...)
		ts, err := template.New(fileName).Funcs(funcs).ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		data[fileName] = ts
	}

	return data, nil
}
