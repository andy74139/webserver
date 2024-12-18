package infra

import (
	"context"

	"go.uber.org/zap"
)

const loggerKey = "logger_key"

var defaultLogger *zap.SugaredLogger

func SetLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func GetLogger(ctx context.Context) *zap.SugaredLogger {
	logger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger)
	if !ok {
		// TODO: handle nil could cause error
		return nil
	}
	return logger
}

func SetDefaultLogger(logger *zap.SugaredLogger) {
	defaultLogger = logger
}

func GetDefaultLogger() *zap.SugaredLogger {
	return defaultLogger
}
