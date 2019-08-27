package services_test

import (
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/processes"
	"github.com/patrickhuber/wrangle/services"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/ui"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Execute", func() {
	var (
		fs         filesystem.FileSystem
		runService services.RunService
		console    ui.MemoryConsole
		cfg        *config.Config
	)
	BeforeEach(func() {
		fs = filesystem.NewMemory()

		cfg = &config.Config{
			Processes: []config.Process{
				config.Process{
					Name: "go",
					Path: "go",
					Args: []string{"version"},
				},
			},
		}
		// create the console
		console = ui.NewMemoryConsole()

		// create run command
		configStoreManager := store.NewManager()
		runService = services.NewRunService(configStoreManager, fs, processes.NewOsFactory(), console)
	})
	It("can run go version process", func() {

		// run the run command
		err := runService.Run(services.NewProcessParams("go", cfg))
		Expect(err).To(BeNil())
	})

	It("can redirect to std out", func() {
		// run the run command
		err := runService.Run(services.NewProcessParams("go", cfg))
		Expect(err).To(BeNil())

		// check something was written to stdout
		buffer := console.OutAsString()
		Expect(buffer).ToNot(BeEmpty())
	})
})
