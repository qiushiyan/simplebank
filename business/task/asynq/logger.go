package asynqamanger

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/qiushiyan/simplebank/foundation/logger"
)

// Logger adapts logger.Logger to asynq logger.
type Logger struct {
	log *logger.Logger
}

func (l *Logger) Debug(args ...interface{}) {
	l.log.Debug(context.Background(), "asynq", args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.log.Error(context.Background(), "asynq", args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.log.Info(context.Background(), "asynq", args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.log.Warn(context.Background(), "asynq", args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.log.Error(context.Background(), "asynq", args...)
}

var _ asynq.Logger = (*Logger)(nil)
