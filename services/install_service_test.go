package services_test

import (	
	"github.com/patrickhuber/wrangle/settings"
	"github.com/patrickhuber/wrangle/feed"
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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type InstallServiceTester interface{
	ExecuteInstallsPackage(platform, downloadFileName string)
	NewInstallServiceCreatesInstance(platform string)
}

type installServiceTester struct{		
	server           *httptest.Server
	fileSystem  filesystem.FileSystem	
}

func NewInstallServiceTester(fs filesystem.FileSystem) InstallServiceTester{
	tester := &installServiceTester {
		fileSystem : fs,		
	}
	return tester
}

func (t *installServiceTester) NewInstallServiceCreatesInstance(platform string){
	paths := t.createPaths(platform)
	manager := t.createManager(paths)
	service, err := services.NewInstallService(platform, t.fileSystem, manager)
	Expect(err).To(BeNil())
	Expect(service).ToNot(BeNil())
}

func (t *installServiceTester) createPaths(platform string) *settings.Paths{
	const wrangleRootPosix = "/opt/wrangle"
	const wrangleRootWindows = "c:/wrangle"

	wrangleRoot := wrangleRootPosix
	if platform == "windows" {
		wrangleRoot = wrangleRootWindows
	}

	paths := &settings.Paths{
		Root : wrangleRoot ,
		Bin: wrangleRoot + "/bin",
		Packages: wrangleRoot + "/packages",
	}
	return paths
}


func (t *installServiceTester) createManager(paths *settings.Paths) packages.Manager{
	console := ui.NewMemoryConsole()

	taskProviders := tasks.NewProviderRegistry()
	taskProviders.Register(tasks.NewExtractProvider(t.fileSystem, console))
	taskProviders.Register(tasks.NewDownloadProvider(t.fileSystem, console))		
	taskProviders.Register(tasks.NewMoveProvider(t.fileSystem, console))
	taskProviders.Register(tasks.NewLinkProvider(t.fileSystem, console))

	feedService := feed.NewFsFeedService(t.fileSystem, paths.Packages)
	contextProvider := packages.NewFsContextProvider(t.fileSystem, paths)
	manager := packages.NewManager(t.fileSystem, feedService, contextProvider, taskProviders)
	return manager
}

func (t *installServiceTester) ExecuteInstallsPackage(platform, downloadFileName string){	
	paths := t.createPaths(platform)
	manager := t.createManager(paths)

	// create the command and execute it
	service, err := services.NewInstallService(platform, t.fileSystem, manager)
	Expect(err).To(BeNil())

	server := fakes.NewHTTPServerWithArchive(
		[]fakes.TestFile{
			{Path: "/test", Data: "this is data"},
			{Path: "/test.exe", Data: "this is data"},
		})
	defer server.Close()

	url := server.URL
	packageVersion := "1.0.0"
    packageName := "test"

	out := filepath.Join("/", downloadFileName)
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	url += downloadFileName

	// create the package manifest
	packageManifest, err := t.createPackageManifest(packageName, packageVersion, platform, url, out, "", "")
	Expect(err).To(BeNil())

	packagePath := filepath.Join(paths.Packages, packageName, packageVersion)
	packageManifestPath := filepath.Join(packagePath, fmt.Sprintf("%s.%s.yml", packageName, packageVersion))
	err = t.fileSystem.Write(packageManifestPath, []byte(packageManifest), 0600)
	Expect(err).To(BeNil())

	err = service.Install(
		&services.InstallServiceRequest{
			Directories: &services.InstallServiceRequestDirectories{
				Bin:      paths.Bin,
				Root:     paths.Root,
				Packages: paths.Packages},
			Package: &services.InstallServiceRequestPackage{
				Name:    packageName,
				Version: packageVersion},
			Feed: &services.InstallServiceRequestFeed{},
		})
	Expect(err).To(BeNil())

	// verify the package is installed?
}

func (t *installServiceTester) createPackageManifest(
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

var _ = Describe("InstallService", func() {
	var (
		fs      filesystem.FileSystem
		manager packages.Manager
		paths *settings.Paths		
	)
	BeforeEach(func() {
		// create command dependencies
		console := ui.NewMemoryConsole()
		fs = filesystem.NewMemory()
		
		taskProviders := tasks.NewProviderRegistry()
		taskProviders.Register(tasks.NewExtractProvider(fs, console))
		taskProviders.Register(tasks.NewDownloadProvider(fs, console))		
		taskProviders.Register(tasks.NewMoveProvider(fs, console))
		taskProviders.Register(tasks.NewLinkProvider(fs, console))
		
		// needs to be cross platform
		paths = &settings.Paths{
			Root : "/opt/wrangle",
			Bin: "/opt/wrangle/bin",
			Packages: "/opt/wrangle/packages",
		}
		feedService := feed.NewFsFeedService(fs, paths.Packages)
		contextProvider := packages.NewFsContextProvider(fs, paths)
		manager = packages.NewManager(fs, feedService, contextProvider, taskProviders)
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
		var(
			tester InstallServiceTester
			platform string
		)
		BeforeEach(func(){
			tester = NewInstallServiceTester(fs)
		})		
		When("Windows", func() {
			BeforeEach(func() {				
				platform = "windows"
			})
			When("Tar", func() {
				It("installs", func() {					
					tester.ExecuteInstallsPackage(platform, "test.tar")
				})
			})
			When("Tgz", func() {
				It("installs", func() {
					tester.ExecuteInstallsPackage(platform, "test.tgz")					
				})
			})
			When("Zip", func() {
				It("installs", func() {					
					tester.ExecuteInstallsPackage(platform, "test.zip")
				})
			})
			When("Binary", func() {
				It("installs", func() {
					tester.ExecuteInstallsPackage(platform, "test.exe")
				})
			})
		})
		When("Linux", func() {
			BeforeEach(func() {
				platform = "linux"
			})
			When("Tar", func() {
				It("installs", func() {					
					tester.ExecuteInstallsPackage(platform, "test.tar")
				})
			})
			When("Tgz", func() {
				It("installs", func() {
					tester.ExecuteInstallsPackage(platform, "test.tgz")					
				})
			})
			When("Zip", func() {
				It("installs", func() {					
					tester.ExecuteInstallsPackage(platform, "test.zip")
				})
			})
			When("Binary", func() {
				It("installs", func() {
					tester.ExecuteInstallsPackage(platform, "test")
				})
			})
		})
		When("Darwin", func() {
			BeforeEach(func() {
				platform = "darwin"
			})
			When("Tar", func() {
				It("installs", func() {					
					tester.ExecuteInstallsPackage(platform, "test.tar")
				})
			})
			When("Tgz", func() {
				It("installs", func() {
					tester.ExecuteInstallsPackage(platform, "test.tgz")					
				})
			})
			When("Zip", func() {
				It("installs", func() {					
					tester.ExecuteInstallsPackage(platform, "test.zip")
				})
			})
			When("Binary", func() {
				It("installs", func() {
					tester.ExecuteInstallsPackage(platform, "test")
				})
			})
		})
	})
})