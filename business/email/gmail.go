package email

import (
	"fmt"
	"net/smtp"
	"testing"

	"github.com/jordan-wright/email"
)

// GmailSender is an email sender that sends emails using Gmail.
type GmailSender struct {
	// Address is the email address to send from, currently only supports Gmail.
	address string
	// Password is the **app password** for the gmail account
	password string
}

func NewGmailSender(address, password string) *GmailSender {
	return &GmailSender{
		address:  address,
		password: password,
	}
}

func (s GmailSender) Send(payload SenderPayload) error {
	var subject string
	var err error

	switch payload.Subject {
	case SubjectWelcome:
		subject = "Welcome to SimpleBank!"
	default:
		return ErrInvalidSubject
	}

	e := email.NewEmail()
	e.From = "SimpleBank <simplebankdev@gmail.com>"
	e.To = []string{payload.To}
	e.Subject = subject

	e.HTML, err = getEmailHTML(payload)

	if err != nil {
		return fmt.Errorf("could not get email template: %w", err)
	}

	if testing.Testing() {
		return nil
	}

	return e.Send(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", s.address, s.password, "smtp.gmail.com"),
	)

}
