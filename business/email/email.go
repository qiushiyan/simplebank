// Package email provides email sending functionality.
package email

import (
	"bytes"
	"encoding/json"
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
		return err
	}
	sp.To = aux.To
	sp.Subject = aux.Subject

	switch {
	case json.Valid(aux.Data):
		var reportData SubjectReportData
		if err := json.Unmarshal(aux.Data, &reportData); err == nil {
			sp.Data = reportData
		} else {
			return fmt.Errorf("error unmarshalling data: %v", err)
		}
	default:
		sp.Data = nil
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
