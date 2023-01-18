package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/sjadczak/webdev-go/lenslocked/controllers"
	"github.com/sjadczak/webdev-go/lenslocked/views"
)

func main() {
	r := chi.NewRouter()

	tpl := views.Must(views.Parse(filepath.Join("templates", "home.gohtml")))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.Parse(filepath.Join("templates", "contact.gohtml")))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.Parse(filepath.Join("templates", "faq.gohtml")))
	r.Get("/faq", controllers.StaticHandler(tpl))

	tpl = views.Must(views.Parse(filepath.Join("templates", "notfound.gohtml")))
	r.NotFound(controllers.StaticHandler(tpl))

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe("127.0.0.1:3000", r)
}
