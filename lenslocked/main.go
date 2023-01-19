package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sjadczak/webdev-go/lenslocked/controllers"
	"github.com/sjadczak/webdev-go/lenslocked/templates"
	"github.com/sjadczak/webdev-go/lenslocked/views"
)

func main() {
	r := chi.NewRouter()

	tpl := views.Must(views.ParseFS(templates.FS, "layout.gohtml", "home.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "contact.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "faq.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "notfound.gohtml"))
	r.NotFound(controllers.StaticHandler(tpl))

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe("127.0.0.1:3000", r)
}
