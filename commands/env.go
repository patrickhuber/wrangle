package commands

import (
	"fmt"
	"os"

	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/ui"
)

type env struct {
	console ui.Console
}

// Env defines an env command
type Env interface {
	Execute() error
}

// NewEnv creates a new env command
func NewEnv(console ui.Console) Env {
	return &env{console: console}
}

func (e *env) Execute() error {
	configFile := os.Getenv(global.ConfigFileKey)
	packagePath := os.Getenv(global.PackagePathKey)
	fmt.Fprintf(e.console.Out(), "%s=%s", global.PackagePathKey, packagePath)
	fmt.Fprintln(e.console.Out())
	fmt.Fprintf(e.console.Out(), "%s=%s", global.ConfigFileKey, configFile)
	fmt.Fprintln(e.console.Out())
	return nil
}
