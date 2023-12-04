package services

import (
	"fmt"

	"github.com/patrickhuber/wrangle/internal/template"

	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/go-xplat/console"
)

type ExportRequest struct {
	Shell string
}

type Export interface {
	Execute(r *ExportRequest) error
}

type export struct {
	shells        map[string]shellhook.Shell
	console       console.Console
	configuration Configuration
}

func NewExport(shells map[string]shellhook.Shell, console console.Console, configuration Configuration) Export {
	return &export{
		shells:        shells,
		console:       console,
		configuration: configuration,
	}
}

func (e *export) Execute(r *ExportRequest) error {
	shell, ok := e.shells[r.Shell]
	if !ok {
		return fmt.Errorf("invalid shell %s", shell)
	}

	cfg, err := e.configuration.Get()
	if err != nil {
		return err
	}

	vars := map[string]string{}

	// loop through the variables and interpolate each against the stores
	for k, v := range cfg.Spec.Environment {
		// set v as a template and extract any vars
		t := template.New(v)
		value, err := t.Evaluate()
		if err != nil {
			return err
		}
		vars[k] = fmt.Sprintf("%v", value)
	}
	rendered := shell.Export(vars)
	_, err = fmt.Fprint(e.console.Out(), rendered)
	return err
}
