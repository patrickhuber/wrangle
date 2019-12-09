package env

import (
	"fmt"
	"sort"

	"github.com/patrickhuber/wrangle/collections"

	"github.com/patrickhuber/wrangle/ui"
)

type envService struct {
	console     ui.Console
	dataService DataService
}

// Service defines an env command
type Service interface {
	Execute() error
	List() map[string]string
}

// NewService creates a new env command
func NewService(console ui.Console, dictionary collections.Dictionary) Service {
	dataService := NewDataService(dictionary)
	return &envService{console: console,
		dataService: dataService}
}

func (e *envService) Execute() error {
	variables := e.List()
	e.print(variables)
	return nil
}

func (e *envService) List() map[string]string {
	return e.dataService.List()
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
