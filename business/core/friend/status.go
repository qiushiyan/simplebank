package friend

import "fmt"

// Status represents the status of a friend request.
type Status struct {
	name string
}

var (
	StatusPending  = Status{"pending"}
	StatusAccepted = Status{"accepted"}
	StatusRejected = Status{"rejected"}
)

var InvalidStatusError = fmt.Errorf("status must be one of pending, accepted, and rejected")

// set of known statuses.
var statuses = map[string]Status{
	StatusPending.name:  StatusPending,
	StatusAccepted.name: StatusAccepted,
	StatusRejected.name: StatusRejected,
}

// ParseStatus parses the string value and returns a status if one exists.
func ParseStatus(value string) (Status, error) {
	status, exists := statuses[value]
	if !exists {
		return Status{}, InvalidStatusError
	}

	return status, nil
}

// MustParseStatus parses the string value and returns a status if one exists. If
// an error occurs the function panics.
func MustParseStatus(value string) Status {
	status, err := ParseStatus(value)
	if err != nil {
		panic(err)
	}

	return status
}

// Name returns the name of the status.
func (s Status) Name() string {
	return s.name
}

// UnmarshalText implement the unmarshal interface for JSON conversions.
func (s *Status) UnmarshalText(data []byte) error {
	status, err := ParseStatus(string(data))
	if err != nil {
		return err
	}

	s.name = status.name
	return nil
}

// MarshalText implement the marshal interface for JSON conversions.
func (s *Status) MarshalText() ([]byte, error) {
	return []byte(s.name), nil
}

// Equal provides support for the go-cmp package and testing.
func (s *Status) Equal(s2 Status) bool {
	return s.name == s2.name
}
