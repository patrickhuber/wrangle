package commands

import (
	"bytes"
	"os"
	"strings"
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

func TestMain(m *testing.M) {
	os.Unsetenv("PSModulePath")
	m.Run()
}

func TestPrintCommand(t *testing.T) {

	t.Run("CanRunCommand", func(t *testing.T) {
		r := require.New(t)
		platform := "linux"

		rendererFactory := renderers.NewFactory(platform, collections.NewDictionary())

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
		configFileContent = strings.Replace(configFileContent, "\t", "  ", -1)
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
		cmd := NewPrint(manager, fileSystem, console, rendererFactory)
		params := &PrintParams{
			Configuration:   cfg,
			EnvironmentName: "lab",
			ProcessName:     "echo"}
		err = cmd.Execute(params)
		r.Nil(err)

		// verify output
		b, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(b)
		r.Equal("export WRANGLE_TEST=value\necho\n", b.String())
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
      WRANGLE_TEST: ((key))`
		configFileContent = strings.Replace(configFileContent, "\t", "  ", -1)
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
		cmd := NewPrint(manager, fileSystem, console, rendererFactory)
		params := &PrintParams{
			Configuration:   cfg,
			EnvironmentName: "lab",
			ProcessName:     "echo"}
		err = cmd.Execute(params)
		r.Nil(err)

		// verify output
		b, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(b)
		r.Equal("export WRANGLE_TEST=value\necho\n", b.String())
	})

	t.Run("CanPrintOutArgumentsInCommand", func(t *testing.T) {
		r := require.New(t)
		content := `
environments:
- name: lab
  processes:
  - name: go
    path: go
    args: 
    - version
`
		platform := "linux"

		rendererFactory := renderers.NewFactory(platform, collections.NewDictionary())

		// create store manager
		manager := store.NewManager()

		fileSystem := afero.NewMemMapFs()
		afero.WriteFile(fileSystem, "/config", []byte(content), 0444)
		console := ui.NewMemoryConsole()

		// load the config
		loader := config.NewLoader(fileSystem)
		cfg, err := loader.Load("/config")
		r.Nil(err)

		// create and run command
		cmd := NewPrint(manager, fileSystem, console, rendererFactory)
		params := &PrintParams{
			Configuration:   cfg,
			EnvironmentName: "lab",
			ProcessName:     "go"}
		err = cmd.Execute(params)
		r.Nil(err)

		// verify output
		b, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(b)
		r.Equal("go version\n", b.String())

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
		cmd := NewPrint(manager, fileSystem, console, rendererFactory)
		params := &PrintParams{
			Configuration:   cfg,
			EnvironmentName: "lab",
			ProcessName:     "echo",
			Shell:           "powershell",
		}
		err = cmd.Execute(params)
		r.Nil(err)

		// verify output
		b, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(b)
		r.Equal("$env:WRANGLE_TEST=\"hello\"\r\necho\r\n", b.String())
	})

}
