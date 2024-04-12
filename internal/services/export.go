package services

import (
	"fmt"

	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/wrangle/internal/envdiff"
	"github.com/patrickhuber/wrangle/internal/global"
)

type Export interface {
	// Execute executes the export of the changes in the context of the given shell
	Execute(shell string, changes []envdiff.Change) error
}

type export struct {
	shells  map[string]shellhook.Shell
	console console.Console
}

func NewExport(
	shells map[string]shellhook.Shell,
	console console.Console,
) Export {
	return &export{
		shells:  shells,
		console: console,
	}
}

func (e *export) Execute(shell string, changes []envdiff.Change) error {
	var err error
	out := e.console.Out()
	s, ok := e.shells[shell]
	if !ok {
		return fmt.Errorf("invalid shell %s", shell)
	}
	for _, change := range changes {
		switch c := change.(type) {
		case envdiff.Add:
			str := s.Export(c.Key, c.Value)
			_, err = fmt.Fprintln(out, str)
		case envdiff.Remove:
			str := s.Unset(c.Key)
			_, err = fmt.Fprintln(out, str)
		case envdiff.Update:
			str := s.Export(c.Key, c.Value)
			_, err = fmt.Fprintln(out, str)
		}
		if err != nil {
			return err
		}
	}

	// save the diff
	diffStr, err := envdiff.Encode(changes)
	if err != nil {
		return err
	}
	exportStr := s.Export(global.EnvDiff, diffStr)
	_, err = fmt.Fprintln(out, exportStr)
	return err
}
