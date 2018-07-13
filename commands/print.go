package commands

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/patrickhuber/wrangle/renderers"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"
)

type print struct {
	fileSystem afero.Fs
	platform   string
	shell      string
	console    ui.Console
	manager    store.Manager
}

// Print represents an environment command
type Print interface {
	Execute(params ProcessParams) error
}

// NewPrint creates a new environment command
func NewPrint(
	manager store.Manager,
	fileSystem afero.Fs,
	platform string,
	shell string,
	console ui.Console) PrintEnv {
	return &print{
		manager:    manager,
		fileSystem: fileSystem,
		platform:   platform,
		console:    console}
}

func (cmd *print) Execute(params ProcessParams) error {

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

	processTemplate, err := store.NewProcessTemplate(cfg, cmd.manager)
	if err != nil {
		return err
	}

	process, err := processTemplate.Evaluate(environmentName, processName)
	if err != nil {
		return err
	}

	renderer, err := renderers.NewFactory().Create(cmd.shell, cmd.platform)
	if err != nil {
		return err
	}

	fmt.Fprint(
		cmd.console.Out(),
		renderer.RenderProcess(
			process.Path,
			process.Args,
			process.Vars))
	return nil
}
