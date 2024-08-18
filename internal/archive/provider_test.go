package archive_test

import (
	"testing"

	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/go-xplat/setup"
	"github.com/patrickhuber/wrangle/internal/archive"
	"github.com/stretchr/testify/require"
)

type TestFile struct {
	Name    string
	Content string
}

func TestProvider(t *testing.T) {
	root := "/gran/parent/child"

	files := []*TestFile{
		{
			Name:    "1.txt",
			Content: "1",
		},
		{
			Name:    "2.txt",
			Content: "2",
		},
	}
	type test struct {
		archiveFile string
	}
	tests := []test{
		{
			archiveFile: "test.tar",
		},
		{
			archiveFile: "test.zip",
		},
		{
			archiveFile: "test.tgz",
		},
	}

	for _, test := range tests {
		t.Run(test.archiveFile, func(t *testing.T) {
			h := setup.NewTest(setup.Platform(platform.Linux))
			path := h.Path
			fs := h.FS
			factory := archive.NewFactory(fs, path)
			provider, err := factory.Select(test.archiveFile)
			require.Nil(t, err)

			var rootedFiles []string
			var names []string
			for _, f := range files {
				rootedFile := path.Join(root, f.Name)
				err := fs.WriteFile(rootedFile, []byte(f.Content), 0644)
				require.Nil(t, err)
				rootedFiles = append(rootedFiles, rootedFile)
				names = append(names, f.Name)
			}

			rootedArchiveFile := path.Join(root, test.archiveFile)
			require.Nil(t, provider.Archive(rootedArchiveFile, rootedFiles...))

			ok, err := fs.Exists(rootedArchiveFile)

			require.Nil(t, err)
			require.True(t, ok)

			for _, f := range rootedFiles {
				require.Nil(t, fs.Remove(f))
			}

			require.Nil(t, provider.Extract(rootedArchiveFile, root, names...))
			for _, f := range files {

				filePath := path.Join(root, f.Name)
				ok, err := fs.Exists(filePath)
				require.Nil(t, err)
				require.True(t, ok, "%s does not exist", filePath)

				b, err := fs.ReadFile(filePath)
				require.Nil(t, err)
				require.Equal(t, f.Content, string(b))
			}
		})
	}
}
