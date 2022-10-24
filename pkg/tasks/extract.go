package tasks

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/wrangle/pkg/archive"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/ilog"
)

type extractProvider struct {
	factory archive.Factory
	logger  ilog.Logger
}

type Extract struct {
	Details *ExtractDetails `yaml:"extract" mapstructure:"extract"`
}

type ExtractDetails struct {
	Archive string `yaml:"archive"`
	Out     string `yaml:"out"`
}

func NewExtractProvider(factory archive.Factory, logger ilog.Logger) Provider {
	return &extractProvider{
		factory: factory,
		logger:  logger,
	}
}

func (p *extractProvider) Type() string {
	return "extract"
}

func (p *extractProvider) Execute(t *Task, ctx *Metadata) error {
	extract, err := p.Encode(t)
	if err != nil {
		return err
	}

	archive := crosspath.Join(ctx.PackageVersionPath, extract.Details.Archive)
	p.logger.Debugf("extracting %s to %s", archive, ctx.PackageVersionPath)

	provider, err := p.factory.Select(extract.Details.Archive)
	if err != nil {
		return err
	}
	return provider.Extract(archive, ctx.PackageVersionPath, extract.Details.Out)
}

func (p *extractProvider) Encode(t *Task) (*Extract, error) {
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
