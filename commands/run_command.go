package commands

import (
	"errors"
	"fmt"

	"github.com/patrickhuber/cli-mgr/processes"
	"github.com/patrickhuber/cli-mgr/store"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/spf13/afero"
	"github.com/urfave/cli"
)

type RunCommand struct {
	configStoreManager store.Manager
	fileSystem         afero.Fs
	processFactory     processes.ProcessFactory
}

func NewRunCommand(
	configStoreManager store.Manager,
	fileSystem afero.Fs,
	processFactory processes.ProcessFactory) *RunCommand {
	return &RunCommand{
		configStoreManager: configStoreManager,
		fileSystem:         fileSystem,
		processFactory:     processFactory}
}

func (cmd *RunCommand) ExecuteCommand(c *cli.Context) error {
	configFile := c.GlobalString("config")
	processName := c.String("name")
	environmenName := c.String("environment")

	if processName == "" {
		return errors.New("process name is required for the run command")
	}

	if environmenName == "" {
		return errors.New("environment name is required for the run command")
	}

	configLoader := config.NewConfigLoader(cmd.fileSystem)
	cfg, err := configLoader.Load(configFile)
	if err != nil {
		return err
	}
	if cfg == nil {
		return errors.New("config is null")
	}
	return cmd.executeConfigItem(cfg, processName, environmenName)
}

func (cmd *RunCommand) executeConfigItem(cfg *config.Config, processName string, environmentName string) error {
	for _, p := range cfg.Processes {
		if p.Name == processName {
			for _, e := range p.Environments {
				if e.Name == environmentName {
					return cmd.execute(&e)
				}
			}
			return fmt.Errorf("unable to find environment '%s' in process '%s'", environmentName, processName)
		}
	}
	return fmt.Errorf("No Processes found in config that match '%s'", processName)
}

func (cmd *RunCommand) execute(processEnvironmentConfig *config.Environment) error {
	process := cmd.processFactory.Create(
		processEnvironmentConfig.Process,
		processEnvironmentConfig.Args,
		processEnvironmentConfig.Vars)
	return process.Dispatch()
}
