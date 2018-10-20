package commands

import (
	"github.com/spf13/afero"
)

type initCommand struct {
	fileSystem afero.Fs
}

// Init defines the init command
type Init interface {
	Execute(configFilePath string) error
}

// NewInitCommand cretes a new init command
func NewInitCommand(fileSystem afero.Fs) Init {
	return &initCommand{
		fileSystem: fileSystem,
	}
}

func (i *initCommand) Execute(configFilePath string) error {
	data := "stores: \nprocesses: \n"
	return afero.WriteFile(i.fileSystem, configFilePath, []byte(data), 0640)
}
