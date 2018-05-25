package commands

import (
	"errors"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/spf13/afero"
	"github.com/urfave/cli"
)

type RunCommand struct {
	configStoreManager *config.ConfigStoreManager
	fileSystem         afero.Fs
}

func NewRunCommand(configStoreManager *config.ConfigStoreManager, fileSystem afero.Fs) *RunCommand {
	return &RunCommand{
		configStoreManager: configStoreManager,
		fileSystem:         fileSystem}
}

func (cmd *RunCommand) ExecuteRunCommand(c *cli.Context) error {
	configFile := c.GlobalString("config")
	configLoader := config.ConfigLoader{FileSystem: cmd.fileSystem}
	cfg, err := configLoader.Load(configFile)
	if err != nil {
		return err
	}
	if cfg == nil {
		return errors.New("cfg is null")
	}
	return nil
}
