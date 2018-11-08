package commands

import (
	"bytes"
	"testing"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/processes"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Execute", func() {
	It("", func() {
		Expect(true).To(BeTrue())
	})
})

func TestRunCommand(t *testing.T) {

	t.Run("CanRunGoVersionProcess", func(t *testing.T) {
		r := require.New(t)

		// write out the config file
		configFileData := `
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
		cfg, err := loader.LoadConfig("/config")
		r.Nil(err)

		// run the run command
		err = runCommand.Execute(
			NewProcessParams(cfg, "go"))
		r.Nil(err)
	})

	t.Run("CanRedirectStdOut", func(t *testing.T) {
		r := require.New(t)

		configFileData := `
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
		cfg, err := loader.LoadConfig("/config")
		r.Nil(err)

		// run the run command
		err = runCommand.Execute(
			NewProcessParams(cfg, "go"))
		r.Nil(err)

		// check something was written to stdout
		buffer := console.Out().(*bytes.Buffer)
		r.NotEmpty(buffer)
	})
}
