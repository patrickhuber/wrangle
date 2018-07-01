package packages

import (
	"github.com/spf13/afero"
)

// Package represents an interface for downloading a binary package
type Package interface {
	URL() string
	Out() string
}

type pkg struct {
	url        string
	outPath    string
	fileSystem afero.Fs
}

// New creates a new package ready for download
func New(url string, outPath string) Package {
	return &pkg{url: url, outPath: outPath}
}

func (p *pkg) URL() string {
	return p.url
}

func (p *pkg) Out() string {
	return p.outPath
}
