package tasks_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/archive"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/ilog"
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
		task        *tasks.Task
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
		task = &tasks.Task{
			Type: "extract",
			Parameters: map[string]interface{}{
				"archive": archiveName,
				"out":     files[0].Name,
			},
		}
	})
	It("can extract tgz", func() {
		archiveName = "archive.tgz"
		task = &tasks.Task{
			Type: "extract",
			Parameters: map[string]interface{}{
				"archive": archiveName,
				"out":     files[0].Name,
			},
		}
	})
	It("can extract tar.gz", func() {
		archiveName = "archive.tar.gz"
		task = &tasks.Task{
			Type: "extract",
			Parameters: map[string]interface{}{
				"archive": archiveName,
				"out":     files[0].Name,
			},
		}
	})
	It("can extract tar", func() {
		archiveName = "archive.tar"
		task = &tasks.Task{
			Type: "extract",
			Parameters: map[string]interface{}{
				"archive": archiveName,
				"out":     files[0].Name,
			},
		}
	})
	When("no out specified", func() {
		It("can extract tar", func() {
			archiveName = "archive.tar"
			task = &tasks.Task{
				Type: "extract",
				Parameters: map[string]interface{}{
					"archive": archiveName,
				},
			}
		})
	})
	AfterEach(func() {
		// file names and archive names are not rooted
		// create the rooted versions
		packageVersionPath := "/"
		rootedFiles := []string{}
		for _, f := range files {
			filePath := crosspath.Join(packageVersionPath, f.Name)
			Expect(fs.Write(filePath, []byte(f.Content), 0644)).To(BeNil())
			rootedFiles = append(rootedFiles, filePath)
		}

		// setup
		logger := ilog.Memory()
		factory := archive.NewFactory(fs)
		provider, err := factory.Select(archiveName)
		Expect(err).To(BeNil())

		// create the test archive
		archivePath := crosspath.Join(packageVersionPath, archiveName)
		Expect(provider.Archive(archivePath, rootedFiles...)).To(BeNil())

		// cleanup so when we roundtrip we see the actual files
		for _, f := range rootedFiles {
			Expect(fs.Remove(f)).To(BeNil())
		}

		extract := tasks.NewExtractProvider(factory, logger)
		Expect(provider).ToNot(BeNil())

		metadata := &tasks.Metadata{}
		err = extract.Execute(task, metadata)
		Expect(err).To(BeNil(), errorStringOrDefault(err))

		for _, f := range files {
			filePath := crosspath.Join(packageVersionPath, f.Name)
			ok, err := fs.Exists(filePath)
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue(), "file %s does not exist", filePath)
			bytes, err := fs.Read(filePath)
			Expect(err).To(BeNil())
			Expect(string(bytes)).To(Equal(f.Content))
		}
	})

})

func errorStringOrDefault(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
