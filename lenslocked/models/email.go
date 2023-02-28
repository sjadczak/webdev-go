package models

import (
	"fmt"

	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender = "support@lenslocked.com"
)

type Email struct {
	From      string
	To        string
	Subject   string
	Plaintext string
	HTML      string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService(cfg SMTPConfig) *EmailService {
	es := EmailService{
		dialer: mail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password),
	}

	return &es
}

type EmailService struct {
	// DefaultSender is used as the default sender when on isn't provided for an email.
	// This is also used in functions where the amil is predetermined, like the forgotten
	// password email.
	DefaultSender string

	// unexported fields
	dialer *mail.Dialer
}

func (es *EmailService) Send(email Email) error {
	msg := mail.NewMessage()

	// Set headers
	msg.SetHeader("To", email.To)
	es.setFrom(msg, email)
	msg.SetHeader("Subject", email.Subject)

	// Set body
	switch {
	case email.Plaintext != "" && email.HTML != "":
		msg.SetBody("text/plain", email.Plaintext)
		msg.AddAlternative("text/html", email.HTML)
	case email.Plaintext != "":
		msg.SetBody("text/plain", email.Plaintext)
	case email.HTML != "":
		msg.SetBody("text/html", email.HTML)
	}

	// Send
	err := es.dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}
	return nil
}

func (es *EmailService) ForgotPassword(to, resetURL string) error {
	email := Email{
		Subject:   "Reset your password",
		To:        to,
		Plaintext: "To reset your password, please visit the following link: " + resetURL,
		HTML: `<p>To reset your password, please visit the following link: <a 
		href="` + resetURL + `">` + resetURL + `</a></p>`,
	}

	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("forgot password email: %w", err)
	}
	return nil
}

func (es *EmailService) setFrom(msg *mail.Message, email Email) {
	var from string
	switch {
	case email.From != "":
		from = email.From
	case es.DefaultSender != "":
		from = es.DefaultSender
	default:
		from = DefaultSender
	}

	msg.SetHeader("From", from)
}
