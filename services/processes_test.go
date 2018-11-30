package services_test

import (
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/ui"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Processes", func() {
	It("can list processes", func() {
		console := ui.NewMemoryConsole()

		cmd := servcies.NewProcessesService(console)
		content := `
processes:
- name: go 
- name: echo
- name: wrangle
- name: dangle
`
		cfg, err := config.DeserializeConfigString(content)
		Expect(err).To(BeNil())

		err = cmd.Execute(cfg)
		Expect(err).To(BeNil())

		value := console.OutAsString()
		Expect(value).To(Equal("name\n----\ngo\necho\nwrangle\ndangle\n"))
	})
})
