package services_test

import (
	"github.com/spf13/afero"
	"github.com/patrickhuber/wrangle/services"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/ui"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Processes", func() {
	It("can list processes", func() {
		console := ui.NewMemoryConsole()
		fs := afero.NewMemMapFs()
		loader := config.NewLoader(fs)
		
		// write config file
		content := `
processes:
- name: go 
- name: echo
- name: wrangle
- name: dangle
`
		err := afero.WriteFile(fs, "/config", []byte(content), 0600)
		Expect(err).To(BeNil())

		// create the service
		service := services.NewProcessesService(console, loader)		
		Expect(service).ToNot(BeNil())

		err = service.List("/config")
		Expect(err).To(BeNil())

		value := console.OutAsString()
		Expect(value).To(Equal("name\n----\ngo\necho\nwrangle\ndangle\n"))
	})
})
