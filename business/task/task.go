// Package task supports asynchronous task processing
package task

import (
	"context"

	asynqamanger "github.com/qiushiyan/simplebank/business/task/asynq"
	taskcommon "github.com/qiushiyan/simplebank/business/task/common"
	simplemanager "github.com/qiushiyan/simplebank/business/task/simple"
)

// Manager defines the behavior for managing asynchronous tasks.
type Manager interface {
	Start() error
	Close() error
	CreateTask(ctx context.Context, taskType string, payload any) (string, error)
	CancelTask(id string) error
	GetTaskState(id string) (*taskcommon.State, error)
}

var _ Manager = (*asynqamanger.AsynqManager)(nil)
var _ Manager = (*simplemanager.SimpleManager)(nil)
