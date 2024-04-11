package taskcommon

type EmailDeliveryPayload struct {
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Template string `json:"template"`
}

func NewEmailDeliveryPayload(to, subject, template string) EmailDeliveryPayload {
	return EmailDeliveryPayload{
		To:       to,
		Subject:  subject,
		Template: template,
	}
}
