package store

import (
	"fmt"
	"text/tabwriter"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/ui"
)

type service struct {
	console ui.Console
}

// Service provides a service over stores
type Service interface {
	List(cfg *config.Config) error
}

// NewService provides a stores service for manipulating stores
func NewService(console ui.Console) Service {
	return &service{
		console: console}
}

func (s *service) List(cfg *config.Config) error {
	w := tabwriter.NewWriter(s.console.Out(), 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "name\ttype")
	fmt.Fprintln(w, "----\t----")
	for _, item := range cfg.Stores {
		fmt.Fprintf(w, "%s\t%s", item.Name, item.StoreType)
		fmt.Fprintln(w)
	}
	return w.Flush()
}
