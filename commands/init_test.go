package commands_test

import (
	"github.com/patrickhuber/wrangle/commands"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/services"
	"github.com/urfave/cli"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Init", func() {
	It("creates config file", func() {
		fileSystem := filesystem.NewMemory()
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
			*commands.CreateInitCommand(app, initService),
		}

		err := app.Run([]string{"wrangle", "init"})
		Expect(err).To(BeNil())

		ok, err := fileSystem.Exists("/config")
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

		data, err := fileSystem.Read("/config")
		Expect(err).To(BeNil())
		Expect(string(data)).To(Equal("stores: \nprocesses: \n"))
	})
})
