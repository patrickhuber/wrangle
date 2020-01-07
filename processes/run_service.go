package processes

import (
	"errors"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/ui"
)

// RunService runs a process defined by params
type RunService interface {
	Run(params ProcessParams) error
}

type runService struct {
	manager        store.Manager
	fileSystem     filesystem.FileSystem
	processFactory Factory
	console        ui.Console
}

// NewRunService - creates a new run command
func NewRunService(
	manager store.Manager,
	fileSystem filesystem.FileSystem,
	processFactory Factory,
	console ui.Console) RunService {
	return &runService{
		manager:        manager,
		fileSystem:     fileSystem,
		processFactory: processFactory,
		console:        console}
}

func (service *runService) Run(params ProcessParams) error {

	processName := params.ProcessName()

	if processName == "" {
		return errors.New("process name is required for the run command")
	}

	cfg := params.Config()
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
