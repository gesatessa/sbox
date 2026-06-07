package main

import (
	"net/http"

	"github.com/gesatessa/sbox/ui"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// fileServer := http.FileServer(http.Dir(app.cfg.static))
	// mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	// dynamic middleware chain
	dynamic := alice.New(app.sessionManager.LoadAndSave, app.authenticate, preventCSRF)

	// alice.ThenFunc() returns `http.Handler` (and not http.HandlerFunc)
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignUp))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignUpPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// authenticated-only routes:
	protected := dynamic.Append(app.requireAuth)
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// standard middleware chain for all routes
	standardMW := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standardMW.Then(mux)
}
