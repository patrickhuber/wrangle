package commands

import (
	"errors"
	"fmt"

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
	environmenName := params.EnvironmentName()

	if processName == "" {
		return errors.New("process name is required for the run command")
	}

	if environmenName == "" {
		return errors.New("environment name is required for the run command")
	}

	configLoader := config.NewConfigLoader(cmd.fileSystem)
	cfg, err := configLoader.Load(configFile)
	if err != nil {
		return err
	}
	if cfg == nil {
		return errors.New("config is null")
	}
	variables, err := cmd.getProcessEnvironmentVariables(cfg, processName, environmenName)
	if err != nil {
		return err
	}
	renderer := NewEvnVarRenderer(cmd.platform)
	fmt.Fprint(cmd.console.Out(), renderer.RenderEnvironment(variables))
	return nil
}

func (cmd *envCommand) getProcessEnvironmentVariables(cfg *config.Config, processName string, environmentName string) (map[string]string, error) {
	for _, p := range cfg.Processes {
		if p.Name == processName {
			for _, e := range p.Environments {
				if e.Name == environmentName {
					return e.Vars, nil
				}
			}
			return nil, fmt.Errorf("unable to find environment '%s' in process '%s'", environmentName, processName)
		}
	}
	return nil, fmt.Errorf("No Processes found in config that match '%s'", processName)
}
