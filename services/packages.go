package services

import (
	"fmt"
	"text/tabwriter"

	"github.com/patrickhuber/wrangle/feed"
	"github.com/patrickhuber/wrangle/ui"
)

type packagesService struct {
	console     ui.Console
	feedService feed.FeedService
}

// PackagesService lists all packages in the configuration
type PackagesService interface {
	List() error
}

// NewPackagesService returns a new packages command object
func NewPackagesService(feedService feed.FeedService, console ui.Console) PackagesService {
	return &packagesService{
		feedService: feedService,
		console:     console}
}

func (service *packagesService) List() error {

	// create the tab writer and write out the header
	w := tabwriter.NewWriter(service.console.Out(), 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "name\tversion")
	fmt.Fprintln(w, "----\t-------")

	response, err := service.feedService.List(&feed.FeedListRequest{})
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
