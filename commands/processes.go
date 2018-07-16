package commands

import (
	"fmt"
	"text/tabwriter"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/ui"
)

type listProcesses struct {
	console ui.Console
}

type ListProcesses interface {
	Execute(cfg *config.Config, environmentName string) error
}

func NewListProcesses(console ui.Console) ListProcesses {
	return &listProcesses{
		console: console}
}

func (cmd *listProcesses) Execute(cfg *config.Config, environmentName string) error {

	w := tabwriter.NewWriter(cmd.console.Out(), 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "name")
	fmt.Fprintln(w, "----")
	for _, environment := range cfg.Environments {
		if environment.Name != environmentName {
			continue
		}
		for _, process := range environment.Processes {
			fmt.Fprintf(w, "%s", process.Name)
			fmt.Fprintln(w)
		}

		w.Flush()
		return nil
	}
	return fmt.Errorf("unable to find environment name '%s'", environmentName)
}
