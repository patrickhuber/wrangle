package commands_test

import (
	"net/http"
	"os"
	"fmt"
	"net/http/httptest"
	
	"github.com/urfave/cli"
	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	
	"github.com/patrickhuber/wrangle/services"
	"github.com/patrickhuber/wrangle/commands"
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/tasks"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/packages"
)

var _ = Describe("Install", func() {
	It("can run with environment variables", func() {
		// rewrite this test to use new package management features
		fs := filesystem.NewMemMapFs()
		console := ui.NewMemoryConsole()		
		variables := collections.NewDictionary()
		loader := config.NewLoader(fs)

		taskProviders := tasks.NewProviderRegistry()
		taskProviders.Register(tasks.NewDownloadProvider(fs, console))
		taskProviders.Register(tasks.NewExtractProvider(fs, console))
		taskProviders.Register(tasks.NewLinkProvider(fs, console))
		taskProviders.Register(tasks.NewMoveProvider(fs, console))

		packagesManager := packages.NewManager(fs, taskProviders)

		installService, err := services.NewInstallService("linux", fs, packagesManager, loader)
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
package:
  name: test
  version: 1.0.0
  platforms:
  - name: linux
    tasks:
    - name: download
      type: download
      params: 
        url: %s
        out: ((package_install_directory))/test.html
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
		})
		Expect(err).To(BeNil())

		err = listFiles(fs, "/")
		Expect(err).To(BeNil())		

		ok, err := afero.Exists(fs, "/packages/test/1.0.0/test.html")
		Expect(err).To(BeNil())			
		Expect(ok).To(BeTrue())
	})

	
})

func listFiles(fs afero.Fs, directory string) error{
	files, err := afero.ReadDir(fs, directory)
	if err != nil{
		return err
	}
	for _, file := range files{
		os.Stdout.WriteString(
			fmt.Sprintf("%s/%s", directory, file.Name()))
	}
	return nil
}