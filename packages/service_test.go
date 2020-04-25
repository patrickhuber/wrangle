package packages_test

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"strings"

	"github.com/patrickhuber/wrangle/feed"
	"github.com/patrickhuber/wrangle/settings"

	"github.com/patrickhuber/wrangle/fakes"
	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/packages"
	"github.com/patrickhuber/wrangle/tasks"
	"github.com/patrickhuber/wrangle/ui"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type ServiceTester interface {
	InstallsPackage(platform, downloadFileName string)
}

type serviceTester struct {
	server  *httptest.Server
	service packages.Service
	fs      filesystem.FileSystem
}

func NewServiceTester(service packages.Service, fs filesystem.FileSystem) ServiceTester {
	tester := &serviceTester{
		service: service,
		fs:      fs,
	}
	return tester
}

func (t *serviceTester) createPaths(platform string) *settings.Paths {
	const wrangleRootPosix = "/opt/wrangle"
	const wrangleRootWindows = "c:/wrangle"

	wrangleRoot := wrangleRootPosix
	if platform == "windows" {
		wrangleRoot = wrangleRootWindows
	}

	paths := &settings.Paths{
		Root:     wrangleRoot,
		Bin:      wrangleRoot + "/bin",
		Packages: wrangleRoot + "/packages",
	}
	return paths
}

func (t *serviceTester) createTaskProviderRegistry() tasks.ProviderRegistry {
	console := ui.NewMemoryConsole()

	taskProviders := tasks.NewProviderRegistry()
	taskProviders.Register(tasks.NewExtractProvider(t.fs, console))
	taskProviders.Register(tasks.NewDownloadProvider(t.fs, console))
	taskProviders.Register(tasks.NewMoveProvider(t.fs, console))
	taskProviders.Register(tasks.NewLinkProvider(t.fs, console))
	return taskProviders
}

func (t *serviceTester) InstallsPackage(platform, downloadFileName string) {
	paths := t.createPaths(platform)
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
	packageManifest, err := t.createPackageManifest(packageName, packageVersion, url, out, "", "")
	Expect(err).To(BeNil())

	packagePath := filepath.Join(paths.Packages, packageName, packageVersion)
	packageManifestPath := filepath.Join(packagePath, fmt.Sprintf("%s.%s.yml", packageName, packageVersion))
	err = t.fs.Write(packageManifestPath, []byte(packageManifest), 0600)
	Expect(err).To(BeNil())

	err = t.service.Install(
		&packages.InstallRequest{
			Directories: &packages.InstallRequestDirectories{
				Bin:      paths.Bin,
				Root:     paths.Root,
				Packages: paths.Packages},
			Package: &packages.InstallRequestPackage{
				Name:    packageName,
				Version: packageVersion},
			Feed:     &packages.InstallRequestFeed{},
			Platform: platform,
		})
	Expect(err).To(BeNil())

	// verify the package is installed?
}

