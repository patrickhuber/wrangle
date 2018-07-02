package packages

import "path/filepath"

// Download represents an interface for downloading a package
type Download interface {
	URL() string
	OutFile() string
	OutFolder() string
	OutPath() string
}

type download struct {
	url       string
	outFile   string
	outFolder string
	outPath   string
}

// NewDownload Creates a new download instance
func NewDownload(url string, out string, outFolder string) Download {
	outPath := filepath.Join(outFolder, out)
	return &download{
		url:       url,
		outFile:   out,
		outFolder: outFolder,
		outPath:   outPath}
}

func (d *download) URL() string {
	return d.url
}

func (d *download) OutFile() string {
	return d.outFile
}

func (d *download) OutFolder() string {
	return d.outFolder
}

func (d *download) OutPath() string {
	return d.outPath
}
