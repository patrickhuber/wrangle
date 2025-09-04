package config

import (
	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
)

type localFactory struct {
	filesystem       fs.FS
	path             filepath.Provider
	os               os.OS
	resolver         config.GlobResolver
	errorIfNotExists bool
}

// NewLocalFactory creates a new local configuration provider.
// TODO: see if one provider per file can be created
func NewLocalFactory(
	fs fs.FS,
	path filepath.Provider,
	os os.OS,
	resolver config.GlobResolver,
	errorIfNotExists bool) config.Factory {

	return &localFactory{
		filesystem:       fs,
		path:             path,
		os:               os,
		resolver:         resolver,
		errorIfNotExists: errorIfNotExists,
	}
}

func (p *localFactory) Providers() ([]config.Provider, error) {
	// get the current working directory
	workingDirectory, err := p.os.WorkingDirectory()
	if err != nil {
		return nil, err
	}

	// use the config.NewGlobUp function to load the local configuration
	glob := config.NewGlobUp(
		p.filesystem,
		p.path,
		p.resolver,
		workingDirectory, ".wrangle.*")

	return glob.Providers()
}
