package renderers

import "fmt"

type factory struct {
	platform string
}

// Factory defines a new renderer factory
type Factory interface {
	Create(shell string) (Renderer, error)
}

// NewFactory creates a new factory
func NewFactory(platform string) Factory {
	return &factory{
		platform: platform,
	}
}

func (f *factory) createFromPlatform(platform string) (Renderer, error) {
	switch platform {
	case "windows":
		return f.createFromShell("powershell")
	default:
		return f.createFromShell("bash")
	}
}

func (f *factory) createFromShell(shell string) (Renderer, error) {
	switch shell {
	case "powershell":
		return NewPowershell(), nil
	case "bash":
		return NewBash(), nil
	default:
		return nil, fmt.Errorf("unrecognized shell '%s'", shell)
	}
}

func (f *factory) Create(shell string) (Renderer, error) {
	if shell == "" {
		return f.createFromPlatform(f.platform)
	}
	return f.createFromShell(shell)
}
