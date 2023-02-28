package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/sjadczak/webdev-go/lenslocked/controllers"
	"github.com/sjadczak/webdev-go/lenslocked/migrations"
	"github.com/sjadczak/webdev-go/lenslocked/models"
	"github.com/sjadczak/webdev-go/lenslocked/templates"
	"github.com/sjadczak/webdev-go/lenslocked/views"
)

func main() {
	// Set up database
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Set up services
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}

	// Set up middleware
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	csrfKey := "a29fghadf092yh3rhglaisdfh2as$@Fas"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		// TODO: fix before deploying
		csrf.Secure(false),
	)

	// Set up handlers
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}

	// Set up router & routes
	r := chi.NewRouter()
	r.Use(csrfMw, umw.SetUser)

	tpl := views.Must(views.ParseFS(templates.FS, "layout.gohtml", "home.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "contact.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "faq.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "signup.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "signin.gohtml"))
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.SignOut)

	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	tpl = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "notfound.gohtml"))
	r.NotFound(controllers.StaticHandler(tpl))

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe("127.0.0.1:3000", r)
}

//func TimerMiddleware(next http.Handler) http.Handler {
//return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//start := time.Now()
//next.ServeHTTP(w, r)
//fmt.Println("Request time:", time.Since(start))
//})
//}
