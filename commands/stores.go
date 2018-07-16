package commands

import (
	"fmt"
	"text/tabwriter"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/ui"
)

type stores struct {
	console ui.Console
}

type Stores interface {
	Execute(cfg *config.Config) error
}

func NewStores(console ui.Console) Stores {
	return &stores{
		console: console}
}

func (cmd *stores) Execute(cfg *config.Config) error {
	w := tabwriter.NewWriter(cmd.console.Out(), 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "name\ttype")
	fmt.Fprintln(w, "----\t----")
	for _, item := range cfg.Stores {
		fmt.Fprintf(w, "%s\t%s", item.Name, item.StoreType)
		fmt.Fprintln(w)
	}
	w.Flush()
	return nil
}
