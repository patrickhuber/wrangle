package commands

import (
	"github.com/pkg/errors"

	"github.com/patrickhuber/cli-mgr/processes"
	"github.com/patrickhuber/cli-mgr/store"
	"github.com/patrickhuber/cli-mgr/ui"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/spf13/afero"
)

// Run represents a run subcommand for the application
type Run interface {
	Execute(params ProcessParams) error
}

type run struct {
	manager        store.Manager
	fileSystem     afero.Fs
	processFactory processes.Factory
	console        ui.Console
}

// NewRun - creates a new run command
func NewRun(
	manager store.Manager,
	fileSystem afero.Fs,
	processFactory processes.Factory,
	console ui.Console) Run {
	return &run{
		manager:        manager,
		fileSystem:     fileSystem,
		processFactory: processFactory,
		console:        console}
}

func (cmd *run) Execute(params ProcessParams) error {

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

	return cmd.execute(environment)
}

func (cmd *run) execute(processConfig *config.Process) error {
	process := cmd.processFactory.Create(
		processConfig.Path,
		processConfig.Args,
		processConfig.Vars,
		cmd.console.Out(),
		cmd.console.Error())
	return process.Dispatch()
}
