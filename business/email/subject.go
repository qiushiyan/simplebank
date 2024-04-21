package email

import "errors"

const (
	SubjectWelcome = "welcome"
	SubjectVerify  = "verify"
)

// templates map the email subject to the template file.
var templates = map[string]string{
	SubjectWelcome: "welcome.html",
	SubjectVerify:  "verify.html",
}
var ErrInvalidSubject = errors.New("invalid email subject, choose from 'welcome' or 'verify'")

type SubjectWelcomeData struct {
	Username string
}

type SubjectVerifyData struct {
	Username string
	Link     string
}
