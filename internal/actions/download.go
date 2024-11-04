package actions

import (
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/mitchellh/mapstructure"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-log"
)

// Download defines the outer structure for a download task
type Download struct {
	Details *DownloadDetails `yaml:"download" mapstructure:"download"`
}

// DownloadDetails define the details for the download
type DownloadDetails struct {
	URL string `yaml:"url"`
	Out string `yaml:"out"`
}

type downloadProvider struct {
	name   string
	logger log.Logger
	fs     fs.FS
	path   filepath.Provider
}

// NewDownloadProvider creates a new download provider
func NewDownloadProvider(logger log.Logger, fs fs.FS, path filepath.Provider) Provider {
	return &downloadProvider{
		name:   "download",
		logger: logger,
		fs:     fs,
		path:   path,
	}
}

func (p *downloadProvider) Decode(object any) (*Action, error) {
	// map structure to Download
	var download = &Download{}
	err := mapstructure.Decode(object, download)
	if err != nil {
		return nil, err
	}
	// map Download to Task
	return &Action{
		Type: p.name,
		Parameters: map[string]any{
			"url": download.Details.URL,
			"out": download.Details.Out,
		},
	}, nil
}

func (p *downloadProvider) Encode(tsk *Action) (*Download, error) {
	var download = &Download{}
	url, err := tsk.GetStringParameter("url")
	if err != nil {
		return nil, err
	}
	out, err := tsk.GetStringParameter("out")
	if err != nil {
		return nil, err
	}
	download.Details = &DownloadDetails{
		URL: url,
		Out: out,
	}
	return download, nil
}

func (p *downloadProvider) Type() string {
	return p.name
}

func (p *downloadProvider) Execute(t *Action, ctx *Metadata) error {
	download, err := p.Encode(t)
	if err != nil {
		return err
	}
	return p.execute(download, ctx)
}

func (p *downloadProvider) execute(download *Download, ctx *Metadata) error {

	// ensure package version path exists
	err := p.fs.MkdirAll(ctx.PackageVersionPath, 0775)
	if err != nil {
		return err
	}

	out := p.path.Join(ctx.PackageVersionPath, download.Details.Out)
	url := download.Details.URL

	p.logger.Debugf("downloading '%s' to '%s'", url, out)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer checkClose(resp.Body, &err)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error downloading '%s'. http status code: '%d'. http status: '%s'",
			url,
			resp.StatusCode,
			resp.Status)
	}

	directory := path.Dir(out)
	p.logger.Debugf("creating %s", directory)
	err = p.fs.MkdirAll(directory, 0775)
	if err != nil {
		return err
	}

	// create the file
	file, err := p.fs.Create(out)
	if err != nil {
		return err
	}

	defer checkClose(file, &err)

	// Write the body to file
	var written int64
	written, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	if written == 0 {
		return fmt.Errorf("zero bytes written to %s", out)
	}
	p.logger.Debugf("%d bytes written to %s", written, out)
	return err
}

// checkClose is used to check the return from Close in a defer
// statement.
func checkClose(c io.Closer, err *error) {
	cerr := c.Close()
	if *err == nil {
		*err = cerr
	}
}
