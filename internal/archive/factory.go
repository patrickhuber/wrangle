package archive

import (
	"fmt"
	"path"
	"strings"

	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
)

type Factory interface {
	Select(archive string) (Provider, error)
}

type factory struct {
	providers map[string]Provider
}

const (
	Tar   = "tar"
	Targz = "tgz"
	Zip   = "zip"
)

func NewFactory(fs fs.FS, path *filepath.Processor) Factory {
	return &factory{
		providers: map[string]Provider{
			Tar:   NewTar(fs, path),
			Targz: NewTarGz(fs, path),
			Zip:   NewZip(fs, path),
		},
	}
}

func (p *factory) Select(archive string) (Provider, error) {
	switch {
	case strings.HasSuffix(archive, ".tgz"), strings.HasSuffix(archive, ".tar.gz"):
		return p.providers[Targz], nil
	case strings.HasSuffix(archive, ".zip"):
		return p.providers[Zip], nil
	case strings.HasSuffix(archive, ".tar"):
		return p.providers[Tar], nil
	}
	return nil, fmt.Errorf("unable to find provider for extension '%s'", path.Ext(archive))
}
