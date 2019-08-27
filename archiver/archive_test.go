package archiver_test

import (
	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/filesystem"
)

type testFile struct {
	folder, name, content string
}

func createFiles(fileSystem filesystem.FileSystem, files []testFile) error {
	for _, f := range files {
		path := filepath.Join(f.folder, f.name)
		path = filepath.ToSlash(path)
		err := fileSystem.Write(path, []byte(f.content), 0666)
		if err != nil {
			return err
		}
	}
	return nil
}
