package services

import (
	"fmt"
	"text/tabwriter"

	"github.com/patrickhuber/wrangle/filepath"

	"github.com/spf13/afero"

	"github.com/patrickhuber/wrangle/ui"
)

type packagesService struct {
	console    ui.Console
	fileSystem afero.Fs
}

// PackagesService lists all packages in the configuration
type PackagesService interface {
	List(packagePath string) error
}

// NewPackagesService returns a new packages command object
func NewPackagesService(fileSystem afero.Fs, console ui.Console) PackagesService {
	return &packagesService{
		fileSystem: fileSystem,
		console:    console}
}

func (service *packagesService) List(packagePath string) error {

	// create the tab writer and write out the header
	w := tabwriter.NewWriter(service.console.Out(), 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "name\tversion")
	fmt.Fprintln(w, "----\t-------")

	path := packagePath
	packageFolders, err := afero.ReadDir(service.fileSystem, path)

	if err != nil {
		return err
	}

	for _, packageFolder := range packageFolders {
		if !packageFolder.IsDir() {
			continue
		}
		path := filepath.Join(path, packageFolder.Name())
		packageVersions, err := afero.ReadDir(service.fileSystem, path)
		if err != nil {
			return err
		}
		for _, packageVersion := range packageVersions {
			fmt.Fprintf(w, "%s\t%s", packageFolder.Name(), packageVersion.Name())
			fmt.Fprintln(w)
		}

	}
	return w.Flush()
}
