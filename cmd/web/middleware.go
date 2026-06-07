package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline';")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src 'self' fonts.gstatic.com;")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		// w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

// because the middleware is a method on the application struct,
// it also has access to the handler dependencies, such as the structured logger.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			method = r.Method
			url    = r.URL.String()
			proto  = r.Proto
		)

		app.logger.Info("request received: ", "ip", ip, "proto", proto, "method", method, "url", url)

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// if panic occurs, recover() will return the value passed to panic() and we can log it.
			pv := recover()

			if pv != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, r, fmt.Errorf("%v", pv))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			// exit. so that no subsequent handlers in the chain are executed.
			return
		}

		// make sure pages requiring authentication are NOT stored in the users' browser cache;
		// or any other intermediary cache.
		w.Header().Add("Cache-Control", "no-store")

		// call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// NOTE: http.Request is designed to be immutable
// If middleware could freely mutate the request object in-place,
// it would be harder to reason about who changed what.
// Also, is it actually copying everything? No. It's a shallow copy.
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if no "authenticatedUserID" value in the session, it'll return 0.
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// otherwise, check if the user ID exists in the database.
		// If it does, then we know the user is authenticated.
		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedCtxKey, true)
			r = r.WithContext(ctx)
		}

		// call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

func preventCSRF(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	// set the base cookie for CSRF token to be HttpOnly and Secure, and with a path of "/".
	// This means the cookie will be sent with all requests to our application, but it won't be accessible via JavaScript.
	// This is a good security practice because it helps to prevent cross-site scripting (XSS) attacks from stealing the CSRF token.
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		// SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}
