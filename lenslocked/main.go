package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/sjadczak/webdev-go/lenslocked/views"
)

func executeTemplate(w http.ResponseWriter, filepath string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl, err := views.Parse(filepath)
	if err != nil {
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}
	tpl.Execute(w, nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "home.gohtml")
	executeTemplate(w, tplPath)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "contact.gohtml")
	executeTemplate(w, tplPath)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "faq.gohtml")
	executeTemplate(w, tplPath)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "notfound.gohtml")
	executeTemplate(w, tplPath)
}

func configureRouter() chi.Router {
	// Allocate new mux
	r := chi.NewRouter()

	// Add routes
	r.Get("/", indexHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.NotFound(notFoundHandler)

	return r
}

func main() {
	router := configureRouter()

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe("127.0.0.1:3000", router)
}
