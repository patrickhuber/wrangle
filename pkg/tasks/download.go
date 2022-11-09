package tasks

import (
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/mitchellh/mapstructure"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
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
	fs     filesystem.FileSystem
}

// NewDownloadProvider creates a new download provider
func NewDownloadProvider(logger log.Logger, fs filesystem.FileSystem) Provider {
	return &downloadProvider{
		name:   "download",
		logger: logger,
		fs:     fs,
	}
}

func (p *downloadProvider) Decode(object interface{}) (*Task, error) {
	// map structure to Download
	var download = &Download{}
	err := mapstructure.Decode(object, download)
	if err != nil {
		return nil, err
	}
	// map Download to Task
	return &Task{
		Type: p.name,
		Parameters: map[string]interface{}{
			"url": download.Details.URL,
			"out": download.Details.Out,
		},
	}, nil
}

func (p *downloadProvider) Encode(tsk *Task) (*Download, error) {
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

func (p *downloadProvider) Execute(t *Task, ctx *Metadata) error {
	download, err := p.Encode(t)
	if err != nil {
		return err
	}
	return p.execute(download, ctx)
}

func (p *downloadProvider) execute(download *Download, ctx *Metadata) error {

	out := crosspath.Join(ctx.PackageVersionPath, download.Details.Out)
	url := download.Details.URL

	p.logger.Printf("downloading '%s' to '%s'", url, out)

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
	err = p.fs.MkdirAll(directory, 0755)
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
