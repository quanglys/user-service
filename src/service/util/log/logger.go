package log

import "go.uber.org/zap"

type Logger struct {
	logger *zap.Logger
}

func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{logger: logger}
}

func (l *Logger) Info(msg string, field ...zap.Field) {
	l.logger.Info(msg, append(field, zap.Any("msg", msg))...)
}

func (l *Logger) Error(msg string, field ...zap.Field) {
	l.logger.Error(msg, append(field, zap.Any("msg", msg))...)
}

func (l *Logger) Debug(msg string, field ...zap.Field) {
	l.logger.Debug(msg, append(field, zap.Any("msg", msg))...)
}

func (l *Logger) Warn(msg string, field ...zap.Field) {
	l.logger.Warn(msg, append(field, zap.Any("msg", msg))...)
}
