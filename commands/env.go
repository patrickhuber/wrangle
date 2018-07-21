package commands

import (
	"fmt"

	"github.com/patrickhuber/wrangle/collections"

	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/ui"
)

type env struct {
	console    ui.Console
	dictionary collections.Dictionary
}

// Env defines an env command
type Env interface {
	Execute() error
}

// NewEnv creates a new env command
func NewEnv(console ui.Console, dictionary collections.Dictionary) Env {
	return &env{console: console,
		dictionary: dictionary}
}

func (e *env) Execute() error {
	configFile, _ := e.dictionary.Lookup(global.ConfigFileKey)
	packagePath, _ := e.dictionary.Get(global.PackagePathKey)
	fmt.Fprintf(e.console.Out(), "%s=%s", global.PackagePathKey, packagePath)
	fmt.Fprintln(e.console.Out())
	fmt.Fprintf(e.console.Out(), "%s=%s", global.ConfigFileKey, configFile)
	fmt.Fprintln(e.console.Out())
	return nil
}
