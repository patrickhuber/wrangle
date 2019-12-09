package packages

import (
	"fmt"
	"text/tabwriter"

	"github.com/patrickhuber/wrangle/feed"
	"github.com/patrickhuber/wrangle/ui"
)

type service struct {
	console     ui.Console
	feedService feed.FeedService
}

// Service lists all packages in the configuration
type Service interface {
	List() error
}

// NewService returns a new packages command object
func NewService(feedService feed.FeedService, console ui.Console) Service {
	return &service{
		feedService: feedService,
		console:     console}
}

func (s *service) List() error {

	// create the tab writer and write out the header
	w := tabwriter.NewWriter(s.console.Out(), 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "name\tversion")
	fmt.Fprintln(w, "----\t-------")

	response, err := s.feedService.List(&feed.FeedListRequest{})
	if err != nil {
		return err
	}

	for _, pkg := range response.Packages {
		for _, ver := range pkg.Versions {
			fmt.Fprintf(w, "%s\t%s", pkg.Name, ver.Version)
			fmt.Fprintln(w)
		}
	}
	return w.Flush()
}
