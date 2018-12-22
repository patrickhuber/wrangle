package services_test

import (
	"github.com/patrickhuber/wrangle/services"
	"fmt"
	"net/http/httptest"
	"strings"

	"github.com/patrickhuber/wrangle/fakes"
	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/tasks"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/packages"

	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"	
)

var _ = Describe("InstallService", func() {
	var (
		fs      filesystem.FsWrapper
		manager packages.Manager
		loader  config.Loader
	)
	BeforeEach(func() {
		// create command dependencies
		console := ui.NewMemoryConsole()
		fs = filesystem.NewMemMapFs()
		loader = config.NewLoader(fs)

		taskProviders := tasks.NewProviderRegistry()
		taskProviders.Register(tasks.NewExtractProvider(fs, console))
		taskProviders.Register(tasks.NewDownloadProvider(fs, console))
		taskProviders.Register(tasks.NewMoveProvider(fs, console))
		taskProviders.Register(tasks.NewLinkProvider(fs, console))

		manager = packages.NewManager(fs, taskProviders)
	})
	Describe("NewInstall", func() {
		It("returns install command", func() {
			platform := "platform"
			service, err := services.NewInstallService(platform, fs, manager, loader)
			Expect(err).To(BeNil())
			Expect(service).ToNot(BeNil())
		})
	})
	Describe("Execute", func() {
		const wrangleRootPosix = "/opt/wrangle"
		const wrangleRootWindows = "c:/wrangle"
		var (
			platform         string
			downloadFileName string
			archive          string
			destination      string
			server           *httptest.Server
		)
		BeforeSuite(func() {
			server = fakes.NewHTTPServerWithArchive(
				[]fakes.TestFile{
					{Path: "/test", Data: "this is data"},
					{Path: "/test.exe", Data: "this is data"},
				})
		})
		AfterSuite(func() {
			server.Close()
		})
		AfterEach(func() {
			url := server.URL
			packageVersion := "1.0.0"
			packageName := "test"
			wrangleRoot := wrangleRootPosix			
			if platform == "windows" {
				wrangleRoot = wrangleRootWindows
			}
			packagesRoot := wrangleRoot + "/packages"
			packagesBin := wrangleRoot + "/bin"

			out := filepath.Join("/", downloadFileName)
			if !strings.HasSuffix(url, "/") {
				url += "/"
			}
			url += downloadFileName

			// create the package manifest
			packageManifest, err := createPackageManifest(packageName, packageVersion, platform, url, out, archive, destination)
			Expect(err).To(BeNil())

			packagePath := filepath.Join(packagesRoot, packageName, packageVersion)
			packageManifestPath := filepath.Join(packagePath, fmt.Sprintf("%s.%s.yml", packageName, packageVersion))
			err = afero.WriteFile(fs, packageManifestPath, []byte(packageManifest), 0600)
			Expect(err).To(BeNil())

			// create the command and execute it
			service, err := services.NewInstallService(platform, fs, manager, loader)
			Expect(err).To(BeNil())

			err = service.Install(wrangleRoot, packagesBin, packagesRoot, packageName, packageVersion)
			Expect(err).To(BeNil())
		})
		When("Windows", func() {
			BeforeEach(func() {
				platform = "windows"
			})
			When("Tar", func() {
				It("installs", func() {
					downloadFileName = "test.tar"
				})
			})
			When("Tgz", func() {
				It("installs", func() {
					downloadFileName = "test.tgz"
				})
			})
			When("Zip", func() {
				It("installs", func() {
					downloadFileName = "test.zip"
				})
			})
			When("Binary", func() {
				It("installs", func() {
					downloadFileName = "test.exe"
				})
			})
		})
		When("Linux", func() {
			BeforeEach(func() {
				platform = "linux"
			})
		})
		When("Darwin", func() {
			BeforeEach(func() {
				platform = "darwin"
			})
		})
	})
})

func createPackageManifest(
	name string,
	version string,
	platform string,
	url string,
	outFile string,
	archive string,
	destination string) (string, error) {
	
		taskList:= []interface{}{
			tasks.NewDownloadTask(url, outFile),
		}
		if len(archive) > 0 {
			extract := tasks.NewExtractTask(archive, destination)
			taskList = append(taskList, extract)
		}
		pkg := &config.Package{
			Name: name,
			Version: version,
			Targets: []config.Target{
				config.Target{
					Platform: platform,
					Architecture: "amd64",
					Tasks: taskList,
				},
			},
		}
		return config.SerializePackage(pkg)

}
