package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/commands"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"
)

var _ = Describe("Packages", func() {
	Describe("Execute", func() {
		It("lists packages from package path", func() {
			console := ui.NewMemoryConsole()
			packagePath := "/opt/wrangle/packages"

			fileSystem := afero.NewMemMapFs()
			afero.WriteFile(fileSystem, "/opt/wrangle/packages/test/0.1.1/test.0.1.1.yml", []byte("this is a package"), 0600)

			command := commands.NewPackages(fileSystem, console, packagePath)
			Expect(command).ToNot(BeNil())
			Expect(command.Execute()).To(BeNil())

			Expect(console.OutAsString()).To(Equal("name\tversion\n"))
		})
	})
})
