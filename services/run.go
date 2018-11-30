package services

import (
	"github.com/pkg/errors"

	"github.com/patrickhuber/wrangle/processes"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/ui"

	"github.com/patrickhuber/wrangle/config"
	"github.com/spf13/afero"
)

// RunService runs a process defined by params
type RunService interface {
	Run(configFile string, params ProcessParams) error
}

type runService struct {
	manager        store.Manager
	fileSystem     afero.Fs
	processFactory processes.Factory
	console        ui.Console
	loader config.Loader
}

// NewRunService - creates a new run command
func NewRunService(
	manager store.Manager,
	fileSystem afero.Fs,
	processFactory processes.Factory,
	console ui.Console,
	loader config.Loader) RunService{
	return &runService{
		manager:        manager,
		fileSystem:     fileSystem,
		processFactory: processFactory,
		console:        console,
		loader: loader}
}

func (service *runService) Run(configFile string, params ProcessParams) error {

	processName := params.ProcessName()

	if processName == "" {
		return errors.New("process name is required for the run command")
	}

	cfg, err := service.loader.LoadConfig(configFile)
	if err  != nil{
		return err
	}
	
	if cfg == nil {
		return errors.New("unable to load configuration")
	}

	processTemplate, err := store.NewProcessTemplate(cfg, service.manager)
	if err != nil {
		return err
	}
	process, err := processTemplate.Evaluate(processName)
	if err != nil {
		return err
	}

	return service.execute(process)
}

func (service *runService) execute(processConfig *config.Process) error {
	process := service.processFactory.Create(
		processConfig.Path,
		processConfig.Args,
		processConfig.Vars,
		service.console.Out(),
		service.console.Error(),
		service.console.In())
	return process.Dispatch()
}
