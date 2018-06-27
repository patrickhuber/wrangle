package commands

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/store"
	"github.com/patrickhuber/cli-mgr/ui"
	"github.com/spf13/afero"
)

type envCommand struct {
	fileSystem afero.Fs
	platform   string
	console    ui.Console
	manager    store.Manager
}

// EnvCommand represents an environment command
type EnvCommand interface {
	ExecuteCommand(params RunCommandParams) error
}

// NewEnvCommand creates a new environment command
func NewEnvCommand(
	manager store.Manager,
	fileSystem afero.Fs,
	platform string,
	console ui.Console) EnvCommand {
	return &envCommand{
		manager:    manager,
		fileSystem: fileSystem,
		platform:   platform,
		console:    console}
}

func (cmd *envCommand) ExecuteCommand(params RunCommandParams) error {

	configFile := params.ConfigFile()
	processName := params.ProcessName()
	environmentName := params.EnvironmentName()

	if processName == "" {
		return errors.New("process name is required for the run command")
	}

	if environmentName == "" {
		return errors.New("environment name is required for the run command")
	}

	configLoader := config.NewLoader(cmd.fileSystem)

	cfg, err := configLoader.Load(configFile)
	if err != nil {
		return errors.Wrap(err, "error running configLoader.Load")
	}
	pipeline := store.NewPipeline(cmd.manager, cfg)
	environment, err := pipeline.Run(processName, environmentName)
	if err != nil {
		return err
	}
	renderer := NewEnvVarRenderer(cmd.platform)
	fmt.Fprint(cmd.console.Out(), renderer.RenderEnvironment(environment.Vars))
	return nil
}
