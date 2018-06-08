package commands

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/spf13/afero"
	"github.com/urfave/cli"
)

type EnvCommand struct {
	fileSystem afero.Fs
}

func NewEnvCommand(fileSystem afero.Fs) *EnvCommand {
	return &EnvCommand{
		fileSystem: fileSystem}
}

func (cmd *EnvCommand) ExecuteCommand(c *cli.Context) error {

	configFile := c.GlobalString("config")
	processName := c.String("name")
	environmenName := c.String("environment")

	if processName == "" {
		return errors.New("process name is required for the run command")
	}

	if environmenName == "" {
		return errors.New("environment name is required for the run command")
	}

	configLoader := config.ConfigLoader{FileSystem: cmd.fileSystem}
	cfg, err := configLoader.Load(configFile)
	if err != nil {
		return err
	}
	if cfg == nil {
		return errors.New("config is null")
	}
	variables, err := cmd.getProcessEnvironmentVariables(cfg, processName, environmenName)
	if err != nil {
		return err
	}
	renderer := NewEvnVarRenderer(runtime.GOOS)
	fmt.Println(renderer.RenderEnvironment(variables))
	return nil
}

func (cmd *EnvCommand) getProcessEnvironmentVariables(cfg *config.Config, processName string, environmentName string) (map[string]string, error) {
	for _, p := range cfg.Processes {
		if p.Name == processName {
			for _, e := range p.Environments {
				if e.Name == environmentName {
					return e.Vars, nil
				}
			}
			return nil, fmt.Errorf("unable to find environment '%s' in process '%s'", environmentName, processName)
		}
	}
	return nil, fmt.Errorf("No Processes found in config that match '%s'", processName)
}
