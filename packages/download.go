package packages

// Download represents an interface for downloading a package
type Download interface {
	URL() string
	Out() string
	OutFolder() string
}

type download struct {
	url       string
	out       string
	outFolder string
}

// NewDownload Creates a new download instance
func NewDownload(url string, out string, outFolder string) Download {
	return &download{url: url, out: out, outFolder: outFolder}
}

func (d *download) URL() string {
	return d.url
}

func (d *download) Out() string {
	return d.out
}

func (d *download) OutFolder() string {
	return d.outFolder
}
