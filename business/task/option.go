package task

import "fmt"

// available task manager options
var (
	OptionSimple = Option{"simple"}
	OptionAsynq  = Option{"asynq"}
)

var options = map[string]Option{
	OptionSimple.name: OptionSimple,
	OptionAsynq.name:  OptionAsynq,
}

// Option represents an available task manager option.
type Option struct {
	name string
}

// ParseOption parses the string value and returns a task option if one exists.
func ParseOption(value string) (Option, error) {
	option, exists := options[value]
	if !exists {
		allOptions := make([]string, 0, len(options))
		for name := range options {
			allOptions = append(allOptions, name)
		}
		return Option{}, fmt.Errorf(
			"invalid task option: %s, available options: %v",
			value,
			allOptions,
		)
	}

	return option, nil
}

// MustParseOption parses the string value and returns a task option if one exists. If
// an error occurs the function panics.
func MustParseOption(value string) Option {
	option, err := ParseOption(value)
	if err != nil {
		panic(err)
	}

	return option
}

// Name returns the name of the task option.
func (o Option) Name() string {
	return o.name
}

// UnmarshalText implement the unmarshal interface for JSON conversions.
func (o *Option) UnmarshalText(data []byte) error {
	option, err := ParseOption(string(data))
	if err != nil {
		return err
	}

	o.name = option.name
	return nil
}

// MarshalText implement the marshal interface for JSON conversions.
func (o Option) MarshalText() ([]byte, error) {
	return []byte(o.name), nil
}

// Equal checks if two task options are equal.
func (o Option) Equal(other Option) bool {
	return o.name == other.name
}
