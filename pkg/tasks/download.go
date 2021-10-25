package tasks

import (
	"fmt"
	"io"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/ilog"
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
	cfg    *config.Config
	logger ilog.Logger
	fs     filesystem.FileSystem
}

// NewDownloadProvider creates a new download provider
func NewDownloadProvider(cfg *config.Config, logger ilog.Logger) Provider {
	return &downloadProvider{
		cfg:  cfg,
		name: "download",
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

func (p *downloadProvider) Execute(t *Task, ctx *Context) error {
	download, err := p.Encode(t)
	if err != nil {
		return err
	}
	return p.execute(download, ctx)
}

func (p *downloadProvider) execute(download *Download, ctx *Context) error {

	out := crosspath.Join(ctx.PackageVersionPath, download.Details.Out)
	url := download.Details.URL

	p.logger.Printf("downloading '%s' to '%s'", url, out)
	p.logger.Println()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error downloading '%s'. http status code: '%d'. http status: '%s'",
			url,
			resp.StatusCode,
			resp.Status)
	}

	// create the file
	file, err := p.fs.Create(out)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the body to file
	_, err = io.Copy(file, resp.Body)

	return err
}
