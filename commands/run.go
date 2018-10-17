package commands

import (
	"github.com/pkg/errors"

	"github.com/patrickhuber/wrangle/processes"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/ui"

	"github.com/patrickhuber/wrangle/config"
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

	if processName == "" {
		return errors.New("process name is required for the run command")
	}

	cfg := params.Config()
	if cfg == nil {
		return errors.New("unable to load configuration")
	}

	processTemplate, err := store.NewProcessTemplate(cfg, cmd.manager)
	if err != nil {
		return err
	}
	process, err := processTemplate.Evaluate(processName)
	if err != nil {
		return err
	}

	return cmd.execute(process)
}

func (cmd *run) execute(processConfig *config.Process) error {
	process := cmd.processFactory.Create(
		processConfig.Path,
		processConfig.Args,
		processConfig.Vars,
		cmd.console.Out(),
		cmd.console.Error(),
		cmd.console.In())
	return process.Dispatch()
}
