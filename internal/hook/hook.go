package hook

import (
	"fmt"

	"github.com/patrickhuber/go-cross/console"
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-shellhook"
)

type service struct {
	shells  map[string]shellhook.Shell
	console console.Console
	env     env.Environment
}

type Service interface {
	Execute(r *Request) error
}

type Request struct {
	Executable string
	Shell      string
}

func NewService(env env.Environment, shells map[string]shellhook.Shell, console console.Console) Service {
	return &service{
		shells:  shells,
		console: console,
		env:     env,
	}
}

func (h *service) Execute(r *Request) error {
	shell, ok := h.shells[r.Shell]
	if !ok {
		return fmt.Errorf("invalid shell %s", shell)
	}
	result, err := shellhook.Hook(shell, &shellhook.Metadata{
		Executable: r.Executable,
		Name:       "wrangle",
		Args:       []string{"export", shell.Name()},
	})
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(h.console.Out(), result)
	return err
}
