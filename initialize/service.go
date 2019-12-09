package initialize

import (
	"github.com/patrickhuber/wrangle/filesystem"
)

type service struct {
	fileSystem filesystem.FileSystem
}

// Service executes directory initialization
type Service interface {
	Init(configFilePath string) error
}

// NewService cretes a new init command
func NewService(fileSystem filesystem.FileSystem) Service {
	return &service{
		fileSystem: fileSystem,
	}
}

// Init creates the config file as well as the settings file
func (i *service) Init(configFilePath string) error {
	data := "stores: \nprocesses: \n"
	return i.fileSystem.Write(configFilePath, []byte(data), 0640)
}
