package processes_test

import (
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/processes"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/store/file"
	"github.com/patrickhuber/wrangle/ui"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Processes", func() {
	It("can list processes", func() {
		console := ui.NewMemoryConsole()

		cfg := &config.Config{
			Processes: []config.Process{
				config.Process{
					Name: "go",
				},
				config.Process{
					Name: "echo",
				},
				config.Process{
					Name: "wrangle",
				},
				config.Process{
					Name: "dangle",
				},
			},
		}

		// create the filesystem
		fileSystem := filesystem.NewMemory()

		// create store manager
		manager := store.NewManager()
		manager.Register(file.NewFileStoreProvider(fileSystem, nil))

		// create the service
		service := processes.NewService(console)
		Expect(service).ToNot(BeNil())

		err := service.List(cfg)
		Expect(err).To(BeNil())

		value := console.OutAsString()
		Expect(value).To(Equal("name\n----\ngo\necho\nwrangle\ndangle\n"))
	})
})
