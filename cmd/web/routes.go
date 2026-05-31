package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(app.cfg.static))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// dynamic middleware chain
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// alice.ThenFunc() returns `http.Handler` (and not http.HandlerFunc)
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))

	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignUp))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignUpPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))

	// standard middleware chain for all routes
	standardMW := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standardMW.Then(mux)
}
