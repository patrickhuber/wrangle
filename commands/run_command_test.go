package commands

import (
	"flag"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/patrickhuber/cli-mgr/config"
)

func TestRunCommand(t *testing.T) {

	t.Run("CanRunGoVersionProcess", func(t *testing.T) {
		r := require.New(t)

		// write out the config file
		configFileData := `
processes:
- name: go
  environments:
  - name: lab
    process: go
    args:
    - -v
`
		fileSystem := afero.NewMemMapFs()
		err := afero.WriteFile(fileSystem, "/config", []byte(configFileData), 0644)
		r.Nil(err)

		// create run command
		configStoreManager := config.NewConfigStoreManager()
		runCommand := NewRunCommand(configStoreManager, fileSystem)

		// global context
		globalSet := flag.NewFlagSet("global", 0)
		globalSet.String("config", "/config", "")
		globalContext := cli.NewContext(nil, globalSet, nil)

		// context
		commandSet := flag.NewFlagSet("cmd", 0)
		commandSet.String("name", "go", "")
		commandSet.String("environment", "lab", "")
		context := cli.NewContext(nil, commandSet, globalContext)

		// run the run command
		err = runCommand.ExecuteCommand(context)
		r.Nil(err)
	})
}
