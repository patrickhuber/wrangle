package services

import (
	"fmt"

	"github.com/patrickhuber/wrangle/collections"

	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/ui"
)

type envService struct {
	console    ui.Console
	dictionary collections.Dictionary
}

// EnvService defines an env command
type EnvService interface {
	Execute() error
}

// NewEnvService creates a new env command
func NewEnvService(console ui.Console, dictionary collections.Dictionary) EnvService {
	return &envService{console: console,
		dictionary: dictionary}
}

func (e *envService) Execute() error {
	configFile, _ := e.dictionary.Lookup(global.ConfigFileKey)
	packagePath, _ := e.dictionary.Get(global.PackagePathKey)
	fmt.Fprintf(e.console.Out(), "%s=%s", global.PackagePathKey, packagePath)
	fmt.Fprintln(e.console.Out())
	fmt.Fprintf(e.console.Out(), "%s=%s", global.ConfigFileKey, configFile)
	fmt.Fprintln(e.console.Out())
	return nil
}
