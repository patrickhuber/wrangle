package processes

import (
	"errors"
	"fmt"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/renderers"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/ui"
)

// PrintService prints out the process information
type PrintService interface {
	Print(params *PrintParams) error
}

type printService struct {
	manager         store.Manager
	rendererFactory renderers.Factory
	console         ui.Console
}

// NewPrintService creates a new print service with the given configuration
func NewPrintService(console ui.Console, manager store.Manager, rendererFactory renderers.Factory) PrintService {
	return &printService{
		manager:         manager,
		console:         console,
		rendererFactory: rendererFactory,
	}
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

func (s *printService) Print(params *PrintParams) error {

	processName := params.ProcessName

	if processName == "" {
		return errors.New("process name is required for the print command")
	}

	cfg := params.Config
	processTemplate, err := store.NewProcessTemplate(cfg, s.manager)
	if err != nil {
		return err
	}

	process, err := processTemplate.Evaluate(processName)
	if err != nil {
		return err
	}

	renderer, err := s.rendererFactory.Create(params.Format)
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
	_, err = fmt.Fprint(s.console.Out(), renderedOutput)
	return err
}
