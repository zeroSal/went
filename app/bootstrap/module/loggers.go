package module

import (
	"context"
	"fmt"
	"went/app/config"

	"github.com/zeroSal/went-logger/logger"
	"go.uber.org/fx"
)

type AuditLogger struct {
	*logger.FileLogger
}

func AuditLoggerProvider(
	lc fx.Lifecycle,
	env *config.Env,
) *AuditLogger {
	path := fmt.Sprintf("%s/audit.log", env.GetLogsDir())
	l := logger.NewFileLogger(path, "audit", logger.LevelInfo)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return l.Close()
		},
	})

	return &AuditLogger{l}
}

type ErrorLogger struct {
	*logger.FileLogger
}

func ErrorLoggerProvider(
	lc fx.Lifecycle,
	env *config.Env,
) *ErrorLogger {
	path := fmt.Sprintf("%s/error.log", env.GetLogsDir())
	l := logger.NewFileLogger(path, "error", logger.LevelError)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return l.Close()
		},
	})

	return &ErrorLogger{l}
}
