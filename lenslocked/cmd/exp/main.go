package main

import (
	"fmt"

	"github.com/sjadczak/webdev-go/lenslocked/models"
)

const (
	host     = "sandbox.smtp.mailtrap.io"
	port     = 587
	username = "9785379ab2eb36"
	password = "2c5cede8eab45a"
)

func main() {
	email := models.Email{
		From:      "test@lenslocked.com",
		To:        "steve@jadczak.com",
		Subject:   "This is a test email",
		Plaintext: "This is the body of the email.",
		HTML:      `<h1>Hello there!</h1><p>This is the email</p><p>Hope you enjoy it</p>`,
	}

	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})

	err := es.Send(email)
	if err != nil {
		panic(err)
	}
	fmt.Println("lenslocked> email sent")
}
