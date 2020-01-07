package processes_test

import (
	"bytes"

	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/processes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/renderers"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/store/file"

	"github.com/patrickhuber/wrangle/ui"
)

var _ = Describe("Execute", func() {
	var (
		expectedOutput     string
		includeProcessInfo bool
	)
	Describe("WithSimpleConfig", func() {
		AfterEach(func() {
			cfg := &config.Config{
				Processes: []config.Process{
					config.Process{
						Name: "echo",
						Path: "echo",
						Vars: map[string]string{"WRANGLE_TEST": "value"},
					},
				},
			}
			RunPrintTest(cfg, "", "echo", includeProcessInfo, expectedOutput)
		})
		Context("WhenNotIncludeProcessAndArgs", func() {
			It("prints only environment", func() {
				includeProcessInfo = false
				expectedOutput = "export WRANGLE_TEST=value\n"
			})
		})
		Context("WhenIncludeProcessAndArgs", func() {
			It("prints process and environment", func() {
				includeProcessInfo = true
				expectedOutput = "export WRANGLE_TEST=value\necho\n"
			})
		})
	})
	Describe("WithStore", func() {
		var (
			expectedOutput     string
			includeProcessInfo bool
		)
		AfterEach(func() {
			// create renderer factory
			rendererFactory := renderers.NewFactory(collections.NewDictionary())

			// create filesystem
			fileSystem := filesystem.NewMemory()

			// create config
			cfg := &config.Config{
				Stores: []config.Store{
					config.Store{
						Name:      "store1",
						StoreType: "file",
						Params:    map[string]string{"path": "/store1"},
					},
				},
				Processes: []config.Process{
					config.Process{
						Name:   "echo",
						Path:   "echo",
						Stores: []string{"store1"},
						Vars:   map[string]string{"WRANGLE_TEST": "((/key))"},
					},
				},
			}
			fileSystem.Write("/store1", []byte("key: value"), 0644)

			// create store manager
			manager := store.NewManager()
			manager.Register(file.NewFileStoreProvider(fileSystem, nil))

			// create console
			console := ui.NewMemoryConsole()

			// create and run command
			service := processes.NewPrintService(console, manager, rendererFactory)
			params := &processes.PrintParams{
				Config:      cfg,
				ProcessName: "echo",
				Format:      "",
				Include: processes.PrintParamsInclude{
					ProcessAndArgs: includeProcessInfo,
				}}
			err := service.Print(params)
			Expect(err).To(BeNil())

			// verify output
			b, ok := console.Out().(*bytes.Buffer)
			Expect(ok).To(BeTrue())
			Expect(b).ToNot(BeNil())
			output := b.String()
			Expect(output).To(Equal(expectedOutput))
		})
		Context("WhenPrintOnlyEnvironment", func() {
			It("prints environment only", func() {
				includeProcessInfo = false
				expectedOutput = "export WRANGLE_TEST=value\n"
			})
		})
		Context("WhenPrintProcess", func() {
			It("prints process and env", func() {
				includeProcessInfo = true
				expectedOutput = "export WRANGLE_TEST=value\necho\n"
			})
		})
	})
	Describe("WithArgs", func() {
		var (
			expectedOutput     string
			includeProcessInfo bool
		)
		AfterEach(func() {
			cfg := &config.Config{
				Processes: []config.Process{
					config.Process{
						Name: "go",
						Path: "go",
						Args: []string{"version"},
					},
				},
			}
			RunPrintTest(cfg, "", "go", includeProcessInfo, expectedOutput)
		})
		Context("WhenPrintOnlyEnvironment", func() {
			It("prints nothing", func() {
				includeProcessInfo = false
				expectedOutput = ""
			})
		})
		Context("WhenPrintProcess", func() {
			It("prints process and arg", func() {
				includeProcessInfo = true
				expectedOutput = "go version\n"
			})
		})
	})
	Describe("WithFormat", func() {
		var (
			format         string
			expectedOutput string
		)
		AfterEach(func() {
			cfg := &config.Config{
				Processes: []config.Process{
					config.Process{
						Name: "go",
						Path: "go",
						Args: []string{"version"},
					},
				},
			}
			RunPrintTest(cfg, format, "go", false, expectedOutput)
		})
		Context("WhenFormatPosix", func() {

			Context("WhenPrintOnlyEnvironment", func() {
				It("prints env", func() {
					format = renderers.PosixFormat
					expectedOutput = ""
				})
			})
			Context("WhenPrintProcess", func() {
				It("prints process and env", func() {
					format = renderers.PosixFormat
					expectedOutput = ""
				})
			})
		})
		Context("WhenFormatPowershell", func() {
			Context("WhenPrintOnlyEnvironment", func() {
				It("prints env", func() {
					format = renderers.PowershellFormat
					expectedOutput = ""
				})
			})
			Context("WhenPrintProcess", func() {
				It("prints process and env", func() {
					format = renderers.PosixFormat
					expectedOutput = ""
				})
			})
		})
	})
})

func RunPrintTest(
	cfg *config.Config,
	format string,
	processName string,
	includeProcessInfo bool,
	expectedOutput string) {

	rendererFactory := renderers.NewFactory(collections.NewDictionary())

	// create store manager
	manager := store.NewManager()

	fileSystem := filesystem.NewMemory()
	console := ui.NewMemoryConsole()

	fileSystem.Write("/store1", []byte("key: value"), 0644)

	// create and run command
	service := processes.NewPrintService(console, manager, rendererFactory)
	params := &processes.PrintParams{
		Config:      cfg,
		ProcessName: processName,
		Format:      format,
		Include: processes.PrintParamsInclude{
			ProcessAndArgs: includeProcessInfo,
		},
	}
	err := service.Print(params)
	Expect(err).To(BeNil())

	// verify output
	b, ok := console.Out().(*bytes.Buffer)
	Expect(ok).To(BeTrue())
	Expect(b).ToNot(BeNil())
	Expect(b.String()).To(Equal(expectedOutput))
}
