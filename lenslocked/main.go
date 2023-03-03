package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"github.com/sjadczak/webdev-go/lenslocked/controllers"
	"github.com/sjadczak/webdev-go/lenslocked/migrations"
	"github.com/sjadczak/webdev-go/lenslocked/models"
	"github.com/sjadczak/webdev-go/lenslocked/templates"
	"github.com/sjadczak/webdev-go/lenslocked/views"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	// TODO: PSQL from ENV
	cfg.PSQL = models.DefaultPostgresConfig()

	// SMTP
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Port = smtpPort
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	// TODO: CSRF from ENV
	cfg.CSRF.Key = "a29fghadf092yh3rhglaisdfh2as$@Fas"
	cfg.CSRF.Secure = false

	// Server config
	cfg.Server.Address = "127.0.0.1:3000"

	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Set up services
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	emailService := models.NewEmailService(models.SMTPConfig{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
	})
	pwrService := &models.PasswordResetService{
		DB: db,
	}

	// Set up middleware
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
	)

	// Set up handlers
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		EmailService:         emailService,
		PasswordResetService: pwrService,
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
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "forgotpassword.gohtml"))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "checkemail.gohtml"))
	usersC.Templates.ResetPassword = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "reset-password.gohtml"))
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.SignOut)
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)

	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	tpl = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "notfound.gohtml"))
	r.NotFound(controllers.StaticHandler(tpl))

	fmt.Printf("Starting the server on %s...\n", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		log.Fatal(err)
	}
}

//func TimerMiddleware(next http.Handler) http.Handler {
//return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//start := time.Now()
//next.ServeHTTP(w, r)
//fmt.Println("Request time:", time.Since(start))
//})
//}
