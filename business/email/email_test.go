package email

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {
	// os.Chdir(os.Getenv("WORKSPACE_DIR"))

	cases := []struct {
		name    string
		payload SenderPayload
		checker func(err error)
	}{
		{
			name: "ok",
			payload: SenderPayload{
				To:      "user@gmail.com",
				Data:    SubjectWelcomeData{Username: "user"},
				Subject: SubjectWelcome,
			},

			checker: func(err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "invalid-subject",
			payload: SenderPayload{
				To:      "user@gmail.com",
				Data:    SubjectWelcomeData{Username: "user"},
				Subject: "unknown-subject",
			},
			checker: func(err error) {
				require.ErrorIs(t, err, ErrInvalidSubject)
			},
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sender := NewGmailSender("", "")
			err := sender.Send(tc.payload)
			tc.checker(err)
		})
	}

}
