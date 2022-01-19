package tasks_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/archive"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/tasks"
	"github.com/spf13/afero"
)

type TestFile struct {
	Name    string
	Content string
}

var _ = Describe("Extract", func() {

	var (
		archiveName string
		fs          filesystem.FileSystem
		files       []*TestFile
	)
	BeforeEach(func() {
		fs = filesystem.FromAferoFS(afero.NewMemMapFs())
		files = []*TestFile{
			{
				Name:    "1.txt",
				Content: "test",
			},
		}
		for _, f := range files {
			Expect(fs.Write(f.Name, []byte(f.Content), 0644)).To(BeNil())
		}
	})
	It("can extract zip", func() {
		archiveName = "archive.zip"
	})
	It("can extract tgz", func() {
		archiveName = "archive.tgz"
	})
	It("can extract tar.gz", func() {
		archiveName = "archive.tar.gz"
	})
	It("can extract tar", func() {
		archiveName = "archive.tar"
	})
	AfterEach(func() {
		names := []string{}
		for _, f := range files {
			Expect(fs.Write(f.Name, []byte(f.Content), 0644)).To(BeNil())
			names = append(names, f.Name)
		}

		factory := archive.NewFactory(fs)
		provider, err := factory.Select(archiveName)
		Expect(err).To(BeNil())

		provider.Archive(archiveName, names...)

		for _, f := range files {
			Expect(fs.Remove(f.Name)).To(BeNil())
		}

		extract := tasks.NewExtractProvider(factory)
		Expect(provider).ToNot(BeNil())

		task := &tasks.Task{
			Type: "extract",
			Parameters: map[string]interface{}{
				"archive": archiveName,
				"out":     files[0].Name,
			},
		}
		metadata := &tasks.Metadata{}
		err = extract.Execute(task, metadata)
		Expect(err).To(BeNil())

		for _, f := range files {
			filePath := "/" + f.Name
			ok, err := fs.Exists(filePath)
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue(), "file %s does not exist", filePath)
			bytes, err := fs.Read(filePath)
			Expect(err).To(BeNil())
			Expect(string(bytes)).To(Equal(f.Content))
		}
	})

})
