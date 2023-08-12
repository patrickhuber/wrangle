package actions_test

import (
	"testing"

	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-xplat/filepath"
	filesystem "github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/go-xplat/setup"
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
	path        *filepath.Processor
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

func TestCanExtract(t *testing.T) {

	files := []*TestFile{
		{
			Name:    "1.txt",
			Content: "test",
		},
	}
	// file names and archive names are not rooted
	// create the rooted versions
	metadata := &actions.Metadata{
		PackageVersionPath: "/some/child/path",
	}

	type test struct {
		archiveName string
		action      *actions.Action
	}

	tests := []test{
		{
			archiveName: "archive.zip",
			action: &actions.Action{
				Type: "extract",
				Parameters: map[string]any{
					"archive": "archive.zip",
					"out":     "1.txt",
				},
			},
		},
		{
			archiveName: "archive.tar",
			action: &actions.Action{
				Type: "extract",
				Parameters: map[string]any{
					"archive": "archive.tar",
					"out":     "1.txt",
				},
			},
		},
		{
			archiveName: "archive.tgz",
			action: &actions.Action{
				Type: "extract",
				Parameters: map[string]any{
					"archive": "archive.tgz",
					"out":     "1.txt",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.archiveName, func(t *testing.T) {
			h := setup.NewTest(setup.Platform(platform.Linux))
			path := h.Path
			fs := h.FS

			rootedFiles := []string{}
			for _, f := range files {
				filePath := path.Join(metadata.PackageVersionPath, f.Name)
				require.Nil(t, fs.WriteFile(filePath, []byte(f.Content), 0644))
				rootedFiles = append(rootedFiles, filePath)
			}

			// setup
			logger := log.Memory()
			factory := archive.NewFactory(fs, path)

			provider, err := factory.Select(test.archiveName)
			require.Nil(t, err)

			// create the test archive
			archivePath := path.Join(metadata.PackageVersionPath, test.archiveName)
			err = provider.Archive(archivePath, rootedFiles...)
			require.Nil(t, err)

			// cleanup so when we roundtrip we see the actual files
			for _, f := range rootedFiles {
				err = fs.Remove(f)
				require.Nil(t, err)
			}

			extract := actions.NewExtractProvider(factory, path, logger)
			require.NotNil(t, provider)

			err = extract.Execute(test.action, metadata)
			require.Nil(t, err, errorStringOrDefault(err))

			for _, f := range files {
				filePath := path.Join(metadata.PackageVersionPath, f.Name)
				ok, err := fs.Exists(filePath)
				require.Nil(t, err)
				require.True(t, ok, "file %s does not exist", filePath)
				bytes, err := fs.ReadFile(filePath)
				require.Nil(t, err)
				require.Equal(t, string(bytes), f.Content)
			}
		})
	}
}

func errorStringOrDefault(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
