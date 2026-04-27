package config

import (
	"fmt"
	"os"
)

type Env struct {
	VarDir string
}

func LoadEnv() *Env {
	env := &Env{
		VarDir: "var",
	}

	home, err := os.UserHomeDir()
    if err != nil {
        home = "~"
    }

	env.VarDir = fmt.Sprintf("%s/.went/var", home)

	return env
}

func (e *Env) Validate() error {
	return nil
}

func (e *Env) GetLogsDir() string {
	return fmt.Sprintf("%s/logs", e.VarDir)
}
