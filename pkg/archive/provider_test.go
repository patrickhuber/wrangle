package archive_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-xplat/filepath"
	filesystem "github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/wrangle/pkg/archive"
)

type TestFile struct {
	Name    string
	Content string
}

var _ = Describe("Provider", func() {
	var (
		fs          filesystem.FS
		path        filepath.Processor
		files       []*TestFile
		provider    archive.Provider
		archiveFile string
	)
	BeforeEach(func() {
		fs = filesystem.NewMemory()
		files = []*TestFile{
			{
				Name:    "1.txt",
				Content: "1",
			},
			{
				Name:    "2.txt",
				Content: "2",
			},
		}

	})
	It("can roundtrip tar", func() {
		provider = archive.NewTar(fs, path)
		archiveFile = "test.tar"
	})
	It("can roundtrip zip", func() {
		provider = archive.NewZip(fs)
		archiveFile = "test.zip"
	})
	It("can roundtrip tgz", func() {
		provider = archive.NewTarGz(fs, path)
		archiveFile = "test.tgz"
	})
	AfterEach(func() {
		Expect(provider).ToNot(BeNil())
		names := []string{}
		for _, f := range files {
			err := fs.WriteFile(f.Name, []byte(f.Content), 0644)
			Expect(err).To(BeNil())
			names = append(names, f.Name)
		}

		Expect(provider.Archive(archiveFile, names...)).To(BeNil())
		ok, err := fs.Exists(archiveFile)

		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

		for _, f := range files {
			Expect(fs.Remove(f.Name)).To(BeNil())
		}

		destination := "/"
		Expect(provider.Extract(archiveFile, "/", names...)).To(BeNil())
		for _, f := range files {
			filePath := path.Join(destination, f.Name)
			ok, err := fs.Exists(filePath)
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue(), "%s does not exist", filePath)
		}
	})
})
