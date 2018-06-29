package commands

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/patrickhuber/cli-mgr/store"
	"github.com/patrickhuber/cli-mgr/ui"
	"github.com/spf13/afero"
)

type printCommand struct {
	fileSystem afero.Fs
	platform   string
	console    ui.Console
	manager    store.Manager
}

// PrintCommand represents an environment command
type PrintCommand interface {
	ExecuteCommand(params RunCommandParams) error
}

// NewPrintCommand creates a new environment command
func NewPrintCommand(
	manager store.Manager,
	fileSystem afero.Fs,
	platform string,
	console ui.Console) PrintCommand {
	return &printCommand{
		manager:    manager,
		fileSystem: fileSystem,
		platform:   platform,
		console:    console}
}

func (cmd *printCommand) ExecuteCommand(params RunCommandParams) error {

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
