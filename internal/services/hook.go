package services

import (
	"fmt"

	"github.com/patrickhuber/go-cross/console"
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-shellhook"
)

type hook struct {
	shells  map[string]shellhook.Shell
	console console.Console
	env     env.Environment
}

type Hook interface {
	Execute(r *HookRequest) error
}

type HookRequest struct {
	Executable string
	Shell      string
}

func NewHook(env env.Environment, shells map[string]shellhook.Shell, console console.Console) Hook {
	return &hook{
		shells:  shells,
		console: console,
		env:     env,
	}
}

func (h *hook) Execute(r *HookRequest) error {
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
