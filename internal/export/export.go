package export

import (
	"fmt"

	"github.com/patrickhuber/go-cross/console"
	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/wrangle/internal/envdiff"
)

type Service interface {
	// Execute executes the export of the changes in the context of the given shell
	Execute(shell string, changes []envdiff.Change) error
}

type service struct {
	shells  map[string]shellhook.Shell
	console console.Console
}

func NewService(
	shells map[string]shellhook.Shell,
	console console.Console,
) Service {
	return &service{
		shells:  shells,
		console: console,
	}
}

func (e *service) Execute(shell string, changes []envdiff.Change) error {
	var err error
	out := e.console.Out()
	s, ok := e.shells[shell]
	if !ok {
		return fmt.Errorf("invalid shell %s", shell)
	}

	// write out the changes as lines
	lines := []string{}
	for _, change := range changes {
		switch c := change.(type) {
		case envdiff.Add:
			str := s.Export(c.Key, c.Value)
			lines = append(lines, str)
		case envdiff.Remove:
			str := s.Unset(c.Key)
			lines = append(lines, str)
		case envdiff.Update:
			str := s.Export(c.Key, c.Value)
			lines = append(lines, str)
		}
	}

	// format the lines into the output
	for _, line := range lines {
		_, err = fmt.Fprintln(out, line)
		if err != nil {
			return err
		}
	}

	return err
}
