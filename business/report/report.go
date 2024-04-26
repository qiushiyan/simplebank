package report

import "context"

// Generator defines behavior for managing reports
type Generator interface {
	New(context.Context, string, any) (string, error)
}
