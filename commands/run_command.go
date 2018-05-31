package commands

import (
	"errors"
	"fmt"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/spf13/afero"
	"github.com/urfave/cli"
)

type RunCommand struct {
	configStoreManager *config.ConfigStoreManager
	fileSystem         afero.Fs
	process            Command
}

func NewRunCommand(
	configStoreManager *config.ConfigStoreManager,
	fileSystem afero.Fs) *RunCommand {
	return &RunCommand{
		configStoreManager: configStoreManager,
		fileSystem:         fileSystem}
}

func (cmd *RunCommand) ExecuteRunCommand(c *cli.Context) error {
	configFile := c.GlobalString("config")
	processName := c.String("name")
	environmenName := c.String("environment")

	configLoader := config.ConfigLoader{FileSystem: cmd.fileSystem}
	cfg, err := configLoader.Load(configFile)
	if err != nil {
		return err
	}
	if cfg == nil {
		return errors.New("cfg is null")
	}

	return executeConfigItem(cfg, processName, environmenName)
}

func executeConfigItem(cfg *config.Config, processName string, environmentName string) error {
	for _, p := range cfg.Processes {
		if p.Name == processName {
			for _, e := range p.Environments {
				if e.Name == environmentName {
					return execute(&p)
				}
			}
			return fmt.Errorf("unable to find environment '%s' in process '%s'", environmentName, processName)
		}
	}
	return nil
}

func execute(process *config.Process) error {
	return nil
}
