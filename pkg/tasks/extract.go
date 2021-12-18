package tasks

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/wrangle/pkg/archive"
)

type extractProvider struct {
	factory archive.Factory
}

type Extract struct {
	Details *ExtractDetails `yaml:"extract" mapstructure:"extract"`
}

type ExtractDetails struct {
	Archive string `yaml:"archive"`
	Out     string `yaml:"out"`
}

func NewExtractProvider(factory archive.Factory) Provider {
	return &extractProvider{
		factory: factory,
	}
}

func (p *extractProvider) Type() string {
	return "extract"
}

func (p *extractProvider) Execute(t *Task, m *Metadata) error {
	extract, err := p.Encode(t)
	if err != nil {
		return err
	}
	provider, err := p.factory.Select(extract.Details.Archive)
	if err != nil {
		return err
	}
	return provider.Extract(extract.Details.Archive, "", extract.Details.Out)
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
	out, err := t.GetStringParameter("out")
	if err != nil {
		return nil, err
	}

	return &Extract{
		Details: &ExtractDetails{
			Archive: archive,
			Out:     out,
		}}, nil
}
