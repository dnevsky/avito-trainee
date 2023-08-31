package logger

import (
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, params ...zapcore.Field)
	Info(msg string, params ...zapcore.Field)
	Warn(msg string, params ...zapcore.Field)
	Error(err error, params ...zapcore.Field)

	ModifyToSentry(client *sentry.Client) error
}

type LogManager struct {
	Logger
}

func NewLogManager() *LogManager {
	return &LogManager{
		NewZapLogger(),
	}
}
