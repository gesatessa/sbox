package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gesatessa/sbox/internal/models"
	"github.com/gesatessa/sbox/internal/validator"
)

// `stract tags` tell th decoder how to map HTML form values into the different struct fields.
// NOTE: type conversions are handled automatically. (expires from string to int)
type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type userForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

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

	data.Form = snippetCreateForm{
		Expires: 7,
	}

	app.render(w, r, http.StatusOK, "create.tpl.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// limit the size of the request body to prevent malicious clients
	// from sending large requests that could consume server resources.
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB

	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		// if there is a problem decoding the form, send 400 bad request response to the client.
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "title cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "title cannot be longer than 100 characters")
	form.CheckField(validator.NotBlank(form.Content), "content", "content field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 30), "expires", "expires must be eitehr: 1, 7 or 30")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// add the flash message if the post created successfully.
	app.sessionManager.Put(r.Context(), "flash", "snippet created successfully.")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) userSignUp(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userForm{}
	app.render(w, r, http.StatusOK, "signup.tpl.html", data)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userForm{}
	app.render(w, r, http.StatusOK, "login.tpl.html", data)
}

func (app *application) userSignUpPost(w http.ResponseWriter, r *http.Request) {
	var form userForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "name cannot be blank")
	// form.CheckField(validator.NotBlank(form.Email), "email", "email cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "enter a valid email address")

	form.CheckField(validator.MinChars(form.Password, 8), "password", "password must be at least 8 charactors long")
	form.CheckField(validator.MaxBytes(form.Password, 72), "password", "password must not be more than 72 bytes long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tpl.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "email address already registered")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tpl.html", data)
		} else {
			app.serverError(w, r, err)
		}

		return
	}
	// add a confirmation flash message that user signup worked.
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please login.")
	app.logger.Info("a new user signed up", "email", form.Email)

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// if the credentials match, we add the user's ID to their session data
// this way, for future requests, we'll know the user is already authenticated, and who that is.
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("credentials are invalid")

			// redisplay the login page
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tpl.html", data)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	// it's good practice to generate a new session ID
	// when the authentication state or privilege levels change for the user.
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// add the ID of the current user to their session, so that they are now "logged in"
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

// renew the session ID & remove the `authenticatedUserID` value from the session
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "you've been logged out.")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
