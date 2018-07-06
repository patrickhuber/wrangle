package archiver_test

import (
	"path/filepath"

	"github.com/spf13/afero"
)

type testFile struct {
	folder, name, content string
}

func createFiles(fileSystem afero.Fs, files []testFile) {
	for _, f := range files {
		path := filepath.Join(f.folder, f.name)
		path = filepath.ToSlash(path)
		afero.WriteFile(fileSystem, path, []byte(f.content), 0444)
	}
}
