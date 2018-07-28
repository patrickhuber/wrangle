package renderers

import (
	"fmt"

	"github.com/patrickhuber/wrangle/collections"
)

type factory struct {
	env collections.Dictionary
}

// Factory defines a new renderer factory
type Factory interface {
	Create(shell string) (Renderer, error)
}

// NewFactory creates a new factory
func NewFactory(env collections.Dictionary) Factory {
	return &factory{
		env: env,
	}
}

func (f *factory) createFromDefault() (Renderer, error) {
	format := PosixFormat
	if _, ok := f.env.Lookup("PSModulePath"); ok {
		format = PowershellFormat
	}
	return f.createFromFormat(format)
}

func (f *factory) createFromFormat(format string) (Renderer, error) {
	switch format {
	case PowershellFormat:
		return NewPowershell(), nil
	case PosixFormat:
		return NewPosix(), nil
	default:
		return nil, fmt.Errorf("unrecognized format '%s'", format)
	}
}

func (f *factory) Create(format string) (Renderer, error) {
	if format == "" {
		return f.createFromDefault()
	}
	return f.createFromFormat(format)
}
