package packages

import (
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/tasks"
)

type manager struct {
	fileSystem    filesystem.FsWrapper
	taskProviders tasks.ProviderRegistry
}

// Manager defines a manager interface
type Manager interface {
	Install(p Package) error
}

// NewManager creates a new package manager
func NewManager(fileSystem filesystem.FsWrapper, taskProviders tasks.ProviderRegistry) Manager {
	return &manager{
		fileSystem:    fileSystem,
		taskProviders: taskProviders}
}

func (manager *manager) Install(p Package) error {
	for _, task := range p.Tasks() {
		provider, err := manager.taskProviders.Get(task.Type())
		if err != nil {
			return err
		}
		err = provider.Execute(task)
		if err != nil {
			return err
		}
	}
	return nil
}
