package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func greeterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	name := chi.URLParam(r, "name")

	for k, v := range r.URL.Query() {
		fmt.Printf("%v: %v (%T)\n", k, v, v)
	}

	msg := fmt.Sprintf(`
	<h1>Dynamic Greetings, again!</h1>
	<p>Hello, %v</p>
	`, name)

	fmt.Fprint(w, msg)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:stevejadczak@gmail.com\">stevejadczak@gmail.com</a>.</p>")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
	<h1>FAQs</h1>
	<ul>
		<li>
			<p><b>Q: Is there a free version?</b></p>
			<p>A: Yes! We offer a free trial for 30 days on any paid plans.</p>
		</li>
		<li>
			<p><b>Q: What are your support hours?</b></p>
			<p>A: We have support staff answering emails 24/7, though response times may be a bit slower on weekends.</p>
		</li>
		<li>
			<p><b>Q: How do I contact support?</b></p>
			<p>A: Email us - <a href="mailto:support@lenslocked.com">support@lenslocked.com</a></p>
		</li>
	</ul>
	`)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>404 Not Found!</h1><p>Go <a href=\"/\">home</a>, you're drunk.</p>")
}

func configureRouter() chi.Router {
	// Allocate new mux
	r := chi.NewRouter()

	// Add routes
	r.Get("/", indexHandler)

	r.Route("/{name}", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/", greeterHandler)
	})

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
