package packages

// Download represents an interface for downloading a package
type Download interface {
	URL() string
	Out() string
}

type download struct {
	url string
	out string
}

// NewDownload Creates a new download instance
func NewDownload(url string, out string) Download {
	return &download{url: url, out: out}
}

func (d *download) URL() string {
	return d.url
}

func (d *download) Out() string {
	return d.out
}