func (t *serviceTester) createPackageManifest(
	name string,
	version string,
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
			{
				Platform:     "windows",
				Architecture: "amd64",
				Tasks:        taskList,
			},
			{
				Platform:     "linux",
				Architecture: "amd64",
				Tasks:        taskList,
			},
			{
				Platform:     "darwin",
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

var _ = Describe("Service", func() {
	var (
		fs      filesystem.FileSystem
		paths   *settings.Paths
		service packages.Service
	)
	BeforeEach(func() {
		// create command dependencies
		console := ui.NewMemoryConsole()
		fs = filesystem.NewMemory()
		filePathTemplate := "/opt/wrangle/packages/test/{{version}}/test.{{version}}.yml"
		contentTemplate := `name: test
version: {{version}}
targets:
- platform: windows
  tasks:
  - download:
      url: https://www.google.com
      out: index.{{version}}.html
- platform: linux
  tasks:
  - download:
      url: https://www.google.com
      out: index.{{version}}.html
- platform: unix
  tasks:
  - download:
      url: https://www.google.com
      out: index.{{version}}.html
`
		versions := []string{"1.0.0", "1.0.1", "1.1.0", "2.0.0"}
		for _, version := range versions {
			path := strings.Replace(filePathTemplate, "{{version}}", version, -1)
			content := strings.Replace(contentTemplate, "{{version}}", version, -1)
			fs.Write(path, []byte(content), 0666)
		}

		taskProviders := tasks.NewProviderRegistry()
		taskProviders.Register(tasks.NewExtractProvider(fs, console))
		taskProviders.Register(tasks.NewDownloadProvider(fs, console))
		taskProviders.Register(tasks.NewMoveProvider(fs, console))
		taskProviders.Register(tasks.NewLinkProvider(fs, console))

		// needs to be cross platform
		paths = &settings.Paths{
			Root:     "/opt/wrangle",
			Bin:      "/opt/wrangle/bin",
			Packages: "/opt/wrangle/packages",
		}
		feedService := feed.NewFsService(fs, paths.Packages)
		contextProvider := packages.NewFsContextProvider(fs, paths)
		interfaceReader := packages.NewYamlInterfaceReader()

		service = packages.NewService(feedService, interfaceReader, contextProvider, taskProviders)
	})
	Describe("Install", func() {
		var (
			tester   ServiceTester
			platform string
		)
		BeforeEach(func() {
			tester = NewServiceTester(service, fs)
		})
		When("Windows", func() {
			BeforeEach(func() {
				platform = "windows"
			})
			When("Tar", func() {
				It("installs", func() {
					tester.InstallsPackage(platform, "test.tar")
				})
			})
			When("Tgz", func() {
				It("installs", func() {
					tester.InstallsPackage(platform, "test.tgz")
				})
			})
			When("Zip", func() {
				It("installs", func() {
					tester.InstallsPackage(platform, "test.zip")
				})
			})
			When("Binary", func() {
				It("installs", func() {
					tester.InstallsPackage(platform, "test.exe")
				})
			})
		})
		When("Linux", func() {
			BeforeEach(func() {
				platform = "linux"
			})
			When("Tar", func() {
				It("installs", func() {
					tester.InstallsPackage(platform, "test.tar")
				})
			})
			When("Tgz", func() {
				It("installs", func() {
					tester.InstallsPackage(platform, "test.tgz")
				})
			})
			When("Zip", func() {
				It("installs", func() {
					tester.InstallsPackage(platform, "test.zip")
				})
			})
			When("Binary", func() {
				It("installs", func() {
					tester.InstallsPackage(platform, "test")
				})
			})
		})
		When("Darwin", func() {
			BeforeEach(func() {
				platform = "darwin"
			})
			When("Tar", func() {
				It("installs", func() {
					tester.InstallsPackage(platform, "test.tar")
				})
			})
			When("Tgz", func() {
				It("installs", func() {
					tester.InstallsPackage(platform, "test.tgz")
				})
			})
			When("Zip", func() {
				It("installs", func() {
					tester.InstallsPackage(platform, "test.zip")
				})
			})
			When("Binary", func() {
				It("installs", func() {
					tester.InstallsPackage(platform, "test")
				})
			})
		})
	})

	Describe("Get", func() {
		Context("WhenVersionSpecified", func() {
			It("loads specified version", func() {
				resp, err := service.Get(
					&packages.GetRequest{
						PackageName:    "test",
						PackageVersion: "1.0.0",
						Platform:       "windows",
					})
				Expect(err).To(BeNil())
				pkg := resp.Package
				Expect(pkg).ToNot(BeNil())
				Expect(pkg.Name()).To(Equal("test"))
				Expect(pkg.Version()).To(Equal("1.0.0"))
				Expect(len(pkg.Tasks())).To(Equal(1))
			})
		})
		Context("WhenVersionNotSpecified", func() {
			It("loads latest version by semver", func() {
				resp, err := service.Get(
					&packages.GetRequest{
						PackageName:    "test",
						PackageVersion: "",
						Platform:       "windows"})
				Expect(err).To(BeNil())
				pkg := resp.Package
				Expect(pkg).ToNot(BeNil())
				Expect(pkg.Name()).To(Equal("test"))
				Expect(pkg.Version()).To(Equal("2.0.0"))
				Expect(len(pkg.Tasks())).To(Equal(1))
			})
			When("latest tag", func() {
				It("loads version from tag", func() {
					fs.Write("/opt/wrangle/packages/test/latest", []byte("1.1.0"), 0666)
					resp, err := service.Get(
						&packages.GetRequest{
							PackageName:    "test",
							PackageVersion: "",
							Platform:       "windows"})
					Expect(err).To(BeNil())
					pkg := resp.Package
					Expect(pkg).ToNot(BeNil())
					Expect(pkg.Name()).To(Equal("test"))
					Expect(pkg.Version()).To(Equal("1.1.0"))
					Expect(len(pkg.Tasks())).To(Equal(1))
				})
			})
		})
		It("should interpolate version", func() {
			resp, err := service.Get(
				&packages.GetRequest{
					PackageName:    "test",
					PackageVersion: "1.1.0",
					Platform:       "windows"})
			Expect(err).To(BeNil())
			pkg := resp.Package
			Expect(pkg).ToNot(BeNil())
			Expect(len(pkg.Tasks())).To(Equal(1))
			task := pkg.Tasks()[0]
			out, ok := task.Params()["out"]
			Expect(ok).To(BeTrue())
			Expect(out).To(Equal("index.1.1.0.html"))
		})
		It("should only return specified platform tasks", func() {
			resp, err := service.Get(
				&packages.GetRequest{
					PackageName:    "test",
					PackageVersion: "1.1.0",
					Platform:       "windows"})
			Expect(err).To(BeNil())
			pkg := resp.Package
			Expect(pkg).ToNot(BeNil())
			Expect(len(pkg.Tasks())).To(Equal(1))
		})
	})
})
