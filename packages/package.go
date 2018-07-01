package packages

// Package represents an interface for a binary package of software
type Package interface {
	Download() Download
	Extract() Extract
}

type pkg struct {
	download Download
	extract  Extract
}

// New creates a new package ready for download
func New(download Download, extract Extract) Package {
	return &pkg{download: download, extract: extract}
}

func (p *pkg) Download() Download {
	return p.download
}

func (p *pkg) Extract() Extract {
	return p.extract
}
