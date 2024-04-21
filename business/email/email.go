// Package email provides email sending functionality.
package email

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"testing"
)

// Sender represents an email sender.
type Sender interface {
	Send(payload SenderPayload) error
}

type SenderPayload struct {
	To      string
	Subject string
	Data    any
}

// returns the interpolated HTML content given the payload
func getEmailHTML(payload SenderPayload) ([]byte, error) {
	t, ok := templates[payload.Subject]
	if !ok {
		return []byte{}, fmt.Errorf("template not found for subject %s", payload.Subject)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return []byte{}, fmt.Errorf("error getting current working directory: %v", err)
	}

	var path string
	if testing.Testing() {
		path = filepath.Join(cwd, "..", "..", "zarf", "email", t)
	} else {
		path = filepath.Join(cwd, "zarf", "email", t)
	}

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return []byte{}, fmt.Errorf("error parsing template: %v", err)
	}

	// Create an instance of WelcomeEmailData with the username.
	data := struct{ Data any }{
		Data: payload.Data,
	}

	// Use a bytes.Buffer to capture the output of the template execution.
	var buf bytes.Buffer

	// Execute the template, passing in the data struct. Write the output to the buffer.
	if err := tmpl.Execute(&buf, data); err != nil {
		return []byte{}, fmt.Errorf("error executing template: %v", err)
	}

	return buf.Bytes(), nil
}
