package commands_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/commands"
	"github.com/patrickhuber/wrangle/feed"
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/services"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"
	"github.com/urfave/cli"
)

var _ = Describe("Packages", func() {
	It("lists local packages when package path is set", func() {
		console := ui.NewMemoryConsole()
		var packageServiceFactory services.PackageServiceFactory = services.NewPackageServiceFactory(console)

		fs := afero.NewMemMapFs()
		fs.Mkdir("/packages", 0600)
		var feedServiceFactory = feed.NewFeedServiceFactory(fs)

		app := cli.NewApp()
		app.Name = "wrangle"
		app.Commands = []cli.Command{
			*commands.CreatePackagesCommand(packageServiceFactory, feedServiceFactory),
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
		var packageServiceFactory services.PackageServiceFactory = services.NewPackageServiceFactory(console)

		fs := afero.NewMemMapFs()
		fs.Mkdir("/packages", 0600)
		var feedServiceFactory = feed.NewFeedServiceFactory(fs)

		app := cli.NewApp()
		app.Name = "wrangle"
		app.Commands = []cli.Command{
			*commands.CreatePackagesCommand(packageServiceFactory, feedServiceFactory),
		}

		os.Unsetenv(global.PackagePathKey)
		err := app.Run([]string{"wrangle", "packages"})

		fmt.Println(console.OutAsString())
		fmt.Println(console.ErrorAsString())

		Expect(err).To(BeNil())
	})
})