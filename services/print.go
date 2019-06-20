package services

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/renderers"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"
)

type printService struct {
	fileSystem      afero.Fs
	console         ui.Console
	manager         store.Manager
	rendererFactory renderers.Factory
}

// PrintParamsInclude defines what additional output to include
type PrintParamsInclude struct {
	ProcessAndArgs bool
}

// PrintParams defines parameters for the print command
type PrintParams struct {
	Config      *config.Config
	ProcessName string
	Format      string
	Include     PrintParamsInclude
}

// PrintService represents an environment command
type PrintService interface {
	Print(params *PrintParams) error
}

// NewPrintService creates a new environment command
func NewPrintService(
	manager store.Manager,
	fileSystem afero.Fs,
	console ui.Console,
	rendererFactory renderers.Factory) PrintService {
	return &printService{
		manager:         manager,
		fileSystem:      fileSystem,
		console:         console,
		rendererFactory: rendererFactory,
	}
}

func (service *printService) Print(params *PrintParams) error {

	processName := params.ProcessName

	if processName == "" {
		return errors.New("process name is required for the print command")
	}

	cfg := params.Config
	processTemplate, err := store.NewProcessTemplate(cfg, service.manager)
	if err != nil {
		return err
	}

	process, err := processTemplate.Evaluate(processName)
	if err != nil {
		return err
	}

	renderer, err := service.rendererFactory.Create(params.Format)
	if err != nil {
		return err
	}

	var renderedOutput string
	if params.Include.ProcessAndArgs {
		renderedOutput = renderer.RenderProcess(
			process.Path,
			process.Args,
			process.Vars)
	} else {
		renderedOutput = renderer.RenderEnvironment(process.Vars)
	}
	_, err = fmt.Fprint(service.console.Out(), renderedOutput)
	return err
}
