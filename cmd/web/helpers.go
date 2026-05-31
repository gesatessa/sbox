package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		url    = r.URL.Path
		trace  = string(debug.Stack())
	)
	app.logger.Error(err.Error(), "method", method, "url", url, "trace", trace)
	// app.logger.Error(err.Error(), "method", method, "url", url)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, statusCode int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s not found", page)
		app.serverError(w, r, err)
		return
	}

	// write the template to a buffer instead of straight to the http.ResponseWriter. This allows us
	// to check for any errors that occur during template execution before we start writing the response.
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	w.Header().Set("Content-Length", "not a number")
	w.WriteHeader(statusCode)
	buf.WriteTo(w)

}

func (app *application) newTemplateData(r *http.Request) templateData {
	// add flash message to the template data, if there is one.
	return templateData{
		CurrentYear: time.Now().Year(),
		Flash:       app.sessionManager.PopString(r.Context(), "flash"),
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var targetErr *form.InvalidDecoderError
		if errors.As(err, &targetErr) {
			panic(err)
		}

		return err
	}

	return nil
}
