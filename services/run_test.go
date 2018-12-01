package services_test

import (
	"github.com/patrickhuber/wrangle/services"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/processes"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Execute", func() {
	var (
		fs afero.Fs 
		runService services.RunService
		console ui.MemoryConsole
	)
	BeforeEach(func(){
		fs = afero.NewMemMapFs()

		// write out the config file
		configFileData := `
processes:
- name: go
  path: go
  args:
  - version
`
		err := afero.WriteFile(fs, "/config", []byte(configFileData), 0644)		
		Expect(err).To(BeNil())

		// create the console
		console = ui.NewMemoryConsole()
		
		// load the config
		loader := config.NewLoader(fs)
		
		// create run command
		configStoreManager := store.NewManager()
		runService = services.NewRunService(configStoreManager, fs, processes.NewOsFactory(), console, loader)
		
		Expect(err).To(BeNil())
		
	})
	It("can run go version process", func() {	

		// run the run command
		err := runService.Run("/config", services.NewProcessParams("go"))
		Expect(err).To(BeNil())
	})

	It("can redirect to std out", func(){		
		// run the run command
		err := runService.Run("/config", services.NewProcessParams("go"))
		Expect(err).To(BeNil())		
		
		// check something was written to stdout
		buffer := console.OutAsString()
		Expect(buffer).ToNot(BeEmpty())
	})
})

