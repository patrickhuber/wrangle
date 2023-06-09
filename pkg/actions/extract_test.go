package actions_test

import (
	"path"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-xplat/filepath"
	filesystem "github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/wrangle/pkg/actions"
	"github.com/patrickhuber/wrangle/pkg/archive"
	"github.com/stretchr/testify/require"
)

type TestFile struct {
	Name    string
	Content string
}

type extractTest struct {
	archiveName string
	fs          filesystem.FS
	path        filepath.Processor
	files       []*TestFile
	action      *actions.Action
}

func SetupExtractTest(t *testing.T) *extractTest {
	fs := filesystem.NewMemory()
	files := []*TestFile{
		{
			Name:    "1.txt",
			Content: "test",
		},
	}
	for _, f := range files {
		require.Nil(t, fs.WriteFile(f.Name, []byte(f.Content), 0644))
	}
	et := &extractTest{
		fs:    fs,
		files: files,
	}
	return et
}

func TestCanExtractZip(t *testing.T) {
	test := SetupExtractTest(t)
	test.archiveName = "archive.zip"
	test.action = &actions.Action{
		Type: "extract",
		Parameters: map[string]interface{}{
			"archive": test.archiveName,
			"out":     test.files[0].Name,
		},
	}
	test.Execute(t)
}

func (et *extractTest) Execute(t *testing.T) {
	// file names and archive names are not rooted
	// create the rooted versions
	packageVersionPath := "/"
	rootedFiles := []string{}
	for _, f := range et.files {
		filePath := path.Join(packageVersionPath, f.Name)
		require.Nil(t, et.fs.WriteFile(filePath, []byte(f.Content), 0644))
		rootedFiles = append(rootedFiles, filePath)
	}

	// setup
	logger := log.Memory()
	factory := archive.NewFactory(et.fs, et.path)
	provider, err := factory.Select(et.archiveName)
	require.Nil(t, err)

	// create the test archive
	archivePath := path.Join(packageVersionPath, et.archiveName)
	err = provider.Archive(archivePath, rootedFiles...)
	require.Nil(t, err)

	// cleanup so when we roundtrip we see the actual files
	for _, f := range rootedFiles {
		err = et.fs.Remove(f)
		require.Nil(t, err)
	}

	extract := actions.NewExtractProvider(factory, logger)
	require.NotNil(t, provider)

	metadata := &actions.Metadata{}
	err = extract.Execute(et.action, metadata)
	Expect(err).To(BeNil(), errorStringOrDefault(err))

	for _, f := range et.files {
		filePath := path.Join(packageVersionPath, f.Name)
		ok, err := et.fs.Exists(filePath)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue(), "file %s does not exist", filePath)
		bytes, err := et.fs.ReadFile(filePath)
		Expect(err).To(BeNil())
		Expect(string(bytes)).To(Equal(f.Content))
	}
}

var _ = Describe("Extract", func() {

	var (
		archiveName string
		fs          filesystem.FS
		path        filepath.Processor
		files       []*TestFile
		task        *actions.Action
	)
	BeforeEach(func() {
		fs = filesystem.NewMemory()
		files = []*TestFile{
			{
				Name:    "1.txt",
				Content: "test",
			},
		}
		for _, f := range files {
			Expect(fs.WriteFile(f.Name, []byte(f.Content), 0644)).To(BeNil())
		}
	})
	It("can extract zip", func() {
		archiveName = "archive.zip"
		task = &actions.Action{
			Type: "extract",
			Parameters: map[string]interface{}{
				"archive": archiveName,
				"out":     files[0].Name,
			},
		}
	})
	It("can extract tgz", func() {
		archiveName = "archive.tgz"
		task = &actions.Action{
			Type: "extract",
			Parameters: map[string]interface{}{
				"archive": archiveName,
				"out":     files[0].Name,
			},
		}
	})
	It("can extract tar.gz", func() {
		archiveName = "archive.tar.gz"
		task = &actions.Action{
			Type: "extract",
			Parameters: map[string]interface{}{
				"archive": archiveName,
				"out":     files[0].Name,
			},
		}
	})
	It("can extract tar", func() {
		archiveName = "archive.tar"
		task = &actions.Action{
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
			task = &actions.Action{
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
			filePath := path.Join(packageVersionPath, f.Name)
			Expect(fs.WriteFile(filePath, []byte(f.Content), 0644)).To(BeNil())
			rootedFiles = append(rootedFiles, filePath)
		}

		// setup
		logger := log.Memory()
		factory := archive.NewFactory(fs, path)
		provider, err := factory.Select(archiveName)
		Expect(err).To(BeNil())

		// create the test archive
		archivePath := path.Join(packageVersionPath, archiveName)
		Expect(provider.Archive(archivePath, rootedFiles...)).To(BeNil())

		// cleanup so when we roundtrip we see the actual files
		for _, f := range rootedFiles {
			Expect(fs.Remove(f)).To(BeNil())
		}

		extract := actions.NewExtractProvider(factory, logger)
		Expect(provider).ToNot(BeNil())

		metadata := &actions.Metadata{}
		err = extract.Execute(task, metadata)
		Expect(err).To(BeNil(), errorStringOrDefault(err))

		for _, f := range files {
			filePath := path.Join(packageVersionPath, f.Name)
			ok, err := fs.Exists(filePath)
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue(), "file %s does not exist", filePath)
			bytes, err := fs.ReadFile(filePath)
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
