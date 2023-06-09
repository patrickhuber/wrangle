package services

import (
	"fmt"

	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/env"
)

type ExportRequest struct {
	Shell string
}

type Export interface {
	Execute(r *ExportRequest) error
}

type export struct {
	env     env.Environment
	shells  map[string]shellhook.Shell
	console console.Console
}

func NewExport(env env.Environment, shells map[string]shellhook.Shell, console console.Console) Export {
	return &export{
		env:     env,
		shells:  shells,
		console: console,
	}
}

func (e *export) Execute(r *ExportRequest) error {
	shell, ok := e.shells[r.Shell]
	if !ok {
		return fmt.Errorf("invalid shell %s", shell)
	}
	vars := e.env.Export()
	rendered := shell.Export(vars)
	_, err := fmt.Fprint(e.console.Out(), rendered)
	return err
}
