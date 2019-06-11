package services_test

import (
	"github.com/patrickhuber/wrangle/templates"
	"bytes"
	"fmt"
	"net/http/httptest"
	"strings"

	"github.com/patrickhuber/wrangle/services"

	"github.com/patrickhuber/wrangle/fakes"
	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/packages"
	"github.com/patrickhuber/wrangle/tasks"
	"github.com/patrickhuber/wrangle/ui"

	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InstallService", func() {
	var (
		fs      filesystem.FsWrapper
		manager packages.Manager
	)
	BeforeEach(func() {
		// create command dependencies
		console := ui.NewMemoryConsole()
		fs = filesystem.NewMemMapFs()

		taskProviders := tasks.NewProviderRegistry()
		taskProviders.Register(tasks.NewExtractProvider(fs, console))
		taskProviders.Register(tasks.NewDownloadProvider(fs, console))
		taskProviders.Register(tasks.NewMoveProvider(fs, console))
		taskProviders.Register(tasks.NewLinkProvider(fs, console))

		templateFactory := templates.NewFactory(templates.NewMacroManagerFactory().Create())
		manager = packages.NewManager(fs, taskProviders, templateFactory)
	})
	Describe("NewInstall", func() {
		It("returns install command", func() {
			platform := "platform"
			service, err := services.NewInstallService(platform, fs, manager)
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
			service, err := services.NewInstallService(platform, fs, manager)
			Expect(err).To(BeNil())

			err = service.Install(
				&services.InstallServiceRequest{
					Directories: &services.InstallServiceRequestDirectories{
						Bin:      packagesBin,
						Root:     wrangleRoot,
						Packages: packagesRoot},
					Package: &services.InstallServiceRequestPackage{
						Name:    packageName,
						Version: packageVersion},
					Feed: &services.InstallServiceRequestFeed{},
				})
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

	taskList := []interface{}{
		tasks.NewDownloadTask(url, outFile),
	}
	if len(archive) > 0 {
		extract := tasks.NewExtractTask(archive)
		taskList = append(taskList, extract)
	}
	pkg := &packages.Manifest{
		Name:    name,
		Version: version,
		Targets: []packages.Target{
			packages.Target{
				Platform:     platform,
				Architecture: "amd64",
				Tasks:        taskList,
			},
		},
	}
	var buffer bytes.Buffer
	writer := packages.NewYamlManifestWriter(&buffer)
	err := writer.Write(pkg)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil

}
