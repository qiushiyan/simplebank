package email

import "errors"

const (
	SubjectWelcome = "welcome"
	SubjectVerify  = "verify"
	SubjectReport  = "report"
)

// templates map the email subject to the template file.
var templates = map[string]string{
	SubjectWelcome: "welcome.html",
	SubjectVerify:  "verify.html",
	SubjectReport:  "report.html",
}

var ErrInvalidSubject = errors.New(
	"invalid email subject, choose from 'welcome', 'verify' or 'report'",
)

type SubjectWelcomeData struct {
	Username string `json:"username"`
}

type SubjectVerifyData struct {
	Username string `json:"username"`
	Link     string `json:"link"`
}

type SubjectReportData struct {
	Username string `json:"username"`
}
