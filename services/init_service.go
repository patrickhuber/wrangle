package services

import "github.com/spf13/afero"

type initService struct {
	fileSystem afero.Fs
}

// InitService executes directory initialization
type InitService interface {
	Init(configFilePath string) error
}

// NewInitService cretes a new init command
func NewInitService(fileSystem afero.Fs) InitService {
	return &initService{
		fileSystem: fileSystem,
	}
}

func (i *initService) Init(configFilePath string) error {
	data := "stores: \nprocesses: \n"
	return afero.WriteFile(i.fileSystem, configFilePath, []byte(data), 0640)
}
