package main

import (
	"bytes"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/patrickhuber/cli-mgr/config"
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

		config := config.Config{
			ConfigSources: []config.ConfigSource{
				{
					Name:             "test",
					ConfigSourceType: "file",
					Params: map[string]string{
						"path": "/store1",
					},
				},
				{
					Name:             "test",
					ConfigSourceType: "file",
					Params: map[string]string{
						"path": "/store2",
					},
				},
			},
			Processes: []config.Process{
				{
					Name: "echo",
					Environments: []config.Environment{
						{
							Name:    "lab",
							Config:  "store1",
							Process: "echo",
							Vars: map[string]string{
								"CLI_MGR_TEST": "((key))",
							},
						},
					},
				},
			},
		}
		configFileContent, err := yaml.Marshal(config)
		r.Nil(err)

		// create files
		err = afero.WriteFile(fileSystem, "/config", configFileContent, 0644)
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
