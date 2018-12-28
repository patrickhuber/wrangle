package archiver_test

import (
	"fmt"

	"github.com/patrickhuber/wrangle/archiver"
	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tar", func() {
	var (
		a       archiver.Archiver
		fs      afero.Fs
		archive = "/tmp/temp.tar"
	)
	BeforeEach(func() {
		fs = afero.NewMemMapFs()

		err := createFiles(fs, []testFile{
			{content: "this is a test", folder: "/tmp", name: "test"},
			{content: "this is a test", folder: "/tmp", name: "test1"},
			{content: "this is a test", folder: "/tmp", name: "test2"},
		})
		Expect(err).To(BeNil())

		a = archiver.NewTarArchiver(fs)
		err = a.Archive(archive, []string{"/tmp/test", "/tmp/test1", "/tmp/test2"})
		Expect(err).To(BeNil())
	})

	Describe("RoundTrip", func() {
		It("Can write and read back a tar file", func() {
			testExtractTar(fs, a, archive, "/out", []string{".*"})
			assertIsFile(fs, "/out/test")
		})
	})

	Describe("Extract", func() {
		Context("WhenFilterIsSet", func() {
			It("extracts only matching files", func() {

				testExtractTar(fs, a, archive, "/out", []string{"^test$"})
				assertExists(fs, "/out/test")
				assertIsFile(fs, "/out/test")
				assertNotExists(fs, "/out/test1")
			})
		})
	})

	AfterEach(func() {
		Expect(fs.RemoveAll("/out")).To(BeNil())
	})
})

func assertExists(fs afero.Fs, filePath string) {
	ok, err := afero.Exists(fs, filePath)
	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue(), fmt.Sprintf("'%s' should exist", filePath))
}

func assertNotExists(fs afero.Fs, filePath string) {
	ok, err := afero.Exists(fs, filePath)
	Expect(err).To(BeNil())
	Expect(ok).To(BeFalse(), fmt.Sprintf("'%s' should not exist", filePath))
}

func assertIsFile(fs afero.Fs, filePath string) {
	ok, err := afero.IsDir(fs, filePath)
	Expect(err).To(BeNil())
	Expect(ok).To(BeFalse(), fmt.Sprintf("'%s' should be a file", filePath))
}

func testExtractTar(fs afero.Fs, a archiver.Archiver, filePath string, out string, files []string) {
	source, err := fs.Open(filePath)
	Expect(err).To(BeNil())
	defer source.Close()

	err = a.Extract(filePath, out, files)
	Expect(err).To(BeNil())
}
