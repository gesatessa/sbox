package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// func (app *application) routes() *http.ServeMux {
// 	mux := http.NewServeMux()

// 	fileServer := http.FileServer(http.Dir(app.cfg.static))
// 	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

// 	mux.HandleFunc("GET /{$}", app.home)
// 	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
// 	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

// 	return mux
// }

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(app.cfg.static))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)

	// standard middleware chain for all routes
	standardMW := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standardMW.Then(mux)
}
