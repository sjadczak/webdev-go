package controllers

import (
	"net/http"
)

type Users struct {
	Templates struct {
		New Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	// view to render
	u.Templates.New.Execute(w, nil)
}
