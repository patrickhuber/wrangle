package commands_test

import (
	"fmt"
	"os"

	"github.com/patrickhuber/wrangle/filesystem"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/commands"
	"github.com/patrickhuber/wrangle/feed"
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/urfave/cli"
)

var _ = Describe("Packages", func() {
	It("lists local packages when package path is set", func() {
		console := ui.NewMemoryConsole()

		fs := filesystem.NewMemory()
		fs.Mkdir("/packages", 0600)
		var feedServiceFactory = feed.NewServiceFactory(fs)

		app := cli.NewApp()
		app.Name = "wrangle"
		app.Commands = []cli.Command{
			commands.CreateListPackagesCommand(console, feedServiceFactory),
		}

		os.Setenv(global.PackagePathKey, "/packages")
		err := app.Run([]string{"wrangle", "packages"})
		os.Unsetenv(global.PackagePathKey)

		fmt.Println(console.OutAsString())
		fmt.Println(console.ErrorAsString())

		Expect(err).To(BeNil())
	})
	It("lists remote packages when package path is not set", func() {
		console := ui.NewMemoryConsole()

		fs := filesystem.NewMemory()
		fs.Mkdir("/packages", 0600)
		var feedServiceFactory = feed.NewServiceFactory(fs)

		app := cli.NewApp()
		app.Name = "wrangle"
		app.Commands = []cli.Command{
			commands.CreateListPackagesCommand(console, feedServiceFactory),
		}

		os.Unsetenv(global.PackagePathKey)
		err := app.Run([]string{"wrangle", "packages"})

		fmt.Println(console.OutAsString())
		fmt.Println(console.ErrorAsString())

		Expect(err).To(BeNil())
	})
})
