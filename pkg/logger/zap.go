package logger

import (
	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	log *zap.Logger
}

func NewZapLogger() *ZapLogger {
	config := zap.NewProductionConfig()
	logger, _ := config.Build()
	return &ZapLogger{log: logger}
}

func (l *ZapLogger) Debug(msg string, params ...zapcore.Field) {
	l.log.Debug(msg, params...)
}

func (l *ZapLogger) Info(msg string, params ...zapcore.Field) {
	l.log.Info(msg, params...)
}

func (l *ZapLogger) Warn(msg string, params ...zapcore.Field) {
	l.log.Warn(msg, params...)
}

func (l *ZapLogger) Error(err error, params ...zapcore.Field) {
	l.log.Error(err.Error(), params...)
}

func (l *ZapLogger) ModifyToSentry(client *sentry.Client) error {
	cfg := zapsentry.Configuration{
		Level:             zapcore.ErrorLevel,
		EnableBreadcrumbs: true,
		BreadcrumbLevel:   zapcore.InfoLevel,
		Tags: map[string]string{
			"component": "system",
		},
	}

	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromClient(client))

	if err != nil {
		return err
	}

	l.log = zapsentry.AttachCoreToLogger(core, l.log)
	l.log = l.log.With(zapsentry.NewScope())

	return nil
}
