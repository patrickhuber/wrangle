package commands_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/commands"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"
)

var _ = Describe("Environments", func() {
	It("can get list of environments", func() {
		fileSystem := afero.NewMemMapFs()
		content := `
environments:
- name: one
- name: two
- name: three
`
		afero.WriteFile(fileSystem, "/test", []byte(content), 0644)

		console := ui.NewMemoryConsole()
		command := commands.NewEnvironments(fileSystem, console)
		loader := config.NewLoader(fileSystem)
		configuration, err := loader.Load("/test")
		Expect(err).To(BeNil())

		err = command.Execute(configuration)
		Expect(err).To(BeNil())

		b, ok := console.Out().(*bytes.Buffer)
		Expect(ok).To(BeTrue())
		Expect(b).ToNot(BeNil())
		Expect(b.String()).To(Equal("name\n----\none\ntwo\nthree\n"))
	})
})
