package commands_test

import (
	"github.com/patrickhuber/wrangle/templates"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/spf13/afero"
	"github.com/urfave/cli"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/commands"	
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/packages"
	"github.com/patrickhuber/wrangle/services"
	"github.com/patrickhuber/wrangle/tasks"
	"github.com/patrickhuber/wrangle/ui"
)

var _ = Describe("Install", func() {
	It("can run with environment variables", func() {
		// rewrite this test to use new package management features
		fs := filesystem.NewMemMapFs()
		console := ui.NewMemoryConsole()
		variables := collections.NewDictionary()		

		taskProviders := tasks.NewProviderRegistry()
		taskProviders.Register(tasks.NewDownloadProvider(fs, console))
		taskProviders.Register(tasks.NewExtractProvider(fs, console))
		taskProviders.Register(tasks.NewLinkProvider(fs, console))
		taskProviders.Register(tasks.NewMoveProvider(fs, console))

		templateFactory := templates.NewFactory(templates.NewMacroManagerFactory().Create())
		packagesManager := packages.NewManager(fs, taskProviders, templateFactory)

		installService, err := services.NewInstallService("linux", fs, packagesManager)
		Expect(err).To(BeNil())

		variables.Set(global.PackagePathKey, "/packages")
		os.Setenv(global.PackagePathKey, "/packages")

		// setup the test server
		message := "this is a message"

		// start the local http server
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(message))
		}))

		defer server.Close()

		content := `
name: test
version: 1.0.0
targets:
- platform: linux
  tasks:
  - download:      
      url: %s
      out: test.html
`
		content = fmt.Sprintf(content, server.URL)

		err = fs.Mkdir("/packages/test/1.0.0", 0666)
		Expect(err).To(BeNil())

		err = afero.WriteFile(fs, "/packages/test/1.0.0/test.1.0.0.yml", []byte(content), 0666)
		Expect(err).To(BeNil())

		app := cli.NewApp()
		app.Flags = []cli.Flag{
			cli.StringFlag{
				Name:   "config, c",
				Usage:  "Load configuration from `FILE`",
				EnvVar: global.ConfigFileKey,
				Value:  "/config",
			},
		}
		app.Commands = []cli.Command{
			*commands.CreateInstallCommand(installService),
		}

		err = app.Run([]string{
			"wrangle",
			"install",
			"test",
			"-v", "1.0.0",
			"-r", "/wrangle",
		})
		Expect(err).To(BeNil())

		ok, err := afero.Exists(fs, "/packages/test/1.0.0/test.html")
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

})
