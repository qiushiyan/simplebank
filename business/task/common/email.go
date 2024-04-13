package taskcommon

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"path/filepath"
	"text/template"

	"github.com/jordan-wright/email"
)

const (
	SubjectWelcome = "welcome"
)

var templates = map[string]string{
	SubjectWelcome: "welcome.html",
}

type EmailDeliveryPayload struct {
	To       string `json:"to"`
	Username string `json:"username"`
	Subject  string `json:"subject"`
}

func NewEmailDeliveryPayload(to, username, subject string) EmailDeliveryPayload {
	return EmailDeliveryPayload{
		To:       to,
		Username: username,
		Subject:  subject,
	}
}

type EmailSender struct {
	address  string
	password string
}

func NewEmailSender(address, password string) *EmailSender {
	return &EmailSender{
		address:  address,
		password: password,
	}
}

func (s EmailSender) Send(payload EmailDeliveryPayload) error {
	var subject string
	var err error

	switch payload.Subject {
	case SubjectWelcome:
		subject = "Welcome to SimpleBank!"
	default:
		subject = "Notification from SimpleBank"
	}

	e := email.NewEmail()
	e.From = "SimpleBank <simplebankdev@gmail.com>"
	e.To = []string{payload.To}
	e.Subject = subject
	e.HTML, err = getEmailHTML(payload)

	if err != nil {
		return fmt.Errorf("could not get email template: %w", err)
	}

	return e.Send(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", s.address, s.password, "smtp.gmail.com"),
	)

}

func getEmailHTML(payload EmailDeliveryPayload) ([]byte, error) {
	t, ok := templates[payload.Subject]
	if !ok {
		return []byte{}, fmt.Errorf("template not found for subject %s", payload.Subject)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return []byte{}, fmt.Errorf("error getting current working directory: %v", err)
	}

	path := filepath.Join(cwd, "zarf", "email", t)

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return []byte{}, fmt.Errorf("error parsing template: %v", err)
	}

	// Create an instance of WelcomeEmailData with the username.
	data := struct{ Username string }{
		Username: payload.Username,
	}

	// Use a bytes.Buffer to capture the output of the template execution.
	var buf bytes.Buffer

	// Execute the template, passing in the data struct. Write the output to the buffer.
	if err := tmpl.Execute(&buf, data); err != nil {
		return []byte{}, fmt.Errorf("error executing template: %v", err)
	}

	return buf.Bytes(), nil
}
