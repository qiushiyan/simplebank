package asynqamanger

import (
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

// Logger adapts a zap logger to asynq logger.
type Logger struct {
	log *zap.SugaredLogger
}

func (l *Logger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

var _ asynq.Logger = (*Logger)(nil)
