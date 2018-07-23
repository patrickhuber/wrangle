package archiver

import (
	"compress/gzip"
	"fmt"
	"io"

	"github.com/spf13/afero"
)

// https://github.com/mholt/archiver/blob/master/targz.go
type targzArchiver struct {
	fileSystem afero.Fs
}

// NewTargzArchiver returns a new targz archiver
func NewTargzArchiver(fileSystem afero.Fs) Archiver {
	return &targzArchiver{fileSystem: fileSystem}
}

func (archive *targzArchiver) Write(output io.Writer, filePaths []string) error {
	return writeTarGz(archive.fileSystem, filePaths, output, "")
}

func writeTarGz(fileSystem afero.Fs, filePaths []string, output io.Writer, dest string) error {
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
