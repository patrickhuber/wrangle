package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/patrickhuber/cli-mgr/processes"
	"github.com/patrickhuber/cli-mgr/ui"
	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/cli-mgr/store"

	"github.com/spf13/afero"
)

func TestMain(t *testing.T) {
	t.Run("CanRunProcess", func(t *testing.T) {

	})
	t.Run("CanChainConfigStores", func(t *testing.T) {
		r := require.New(t)

		// create dependencies
		fileSystem := afero.NewMemMapFs()
		storeManager := store.NewManager()
		processFactory := processes.NewOsProcessFactory() // change to fake process factory?
		console := ui.NewMemoryConsole()

		// create config file
		configFileContent := `
---
config-sources:
- name: store1
  type: file
  path: /store1
- name: store2
  type: file
  path: /store2
processes:
- name: echo
  environments:
  - name: lab
    process: echo
    env:
      CLI_MGR_TEST: ((key1))`
		configFileContent = strings.Replace(configFileContent, "\t", "  ", -1)

		// create files
		err := afero.WriteFile(fileSystem, "/config", []byte(configFileContent), 0644)
		r.Nil(err)

		err = afero.WriteFile(fileSystem, "/store1", []byte("key: ((key1))"), 0644)
		r.Nil(err)

		err = afero.WriteFile(fileSystem, "/store2", []byte("key1: value"), 0644)
		r.Nil(err)

		// create cli
		app, err := createApplication(
			storeManager,
			fileSystem,
			processFactory,
			console,
			"linux")
		r.Nil(err)

		// run command
		args := []string{
			"cli-mgr",
			"-c", "/config",
			"env",
			"-n", "echo",
			"-e", "lab"}
		err = app.cliApplication.Run(args)
		r.Nil(err)

		// get the output, validate the chaining works
		buffer, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(buffer)

		r.Equal("export CLI_MGR_TEST=value\n", buffer.String())
	})
}
