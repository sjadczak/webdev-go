package controllers

import (
	"fmt"
	"net/http"
)

type Users struct {
	Templates struct {
		New Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(4000 << 10)
	if err != nil {
		panic(err)
	}
	f, handler, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintf(w, "Email:    %v\n", r.FormValue("email"))
	fmt.Fprintf(w, "Password: %v\n", r.FormValue("password"))
	fmt.Fprintf(w, "Filename: %v\n", handler.Filename)
}
