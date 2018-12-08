package commands_test

import (
	"github.com/patrickhuber/wrangle/global"
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
		
		defaultConfigPath := "/config"

		app := cli.NewApp()
		app.Flags = []cli.Flag{
			cli.StringFlag{
				Name:   "config, c",
				Usage:  "Load configuration from `FILE`",
				EnvVar: global.ConfigFileKey,
				Value:  defaultConfigPath,
			},
		}
		app.Commands = []cli.Command{
			*commands.CreateInitCommand(initService),
		}
		
		err := app.Run([]string{"wrangle", "init"})
		Expect(err).To(BeNil())

		ok, err := afero.Exists(fileSystem, "/config")
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

		data, err := afero.ReadFile(fileSystem, "/config")
		Expect(err).To(BeNil())
		Expect(string(data)).To(Equal("stores: \nprocesses: \n"))
	})
})
