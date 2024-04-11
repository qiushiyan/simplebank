package simplemanager

import (
	"time"

	"go.uber.org/zap"
)

type SimpleManager struct {
	log  *zap.SugaredLogger
	quit chan struct{}
}

func New(log *zap.SugaredLogger) *SimpleManager {
	return &SimpleManager{}
}

func (m *SimpleManager) Start() error {
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				select {

				default:
					// No task available, wait for the next tick
				}
			}
		}
	}()

	m.waitForSignal()

	return nil
}

func (m *SimpleManager) Close() error {
	return nil
}

func (m *SimpleManager) CreateTask(taskType string, payload any) (string, error) {
	return "", nil
}

func (m *SimpleManager) CancelTask(id string) error {
	return nil
}

func (m *SimpleManager) GetTaskStatus(id string) (string, error) {
	return "", nil
}

func (m *SimpleManager) waitForSignal() {
	<-m.quit
}
