package bootstrap

import (
	"fmt"
	"os"
	"went/app/bootstrap/module"
	"went/app/config"

	"go.uber.org/fx"
)

var Init = fx.Options(
	fx.Invoke(InitWorkingDirs),
	fx.Invoke(ValidateEnv),
	fx.Invoke(InitLoggers),
)

func InitWorkingDirs(env *config.Env) error {
	dirs := []string{
		env.GetLogsDir(),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("cannot create dir %s", dir)
		}
	}

	return nil
}

func ValidateEnv(env *config.Env) error {
	return env.Validate()
}

func InitLoggers(
	auditLogger *module.AuditLogger,
	errorLogger *module.ErrorLogger,
) error {
	if err := auditLogger.Init(); err != nil {
		return err
	}
	if err := errorLogger.Init(); err != nil {
		return err
	}
	return nil
}
