package services_test

import (
	"bytes"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/processes"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Execute", func() {
	It("can run go version process", func() {
		
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
		Expect(err).To(BeNil())

		// create the console
		console := ui.NewMemoryConsole()

		// create run command
		configStoreManager := store.NewManager()
		runCommand := NewRun(configStoreManager, fileSystem, processes.NewOsFactory(), console)

		// load the config
		loader := config.NewLoader(fileSystem)
		cfg, err := loader.LoadConfig("/config")

		Expect(err).To(BeNil())

		
		// run the run command
		err = runCommand.Execute(NewProcessParams(cfg, "go"))

		Expect(err).To(BeNil())
	})

	It("can redirect to std out", func(){
		configFileData := `
		processes:
		- name: go
		  path: go
		  args: 
		  - version 
		`
		fileSystem := afero.NewMemMapFs()
		err := afero.WriteFile(fileSystem, "/config", []byte(configFileData), 0644)
		
		Expect(err).To(BeNil())

		// create the console
		console := ui.NewMemoryConsole()

		// create run command
		configStoreManager := store.NewManager()
		runCommand := services.NewRunService(configStoreManager, fileSystem, processes.NewOsFactory(), console)

		// load the config
		loader := config.NewLoader(fileSystem)
		cfg, err := loader.LoadConfig("/config")

				
		Expect(err).To(BeNil())
		
		// run the run command
		err = runCommand.Execute(
			NewProcessParams(cfg, "go"))

			
		Expect(err).To(BeNil())		
		
		// check something was written to stdout
		buffer := console.Out().(*bytes.Buffer)
		Expect(buffer).ToNot(BeEmpty())
	})
})

