package archiver

import (
	"github.com/patrickhuber/wrangle/filesystem"
	"compress/gzip"
	"fmt"
)

// https://github.com/mholt/archiver/blob/master/targz.go
type tgzArchiver struct {
	fileSystem filesystem.FileSystem
}

// NewTargz returns a new targz archiver
func NewTargz(fileSystem filesystem.FileSystem) Archiver {
	return &tgzArchiver{fileSystem: fileSystem}
}

func (tgz *tgzArchiver) Archive(archive string, filePaths []string) error {
	return tgz.writeTarGz(archive, filePaths)
}

func (tgz *tgzArchiver) writeTarGz(archive string, filePaths []string) error {
	file, err := tgz.fileSystem.Create(archive)
	if err != nil {
		return err
	}

	gzw := gzip.NewWriter(file)
	defer gzw.Close()

	return NewTarArchiver(tgz.fileSystem).ArchiveWriter(gzw, filePaths)
}

func (tgz *tgzArchiver) Extract(archive string, destination string, files []string) error {

	file, err := tgz.fileSystem.Open(archive)
	if err != nil {
		return err
	}

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("error decompressing: %v", err)
	}
	defer gzr.Close()

	return NewTarArchiver(tgz.fileSystem).ExtractReader(gzr, destination, files)
}
