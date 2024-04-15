package taskcommon

import (
	"errors"
)

type State struct {
	// Task ID
	Id string `json:"id"`
	// Task Type
	Type string `json:"type"`
	// One of "in_progress", "completed", "failed", or "other"
	Status Status `json:"status"`
	// The Error message if the task failed, otherwise omitted
	Error string `json:"error,omitempty"`
}

var (
	StatusInProgress = Status{"in_progress"}
	StatusCompleted  = Status{"completed"}
	StatusFailed     = Status{"failed"}
	StatusOther      = Status{"other"}
)

// Set of known statuses.
var statuses = map[string]Status{
	StatusInProgress.name: StatusInProgress,
	StatusCompleted.name:  StatusCompleted,
	StatusFailed.name:     StatusFailed,
	StatusOther.name:      StatusOther,
}

var ErrInvalidStatus = errors.New("invalid status")

// Status represents the status of a task.
type Status struct {
	name string
}

// ParseStatus parses the string value and returns a status if one exists.
func ParseStatus(value string) (Status, error) {
	status, exists := statuses[value]
	if !exists {
		return Status{}, ErrInvalidStatus
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
func (s Status) MarshalText() ([]byte, error) {
	return []byte(s.name), nil
}
