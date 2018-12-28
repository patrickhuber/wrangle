package packages_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/patrickhuber/wrangle/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/fakes"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/packages"
	"github.com/patrickhuber/wrangle/tasks"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"
)

var _ = Describe("Manager", func() {
	var (
		manager packages.Manager
		fs      filesystem.FsWrapper
		context packages.PackageContext
	)
	BeforeEach(func() {
		fs = filesystem.NewMemMapFs()
		console := ui.NewMemoryConsole()

		registry := tasks.NewProviderRegistry()
		registry.Register(tasks.NewDownloadProvider(fs, console))
		registry.Register(tasks.NewExtractProvider(fs, console))
		registry.Register(tasks.NewLinkProvider(fs, console))

		manager = packages.NewManager(fs, registry)

		context = packages.NewDefaultContext("/wrangle", "test", "1.0.0")
	})
	Describe("Download", func() {
		Context("WhenDownloadSucceeds", func() {
			var (
				server   *httptest.Server
				fileName string
			)
			AfterEach(func() {
				server = fakes.NewHTTPServer()
				defer server.Close()

				url := server.URL
				if !strings.HasSuffix(url, "/") {
					url += "/"
				}
				url += fileName

				download := tasks.NewDownloadTask(url, fileName)
				p := packages.New("", "", context, download)

				err := manager.Install(p)
				Expect(err).To(BeNil())

				expected := filepath.Join(context.PackageVersionPath(), fileName)
				ok, err := afero.Exists(fs, expected)
				Expect(err).To(BeNil())
				Expect(ok).To(BeTrue())
			})
			It("can download binary package", func() {
				fileName = "data"
			})
			It("can download tgz package", func() {
				fileName = "test.tgz"
			})
			It("can download zip package", func() {
				fileName = "test.zip"
			})
			It("can download tar package", func() {
				fileName = "test.tar"
			})
			It("can download tar.gz package", func() {
				fileName = "test.tar.gz"
			})
		})
		Context("WhenDownloadFails", func() {
			var (
				server *httptest.Server
			)
			BeforeEach(func() {
				// start the local http server
				server = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					rw.WriteHeader(404)
					rw.Write([]byte("failure"))
				}))
				defer server.Close()
			})
			It("does not write a file", func() {
				pkg := packages.New(
					"", "", context,
					tasks.NewDownloadTask(server.URL, "/test/file"))

				err := manager.Install(pkg)
				Expect(err).ToNot(BeNil())

				ok, err := afero.Exists(fs, "/test/file")
				Expect(err).To(BeNil())
				Expect(ok).To(BeFalse())
			})
			It("returns an error", func() {
				pkg := packages.New(
					"", "", context,
					tasks.NewDownloadTask(server.URL, ""))

				err := manager.Install(pkg)
				Expect(err).ToNot(BeNil())
			})
		})
	})
	Describe("Extract", func() {
		var (
			server   *httptest.Server
			fileName string
		)
		Context("WhenFlat", func() {
			AfterEach(func() {
				server = fakes.NewHTTPServer()
				defer server.Close()

				url := server.URL
				if !strings.HasSuffix(url, "/") {
					url += "/"
				}
				url += fileName

				archive := fileName
				download := tasks.NewDownloadTask(url, archive)
				extract := tasks.NewExtractTask(archive)
				p := packages.New("", "", context, download, extract)

				err := manager.Install(p)
				Expect(err).To(BeNil())

				expected := filepath.Join(context.PackageVersionPath(), "data")

				ok, err := afero.Exists(fs, expected)
				Expect(err).To(BeNil())
				Expect(ok).To(BeTrue())
			})
			It("can extract tgz", func() {
				fileName = "test.tgz"
			})
			It("can extract tar.gz", func() {
				fileName = "test.tar.gz"
			})
			It("can extract tar", func() {
				fileName = "test.tar"
			})
			It("can extract zip", func() {
				fileName = "test.zip"
			})
		})
		Context("WhenNested", func() {
			It("can extract", func() {
				var files = []fakes.TestFile{
					{Path: "/test1", Data: "not right"},
					{Path: "/parent/child", Data: "test\n"},
					{Path: "/parent/test2/grand-child", Data: "also not right"},
				}
				server := fakes.NewHTTPServerWithArchive(files)
				defer server.Close()

			})
		})
	})
	Describe("Symlink", func() {
		var (
			server *httptest.Server
		)
		BeforeEach(func() {
			// start the local http server
			server = fakes.NewHTTPServer()
			defer server.Close()

			url := server.URL
			if !strings.HasSuffix(url, "/") {
				url += "/"
			}
			url += "data"

			pkg := packages.New("test", "1.0.0",
				context,
				tasks.NewDownloadTask(url, "data"),
				tasks.NewLinkTask("data", "symlink"))

			err := manager.Install(pkg)
			Expect(err).To(BeNil())
		})
		Context("WhenSymlinkExists", func() {
			It("deletes existing symlink", func() {
				afero.WriteFile(fs, "/out/symlink", []byte(""), 0666)
			})
		})
	})

	Describe("Load", func() {
		BeforeEach(func() {
			filePathTemplate := "/wrangle/packages/test/{{version}}/test.{{version}}.yml"
			contentTemplate := `name: test
version: {{version}}
targets:
- platform: windows
  tasks:
  - download:
      url: https://www.google.com
      out: index.((version)).html
`
			versions := []string{"1.0.0", "1.0.1", "1.1.0", "2.0.0"}
			for _, version := range versions {
				path := strings.Replace(filePathTemplate, "{{version}}", version, -1)
				content := strings.Replace(contentTemplate, "{{version}}", version, -1)
				afero.WriteFile(fs, path, []byte(content), 0666)
			}
		})
		Context("WhenVersionSpecified", func() {
			It("loads specified version", func() {
				pkg, err := manager.Load("/wrangle", "/wrangle/bin", "/wrangle/packages", "test", "1.0.0")
				Expect(err).To(BeNil())
				Expect(pkg).ToNot(BeNil())
				Expect(pkg.Name()).To(Equal("test"))
				Expect(pkg.Version()).To(Equal("1.0.0"))
				Expect(len(pkg.Tasks())).To(Equal(1))
			})
		})
		Context("WhenVersionNotSpecified", func() {
			It("loads latest version by semver", func() {
				pkg, err := manager.Load("/wrangle", "/wrangle/bin", "/wrangle/packages", "test", "")
				Expect(err).To(BeNil())
				Expect(pkg).ToNot(BeNil())
				Expect(pkg.Name()).To(Equal("test"))
				Expect(pkg.Version()).To(Equal("2.0.0"))
				Expect(len(pkg.Tasks())).To(Equal(1))
			})
			/* When("latest file present", func(){
				It("loads specified version in file", func(){
					afero.WriteFile(fs, "/packages/test/latest", []byte("1.1.0"), 0666)
					pkg, err := manager.Load("/packages", "test", "")
					Expect(err).To(BeNil())
					Expect(pkg).ToNot(BeNil())
					Expect(pkg.Name()).To(Equal("test"))
					Expect(pkg.Version()).To(Equal("1.1.0"))
					Expect(len(pkg.Tasks())).To(Equal(0))
				})
			}) */
		})
		It("should interpolate version", func() {
			pkg, err := manager.Load("/wrangle", "/wrangle/bin", "/wrangle/packages", "test", "1.1.0")
			Expect(err).To(BeNil())
			Expect(pkg).ToNot(BeNil())
			Expect(len(pkg.Tasks())).To(Equal(1))
			task := pkg.Tasks()[0]
			out, ok := task.Params().Lookup("out")
			Expect(ok).To(BeTrue())
			Expect(out).To(Equal("index.1.1.0.html"))
		})
	})
})
