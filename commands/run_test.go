package commands

import (
	"bytes"
	"testing"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/processes"
	"github.com/patrickhuber/cli-mgr/store"
	"github.com/patrickhuber/cli-mgr/ui"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestRunCommand(t *testing.T) {

	t.Run("CanRunGoVersionProcess", func(t *testing.T) {
		r := require.New(t)

		// write out the config file
		configFileData := `
environments:
- name: lab
  processes:
  - name: go
    path: go
    args:
    - version
`
		fileSystem := afero.NewMemMapFs()
		err := afero.WriteFile(fileSystem, "/config", []byte(configFileData), 0644)
		r.Nil(err)

		// create the console
		console := ui.NewMemoryConsole()

		// create run command
		configStoreManager := store.NewManager()
		runCommand := NewRun(configStoreManager, fileSystem, processes.NewOsFactory(), console)

		// load the config
		loader := config.NewLoader(fileSystem)
		cfg, err := loader.Load("/config")
		r.Nil(err)

		// run the run command
		err = runCommand.Execute(
			NewProcessParams(cfg, "lab", "go"))
		r.Nil(err)
	})

	t.Run("CanRedirectStdOut", func(t *testing.T) {
		r := require.New(t)

		configFileData := `
environments:
- name: lab
  processes:
  - name: go
    path: go
    args: 
    - version 
`
		fileSystem := afero.NewMemMapFs()
		err := afero.WriteFile(fileSystem, "/config", []byte(configFileData), 0644)
		r.Nil(err)

		// create the console
		console := ui.NewMemoryConsole()

		// create run command
		configStoreManager := store.NewManager()
		runCommand := NewRun(configStoreManager, fileSystem, processes.NewOsFactory(), console)

		// load the config
		loader := config.NewLoader(fileSystem)
		cfg, err := loader.Load("/config")
		r.Nil(err)

		// run the run command
		err = runCommand.Execute(
			NewProcessParams(cfg, "lab", "go"))
		r.Nil(err)

		// check something was written to stdout
		buffer := console.Out().(*bytes.Buffer)
		r.NotEmpty(buffer)
	})
}
