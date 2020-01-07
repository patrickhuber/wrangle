package processes

import (
	"fmt"
	"text/tabwriter"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/ui"
)

// Service provides a service for processes
type Service interface {
	List(cfg *config.Config) error
}

// NewService creates a new processes service
func NewService(console ui.Console) Service {
	return &service{
		console: console,
	}
}

type service struct {
	console ui.Console
}

func (s *service) List(cfg *config.Config) error {

	w := tabwriter.NewWriter(s.console.Out(), 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "name")
	fmt.Fprintln(w, "----")
	for _, process := range cfg.Processes {
		fmt.Fprintf(w, "%s", process.Name)
		fmt.Fprintln(w)
	}
	return w.Flush()
}
