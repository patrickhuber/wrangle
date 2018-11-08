package commands

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/renderers"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/store/file"

	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"
)

var _ = Describe("Execute", func() {
	var (
		expectedOutput     string
		includeProcessInfo bool
	)
	Describe("WithSimpleConfig", func() {
		AfterEach(func() {
			// create config file
			configFileContent := `
---
stores:
processes:
- name: echo
  path: echo
  env:
    WRANGLE_TEST: value`
			RunPrintTest(configFileContent, "", "echo", includeProcessInfo, expectedOutput)
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
			rendererFactory := renderers.NewFactory(collections.NewDictionary())

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
			cfg, err := loader.LoadConfig("/config")
			Expect(err).To(BeNil())

			// create and run command
			cmd := NewPrint(manager, fileSystem, console, rendererFactory)
			params := &PrintParams{
				Configuration: cfg,
				ProcessName:   "echo",
				Format:        "",
				Include: PrintParamsInclude{
					ProcessAndArgs: includeProcessInfo,
				}}
			err = cmd.Execute(params)
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
			content := `
processes:
- name: go
  path: go
  args: 
  - version
`
			RunPrintTest(content, "", "go", includeProcessInfo, expectedOutput)
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
			content := `
processes:
- name: go
  path: go
  args: 
  - version
`
			RunPrintTest(content, format, "go", false, expectedOutput)
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
	content string,
	format string,
	processName string,
	includeProcessInfo bool,
	expectedOutput string) {

	rendererFactory := renderers.NewFactory(collections.NewDictionary())

	// create store manager
	manager := store.NewManager()

	fileSystem := afero.NewMemMapFs()
	afero.WriteFile(fileSystem, "/config", []byte(content), 0444)
	console := ui.NewMemoryConsole()

	// load the config
	loader := config.NewLoader(fileSystem)
	cfg, err := loader.LoadConfig("/config")
	Expect(err).To(BeNil())

	// create and run command
	cmd := NewPrint(manager, fileSystem, console, rendererFactory)
	params := &PrintParams{
		Configuration: cfg,
		ProcessName:   processName,
		Format:        format,
		Include: PrintParamsInclude{
			ProcessAndArgs: includeProcessInfo,
		},
	}
	err = cmd.Execute(params)
	Expect(err).To(BeNil())

	// verify output
	b, ok := console.Out().(*bytes.Buffer)
	Expect(ok).To(BeTrue())
	Expect(b).ToNot(BeNil())
	Expect(b.String()).To(Equal(expectedOutput))
}
