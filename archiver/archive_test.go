package archiver_test

import (
	"github.com/patrickhuber/wrangle/filepath"

	"github.com/spf13/afero"
)

type testFile struct {
	folder, name, content string
}

func createFiles(fileSystem afero.Fs, files []testFile) error {
	for _, f := range files {
		path := filepath.Join(f.folder, f.name)
		path = filepath.ToSlash(path)
		err := afero.WriteFile(fileSystem, path, []byte(f.content), 0666)
		if err != nil {
			return err
		}
	}
	return nil
}
