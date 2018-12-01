package commands_test

import (
	"github.com/urfave/cli"
	"github.com/patrickhuber/wrangle/services"
	"github.com/patrickhuber/wrangle/commands"
	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Init", func() {
	It("creates config file", func() {
		fileSystem := afero.NewMemMapFs()
		initService := services.NewInitService(fileSystem)
		initCommand := commands.CreateInitCommand(initService)

		context := &cli.Context{}
		initCommand.Action.(func (context *cli.Context)error)(context)
		ok, err := afero.Exists(fileSystem, "/test")
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

		data, err := afero.ReadFile(fileSystem, "/test")
		Expect(err).To(BeNil())
		Expect(string(data)).To(Equal("stores: \nprocesses: \n"))
	})
})
