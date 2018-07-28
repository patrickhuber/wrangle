package commands

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/renderers"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"
)

type printEnv struct {
	fileSystem      afero.Fs
	rendererFactory renderers.Factory
	console         ui.Console
	manager         store.Manager
}

// PrintEnvParams defines parameters for the print env command
type PrintEnvParams struct {
	Configuration   *config.Config
	EnvironmentName string
	ProcessName     string
	Format          string
}

// PrintEnv represents an environment command
type PrintEnv interface {
	Execute(*PrintEnvParams) error
}

// NewPrintEnv creates a new environment command
func NewPrintEnv(
	manager store.Manager,
	fileSystem afero.Fs,
	console ui.Console,
	rendererFactory renderers.Factory) PrintEnv {
	return &printEnv{
		manager:         manager,
		fileSystem:      fileSystem,
		rendererFactory: rendererFactory,
		console:         console}
}

func (cmd *printEnv) Execute(
	params *PrintEnvParams) error {

	if params.ProcessName == "" {
		return errors.New("process name is required for the run command")
	}
	processName := params.ProcessName

	if params.EnvironmentName == "" {
		return errors.New("environment name is required for the run command")
	}
	environmentName := params.EnvironmentName

	if params.Configuration == nil {
		return errors.New("unable to load configuration")
	}
	cfg := params.Configuration

	processTemplate, err := store.NewProcessTemplate(cfg, cmd.manager)
	if err != nil {
		return err
	}

	process, err := processTemplate.Evaluate(environmentName, processName)
	if err != nil {
		return err
	}

	renderer, err := cmd.rendererFactory.Create(params.Format)
	if err != nil {
		return err
	}

	fmt.Fprint(cmd.console.Out(), renderer.RenderEnvironment(process.Vars))

	return nil
}
