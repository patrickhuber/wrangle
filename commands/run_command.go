package commands

import (
	"github.com/pkg/errors"

	"github.com/patrickhuber/cli-mgr/processes"
	"github.com/patrickhuber/cli-mgr/store"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/spf13/afero"
)

// RunCommand represents a run subcommand for the application
type RunCommand interface {
	ExecuteCommand(params RunCommandParams) error
}

type runCommand struct {
	manager        store.Manager
	fileSystem     afero.Fs
	processFactory processes.ProcessFactory
}

// NewRunCommand - creates a new run command
func NewRunCommand(
	manager store.Manager,
	fileSystem afero.Fs,
	processFactory processes.ProcessFactory) RunCommand {
	return &runCommand{
		manager:        manager,
		fileSystem:     fileSystem,
		processFactory: processFactory}
}

func (cmd *runCommand) ExecuteCommand(params RunCommandParams) error {

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

	return cmd.execute(environment)
}

func (cmd *runCommand) execute(processEnvironmentConfig *config.Environment) error {
	process := cmd.processFactory.Create(
		processEnvironmentConfig.Process,
		processEnvironmentConfig.Args,
		processEnvironmentConfig.Vars)
	return process.Dispatch()
}
