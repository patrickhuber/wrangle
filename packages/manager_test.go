package packages_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

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
	)
	BeforeEach(func() {
		fs = filesystem.NewMemMapFs()
		console := ui.NewMemoryConsole()

		registry := tasks.NewProviderRegistry()
		registry.Register(tasks.NewDownloadProvider(fs, console))
		registry.Register(tasks.NewExtractProvider(fs, console))
		registry.Register(tasks.NewLinkProvider(fs, console))

		manager = packages.NewManager(fs, registry)
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

				download := tasks.NewDownloadTask(url, fmt.Sprintf("/out/%v", fileName))
				p := packages.New("", "", download)

				err := manager.Install(p)
				Expect(err).To(BeNil())

				ok, err := afero.Exists(fs, "/out")
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
					"", "",
					tasks.NewDownloadTask(server.URL, "/test/file"))

				err := manager.Install(pkg)
				Expect(err).ToNot(BeNil())

				ok, err := afero.Exists(fs, "/test/file")
				Expect(err).To(BeNil())
				Expect(ok).To(BeFalse())
			})
			It("returns an error", func() {
				pkg := packages.New(
					"", "",
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

				archive := fmt.Sprintf("/download/%v", fileName)
				download := tasks.NewDownloadTask( url, archive)
				extract := tasks.NewExtractTask( archive, "/extract")
				p := packages.New("", "", download, extract)

				err := manager.Install(p)
				Expect(err).To(BeNil())

				ok, err := afero.Exists(fs, "/extract")
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

			pkg := packages.New("", "",
				tasks.NewDownloadTask(url, "/out/data"),
				tasks.NewLinkTask( "/out/data", "/out/symlink"))

			err := manager.Install(pkg)
			Expect(err).To(BeNil())
		})
		Context("WhenSymlinkExists", func() {
			It("deletes existing symlink", func() {
				afero.WriteFile(fs, "/out/symlink", []byte(""), 0666)
			})
		})
	})
})
