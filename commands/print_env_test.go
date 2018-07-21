package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/renderers"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/store/file"

	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"
)

func TestPrintEnvCommand(t *testing.T) {
	t.Run("CanRunCommand", func(t *testing.T) {
		r := require.New(t)

		// create filesystem
		fileSystem := afero.NewMemMapFs()

		// create config file
		configFileContent := `
---
stores:
environments:
- name: lab
  processes:
  - name: echo
    path: echo
    env:
      WRANGLE_TEST: value`
		afero.WriteFile(fileSystem, "/config", []byte(configFileContent), 0644)

		platform := "linux"
		// create renderer factory
		rendererFactory := renderers.NewFactory(platform, collections.NewDictionary())

		// create store manager
		manager := store.NewManager()
		manager.Register(file.NewFileStoreProvider(fileSystem, nil))

		// create console
		console := ui.NewMemoryConsole()

		// load the config
		loader := config.NewLoader(fileSystem)
		cfg, err := loader.Load("/config")
		r.Nil(err)

		// create and run command
		cmd := NewPrintEnv(manager, fileSystem, console, rendererFactory)
		params := &PrintEnvParams{
			Configuration:   cfg,
			EnvironmentName: "lab",
			ProcessName:     "echo",
			Shell:           ""}
		err = cmd.Execute(params)
		r.Nil(err)

		// verify output
		b, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(b)
		r.Equal("export WRANGLE_TEST=value\n", b.String())
	})

	t.Run("CanRunReplaceVariable", func(t *testing.T) {
		r := require.New(t)

		platform := "linux"
		rendererFactory := renderers.NewFactory(platform, collections.NewDictionary())

		// create filesystem
		fileSystem := afero.NewMemMapFs()

		// create config file
		configFileContent := `
---
stores:
- name: store1
  type: file
  params: 
    path: /store1
environments:
- name: lab
  processes:
  - name: echo
    path: echo
    stores: [ store1 ]
    env:
      WRANGLE_TEST: ((/key))`
		afero.WriteFile(fileSystem, "/config", []byte(configFileContent), 0644)
		afero.WriteFile(fileSystem, "/store1", []byte("key: value"), 0644)

		// create store manager
		manager := store.NewManager()
		manager.Register(file.NewFileStoreProvider(fileSystem, nil))

		// create console
		console := ui.NewMemoryConsole()

		// load the config
		loader := config.NewLoader(fileSystem)
		cfg, err := loader.Load("/config")
		r.Nil(err)

		// create and run command
		cmd := NewPrintEnv(manager, fileSystem, console, rendererFactory)
		params := &PrintEnvParams{
			Configuration:   cfg,
			EnvironmentName: "lab",
			ProcessName:     "echo",
			Shell:           ""}
		err = cmd.Execute(params)
		r.Nil(err)

		// verify output
		b, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(b)
		r.Equal("export WRANGLE_TEST=value\n", b.String())
	})

	t.Run("ShellOverridesPlatform", func(t *testing.T) {
		r := require.New(t)
		platform := "linux"
		rendererFactory := renderers.NewFactory(platform, collections.NewDictionary())

		// create filesystem
		fileSystem := afero.NewMemMapFs()

		// create config file
		configFileContent := `
---
environments:
- name: lab
  processes:
  - name: echo
    path: echo    
    env:
      WRANGLE_TEST: hello`

		afero.WriteFile(fileSystem, "/config", []byte(configFileContent), 0644)

		// create store manager
		manager := store.NewManager()
		manager.Register(file.NewFileStoreProvider(fileSystem, nil))

		// create console
		console := ui.NewMemoryConsole()

		// load the config
		loader := config.NewLoader(fileSystem)
		cfg, err := loader.Load("/config")
		r.Nil(err)

		// create and run command
		cmd := NewPrintEnv(manager, fileSystem, console, rendererFactory)
		params := &PrintEnvParams{
			Configuration:   cfg,
			EnvironmentName: "lab",
			ProcessName:     "echo",
			Shell:           "powershell"}
		err = cmd.Execute(params)
		r.Nil(err)

		// verify output
		b, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(b)
		r.Equal("$env:WRANGLE_TEST=\"hello\"\r\n", b.String())
	})
}
