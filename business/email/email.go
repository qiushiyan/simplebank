// Package email provides email sending functionality.
package email

import (
	"bytes"
	"encoding/json"
	"errors"
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
	To      string `json:"to"`
	Subject string `json:"subject"`
	Data    any    `json:"data"`
}

func (sp *SenderPayload) UnmarshalJSON(data []byte) error {
	type Alias SenderPayload
	aux := &struct {
		Data json.RawMessage `json:"data"`
		*Alias
	}{
		Alias: (*Alias)(sp),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("error unmarshalling SenderPayload: %v", err)
	}
	sp.To = aux.To
	sp.Subject = aux.Subject

	if json.Valid(aux.Data) {
		var err error
		switch sp.Subject {
		case SubjectWelcome:
			var Data SubjectWelcomeData
			err = json.Unmarshal(aux.Data, &Data)
			sp.Data = Data
		case SubjectVerify:
			var Data SubjectVerifyData
			err = json.Unmarshal(aux.Data, &Data)
			sp.Data = Data
		case SubjectReport:
			var Data SubjectReportData
			err = json.Unmarshal(aux.Data, &Data)
			sp.Data = Data
		default:
			return fmt.Errorf("unknown subject type: %s", sp.Subject)
		}
		if err != nil {
			sp.Data = nil
			return fmt.Errorf("error unmarshalling data for %s: %v", sp.Subject, err)
		}
	} else {
		sp.Data = nil
		return errors.New("can't unmarshal invalid Data property in SenderPayload")
	}

	return nil
}

// returns the interpolated HTML content given the payload
func getEmailHTML(payload *SenderPayload) ([]byte, error) {
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
		path = filepath.Join(cwd, "..", "..", "zarf", "email", "templates", t)
	} else {
		path = filepath.Join(cwd, "zarf", "email", "templates", t)
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
