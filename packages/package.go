package packages

import (
	"strings"
)

// Package represents an interface for a binary package of software
type Package interface {
	Download() Download
	Extract() Extract
	Version() string
	Alias() string
}

type pkg struct {
	download Download
	extract  Extract
	version  string
	alias    string
	name     string
}

// New creates a new package ready for download
func New(name string, version string, alias string, download Download, extract Extract) Package {
	return &pkg{
		download: interpolateDownload(version, download),
		extract:  interpolateExtract(version, extract),
		version:  version,
		alias:    alias,
		name:     name}
}

func (p *pkg) Download() Download {
	return p.download
}

func (p *pkg) Extract() Extract {
	return p.extract
}

func (p *pkg) Version() string {
	return p.version
}

func (p *pkg) Alias() string {
	return p.alias
}

func interpolateDownload(version string, download Download) Download {
	if download == nil {
		return nil
	}
	url := replaceVersion(download.URL(), version)
	outFile := replaceVersion(download.OutFile(), version)
	outFolder := replaceVersion(download.OutFolder(), version)
	return NewDownload(url, outFolder, outFile)
}

func interpolateExtract(version string, extract Extract) Extract {
	if extract == nil {
		return nil
	}
	filter := replaceVersion(extract.Filter(), version)
	outFile := replaceVersion(extract.OutFile(), version)
	outFolder := replaceVersion(extract.OutFolder(), version)
	return NewExtract(filter, outFolder, outFile)
}

func replaceVersion(input string, version string) string {
	return strings.Replace(input, "((version))", version, -1)
}
