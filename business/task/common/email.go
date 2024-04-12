package taskcommon

type EmailDeliveryPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
}

func NewEmailDeliveryPayload(to, subject string) EmailDeliveryPayload {
	return EmailDeliveryPayload{
		To:      to,
		Subject: subject,
	}
}
