package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sjadczak/webdev-go/lenslocked/models"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("lenslocked> Error loading .env file")
	}

	host := os.Getenv("SMTP_HOST")
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatal("lenslocked> SMTP Port not found in ENV")
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	email := models.Email{
		To:        "test@test.com",
		Subject:   "This is a test email",
		Plaintext: "This is the body of the email",
		HTML:      `<h1>Hello there buddy!</h1><p>This is an email</p><p>Hope you enjoy it</p>`,
	}

	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})

	err = es.Send(email)
	if err != nil {
		panic(err)
	}
	fmt.Println("lenslocked> message sent.")

	//msg :=
	//msg.SetHeader("To", to)
	//msg.SetHeader("From", from)
	//msg.SetHeader("Subject", subject)
	//msg.SetBody("text/plain", plaintext)
	//msg.AddAlternative("text/html", html)

	//msg.WriteTo(os.Stdout)

	//dialer := mail.NewDialer(host, port, username, password)
	//err := dialer.DialAndSend(msg)
	//if err != nil {
	//panic(err)
	//}

	//fmt.Println("Message sent.")
	//sender, err := dialer.Dial()
	//if err != nil {
	//panic(err)
	//}
	//defer sender.Close()

	//sender.Send(from, []string{to}, msg)

}
