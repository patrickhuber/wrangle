package services

import (
	"fmt"
	"sort"

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
	List() map[string]string
}

// NewEnvService creates a new env command
func NewEnvService(console ui.Console, dictionary collections.Dictionary) EnvService {
	return &envService{console: console,
		dictionary: dictionary}
}

func (e *envService) Execute() error {
	variables := e.List()
	e.print(variables)
	return nil
}

func (e *envService) List() map[string]string {
	keys := []string{
		global.BinPathKey,
		global.CachePathKey,
		global.ConfigFileKey,
	}
	variables := map[string]string{}
	for _, k := range keys {
		value, _ := e.dictionary.Get(k)
		variables[k] = value
	}
	return variables
}

func (e *envService) print(variables map[string]string) {
	keys := make([]string, len(variables))
	i := 0
	for k := range variables {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(e.console.Out(), "%s=%s", k, variables[k])
		fmt.Fprintln(e.console.Out())
	}
}
