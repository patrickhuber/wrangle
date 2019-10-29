package services

import (
	"github.com/patrickhuber/wrangle/filesystem"
)

type initService struct {
	fileSystem filesystem.FileSystem
}

// InitService executes directory initialization
type InitService interface {
	Init(configFilePath string) error
}

// NewInitService cretes a new init command
func NewInitService(fileSystem filesystem.FileSystem) InitService {
	return &initService{
		fileSystem: fileSystem,
	}
}

// Init creates the config file as well as the settings file
func (i *initService) Init(configFilePath string) error {
	data := "stores: \nprocesses: \n"
	return i.fileSystem.Write(configFilePath, []byte(data), 0640)
}
