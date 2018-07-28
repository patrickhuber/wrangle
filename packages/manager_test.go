package packages

import (
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/fakes"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/spf13/afero"
)

var _ = Describe("Manager", func() {
	Describe("Download", func() {
		It("can download binary package", func() {
			testDownloadFile("data")
		})
		It("can download tgz package", func() {
			testDownloadFile("test.tgz")
		})
		It("can download zip package", func() {
			testDownloadFile("test.zip")
		})
		It("can download tar package", func() {
			testDownloadFile("test.tar")
		})
		Context("WhenDownloadFails", func() {
			It("does not write a file", func() {
				fileSystem := filesystem.NewMemMapFs()

				// start the local http server
				server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					rw.WriteHeader(404)
					rw.Write([]byte("failure"))
				}))

				defer server.Close()

				pkg := New(
					"", "", "",
					NewDownload(server.URL, "/test", "file"),
					nil)

				manager := NewManager(fileSystem)

				err := manager.Download(pkg)
				Expect(err).ToNot(BeNil())

				ok, err := afero.Exists(fileSystem, "/test/file")
				Expect(err).To(BeNil())
				Expect(ok).To(BeFalse())
			})
			It("returns an error", func() {
				fileSystem := filesystem.NewMemMapFs()

				// start the local http server
				server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					rw.WriteHeader(404)
					rw.Write([]byte("failure"))
				}))

				defer server.Close()

				pkg := New(
					"", "", "",
					NewDownload(server.URL, "", ""),
					nil)

				manager := NewManager(fileSystem)

				err := manager.Download(pkg)
				Expect(err).ToNot(BeNil())
			})
		})
	})
	Describe("Extract", func() {
		It("can extract tgz", func() {
			testExtractFile("test.tgz")
		})
		It("can extract tar.gz", func() {
			testExtractFile("test.tar.gz")
		})
		It("can extract tar", func() {
			testExtractFile("test.tar")
		})
		It("can extract zip", func() {
			testExtractFile("test.zip")
		})
		It("can extract nested file from tar", func() {
			var files = []fakes.TestFile{
				{Path: "/test1", Data: "not right"},
				{Path: "/parent/child", Data: "test\n"},
				{Path: "/parent/test2/grand-child", Data: "also not right"},
			}

			testDownloadExtractAndLink(files, "out.tar", "/child", "child", "link")
		})
	})
	Describe("Symlink", func() {
		It("can create symlink for binary", func() {
			// start the local http server
			server := fakes.NewHTTPServerWithArchive([]fakes.TestFile{{Path: "/data", Data: "this is data"}})
			defer server.Close()

			url := server.URL
			if !strings.HasSuffix(url, "/") {
				url += "/"
			}
			url += "data"

			pkg := New("", "", "symlink",
				NewDownload(url, "/out", "data"),
				nil)

			fs := filesystem.NewMemMapFs()
			manager := NewManager(fs)

			err := manager.Download(pkg)
			Expect(err).To(BeNil())

			err = manager.Link(pkg)
			Expect(err).To(BeNil())
		})
		Context("WhenSymlinkExists", func() {
			It("deletes existing symlink", func() {
				fs := filesystem.NewMemMapFs()

				oldname := "/out/existing"
				newname := "/out/symlink"

				afero.WriteFile(fs, oldname, []byte(""), 0666)

				err := fs.Symlink(oldname, newname)
				Expect(err).To(BeNil())

				// start the local http server
				server := fakes.NewHTTPServerWithArchive([]fakes.TestFile{{Path: "/data", Data: "this is data"}})
				defer server.Close()

				url := server.URL
				if !strings.HasSuffix(url, "/") {
					url += "/"
				}
				url += "data"

				pkg := New("", "", "symlink",
					NewDownload(url, "/out", "data"),
					nil)

				manager := NewManager(fs)

				err = manager.Download(pkg)
				Expect(err).To(BeNil())

				err = manager.Link(pkg)
				Expect(err).To(BeNil())
			})
		})
	})
})

func testDownloadExtractAndLink(files []fakes.TestFile, downloadOut, extractFilter, extractOut, alias string) {
	server := fakes.NewHTTPServerWithArchive(files)
	defer server.Close()

	url := server.URL
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	url += downloadOut

	download := NewDownload(url, "/out", downloadOut)
	var extract Extract
	if extractFilter != "" && extractOut != "" {
		extract = NewExtract(extractFilter, "/out", extractOut)
	}

	p := New("", "", alias, download, extract)

	fs := filesystem.NewMemMapFsWrapper(afero.NewMemMapFs())
	manager := NewManager(fs)

	err := manager.Download(p)
	Expect(err).To(BeNil())

	if extract != nil {
		err = manager.Extract(p)
		Expect(err).To(BeNil())

		ok, err := afero.Exists(fs, extract.OutPath())
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	}

	if alias == "" {
		return
	}
	err = manager.Link(p)
	Expect(err).To(BeNil())
}

func testDownloadFile(fileName string) {
	server := fakes.NewHTTPServer()
	defer server.Close()

	url := server.URL
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	url += fileName

	download := NewDownload(url, "/out", fileName)
	p := New("", "", "", download, nil)

	fs := filesystem.NewMemMapFs()

	manager := NewManager(fs)
	err := manager.Download(p)
	Expect(err).To(BeNil())

	ok, err := afero.Exists(fs, download.OutPath())
	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue())
}

func testExtractFile(fileName string) {
	server := fakes.NewHTTPServer()
	defer server.Close()

	extract := NewExtract(".*", "/out", "test.txt")
	url := server.URL
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	url += fileName
	p := New("", "", "",
		NewDownload(url, "/in", fileName),
		extract)

	fs := filesystem.NewMemMapFs()

	manager := NewManager(fs)

	err := manager.Download(p)
	Expect(err).To(BeNil())

	err = manager.Extract(p)
	Expect(err).To(BeNil())

	ok, err := afero.Exists(fs, extract.OutPath())
	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue())
}
