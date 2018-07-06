package archiver

import (
	"compress/gzip"
	"fmt"
	"io"

	"github.com/patrickhuber/cli-mgr/filesystem"
)

// https://github.com/mholt/archiver/blob/master/targz.go
type targzArchiver struct {
	fileSystem filesystem.FsWrapper
}

// NewTargzArchiver returns a new targz archiver
func NewTargzArchiver(fileSystem filesystem.FsWrapper) Archiver {
	return &targzArchiver{fileSystem: fileSystem}
}

func (archive *targzArchiver) Write(output io.Writer, filePaths []string) error {
	return writeTarGz(archive.fileSystem, filePaths, output, "")
}

func writeTarGz(fileSystem filesystem.FsWrapper, filePaths []string, output io.Writer, dest string) error {
	gzw := gzip.NewWriter(output)
	defer gzw.Close()

	return writeTar(fileSystem, filePaths, gzw, dest)
}

func (archive *targzArchiver) Read(input io.Reader, destination string) error {
	gzr, err := gzip.NewReader(input)
	if err != nil {
		return fmt.Errorf("error decompressing: %v", err)
	}
	defer gzr.Close()

	return NewTarArchiver(archive.fileSystem).Read(gzr, destination)
}
