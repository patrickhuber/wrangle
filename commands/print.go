package commands

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/patrickhuber/cli-mgr/store"
	"github.com/patrickhuber/cli-mgr/ui"
	"github.com/spf13/afero"
)

type print struct {
	fileSystem afero.Fs
	platform   string
	console    ui.Console
	manager    store.Manager
}

// Print represents an environment command
type Print interface {
	Execute(params RunCommandParams) error
}

// NewPrint creates a new environment command
func NewPrint(
	manager store.Manager,
	fileSystem afero.Fs,
	platform string,
	console ui.Console) Print {
	return &print{
		manager:    manager,
		fileSystem: fileSystem,
		platform:   platform,
		console:    console}
}

func (cmd *print) Execute(params RunCommandParams) error {

	processName := params.ProcessName()
	environmentName := params.EnvironmentName()

	if processName == "" {
		return errors.New("process name is required for the run command")
	}

	if environmentName == "" {
		return errors.New("environment name is required for the run command")
	}

	cfg := params.Config()
	if cfg == nil {
		return errors.New("unable to load configuration")
	}

	pipeline := store.NewPipeline(cmd.manager, cfg)
	environment, err := pipeline.Run(environmentName, processName)
	if err != nil {
		return err
	}
	renderer := NewEnvVarRenderer(cmd.platform)
	fmt.Fprint(cmd.console.Out(), renderer.RenderEnvironment(environment.Vars))
	return nil
}
