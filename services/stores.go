package services

import (
	"fmt"
	"text/tabwriter"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/ui"
)

type stores struct {
	console ui.Console
	loader  config.Loader
}

// StoresService provides a service over stores
type StoresService interface {
	List(configFile string) error
}

// NewStoresService provides a stores service for manipulating stores
func NewStoresService(console ui.Console, loader config.Loader) StoresService {
	return &stores{
		console: console,
		loader:  loader}
}

func (service *stores) List(configFile string) error {
	w := tabwriter.NewWriter(service.console.Out(), 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "name\ttype")
	fmt.Fprintln(w, "----\t----")
	cfg, err := service.loader.LoadConfig(configFile)
	if err != nil {
		return err
	}
	for _, item := range cfg.Stores {
		fmt.Fprintf(w, "%s\t%s", item.Name, item.StoreType)
		fmt.Fprintln(w)
	}
	return w.Flush()
}
