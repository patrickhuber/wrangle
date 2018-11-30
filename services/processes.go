package services

import (
	"fmt"
	"text/tabwriter"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/ui"
)

// ProcessesService provides a service for processes
type ProcessesService interface {
	List(configFile string) error
}

// NewProcessesService creates a new processes service
func NewProcessesService(console ui.Console, loader config.Loader) ProcessesService {
	return &processesService{
		console: console,
		loader:  loader,
	}
}

type processesService struct {
	console ui.Console
	loader  config.Loader
}

func (service *processesService) List(configFile string) error {

	w := tabwriter.NewWriter(service.console.Out(), 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "name")
	fmt.Fprintln(w, "----")
	cfg, err := service.loader.LoadConfig(configFile)
	if err != nil {
		return err
	}
	for _, process := range cfg.Processes {
		fmt.Fprintf(w, "%s", process.Name)
		fmt.Fprintln(w)
	}
	return w.Flush()
}
