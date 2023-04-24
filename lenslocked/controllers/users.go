package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/sjadczak/webdev-go/lenslocked/context"
	"github.com/sjadczak/webdev-go/lenslocked/models"
)

type Users struct {
	Templates struct {
		New            Template
		SignIn         Template
		ForgotPassword Template
		CheckYourEmail Template
		ResetPassword  Template
	}
	SessionService       *models.SessionService
	UserService          *models.UserService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
}

// Serve user registration
func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, r, data)
}

// Create new user from registration form
func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		// TODO: long term display warning about not being able to sign user in.
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

// Server user signin
func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}

	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
}

// Verify user credentials
func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")

	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	q := r.URL.Query()
	next := q.Get("next")
	if next == "" {
		next = "/users/me"
	}

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, next, http.StatusFound)
}

// Serve current user info
func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := context.User(ctx)
	if user == nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	fmt.Fprintf(w, "Current user: %v\n", user.Email)
}

// Sign current user out
func (u Users) SignOut(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

// Serve reset password page
func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.ForgotPassword.Execute(w, r, data)
}

// Handle password reset submission
func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		// TODO: Handle other cases in the future. For instance, if a user does not exist
		// with that email address.
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	vals := url.Values{
		"token": {pwReset.Token},
	}
	resetURL := "http://localhost:3000/reset-pw?" + vals.Encode()

	err = u.EmailService.ForgotPassword(data.Email, resetURL)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	// Don't render reset token here!!
	u.Templates.CheckYourEmail.Execute(w, r, data)
}

func (u Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token")
	u.Templates.ResetPassword.Execute(w, r, data)
}

func (u Users) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}
	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")

	user, err := u.PasswordResetService.Consume(data.Token)
	if err != nil {
		// TODO: distinguish between types of errors
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	// TODO: Update the user's password

	// Sign the use rin now that their password has been reset.
	// Any errros from this point onwards should redirect the user to the sign in page
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

type UserMiddleware struct {
	SessionService *models.SessionService
}

// If user has cookies, add current user to context for downstream handlers
func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := readCookie(r, CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := umw.SessionService.User(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// If user found in context, process. Else redirect to signing.
// Protects access to downstream handlers
func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := context.User(ctx)
		if user == nil {
			target := r.URL.Path
			path := *r.URL
			path.Path = "/signin"
			q := path.Query()
			q.Add("next", target)
			path.RawQuery = q.Encode()

			http.Redirect(w, r, path.String(), http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
