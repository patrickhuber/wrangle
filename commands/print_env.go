package commands

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/patrickhuber/wrangle/renderers"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"
)

type printEnv struct {
	fileSystem afero.Fs
	platform   string
	shell      string
	console    ui.Console
	manager    store.Manager
}

// PrintEnv represents an environment command
type PrintEnv interface {
	Execute(params ProcessParams) error
}

// NewPrintEnv creates a new environment command
func NewPrintEnv(
	manager store.Manager,
	fileSystem afero.Fs,
	platform string,
	shell string,
	console ui.Console) PrintEnv {
	return &printEnv{
		manager:    manager,
		fileSystem: fileSystem,
		platform:   platform,
		shell:      shell,
		console:    console}
}

func (cmd *printEnv) Execute(params ProcessParams) error {

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

	fmt.Fprint(cmd.console.Out(), renderer.RenderEnvironment(process.Vars))

	return nil
}
