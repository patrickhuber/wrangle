package actions

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/wrangle/pkg/archive"
)

type extractProvider struct {
	factory archive.Factory
	logger  log.Logger
	path    filepath.Processor
}

type Extract struct {
	Details *ExtractDetails `yaml:"extract" mapstructure:"extract"`
}

type ExtractDetails struct {
	Archive string `yaml:"archive"`
	Out     string `yaml:"out"`
}

func NewExtractProvider(factory archive.Factory, path filepath.Processor, logger log.Logger) Provider {
	return &extractProvider{
		factory: factory,
		logger:  logger,
		path:    path,
	}
}

func (p *extractProvider) Type() string {
	return "extract"
}

func (p *extractProvider) Execute(t *Action, ctx *Metadata) error {
	extract, err := p.Encode(t)
	if err != nil {
		return err
	}

	archive := p.path.Join(ctx.PackageVersionPath, extract.Details.Archive)
	p.logger.Debugf("extracting %s to %s", archive, ctx.PackageVersionPath)

	provider, err := p.factory.Select(extract.Details.Archive)
	if err != nil {
		return err
	}
	return provider.Extract(archive, ctx.PackageVersionPath, extract.Details.Out)
}

func (p *extractProvider) Encode(t *Action) (*Extract, error) {
	if strings.TrimSpace(t.Type) == "" {
		return nil, fmt.Errorf("task Type is empty")
	}
	if t.Type != p.Type() {
		return nil, fmt.Errorf("invalid task type, expected '%s' found '%s'", p.Type(), t.Type)
	}
	archive, err := t.GetStringParameter("archive")
	if err != nil {
		return nil, err
	}

	out, ok, err := t.GetOptionalStringParameter("out")
	if err != nil {
		return nil, err
	}
	if !ok {
		out = "."
	}

	return &Extract{
		Details: &ExtractDetails{
			Archive: archive,
			Out:     out,
		}}, nil
}
