// Package task supports asynchronous task processing
package task

import (
	asynqamanger "github.com/qiushiyan/simplebank/business/task/asynq"
	simplemanager "github.com/qiushiyan/simplebank/business/task/simple"
)

// Manager defines the behavior for managing asynchronous tasks.
type Manager interface {
	Start() error
	Close() error
	CreateTask(taskType string, payload any) (string, error)
	CancelTask(id string) error
	GetTaskStatus(id string) (string, error)
}

var _ Manager = (*asynqamanger.AsynqManager)(nil)
var _ Manager = (*simplemanager.SimpleManager)(nil)
