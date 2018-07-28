package archiver

import (
	"compress/gzip"
	"fmt"
	"io"

	"github.com/spf13/afero"
)

// https://github.com/mholt/archiver/blob/master/targz.go
type tgzArchiver struct {
	fileSystem afero.Fs
}

// NewTargzArchiver returns a new targz archiver
func NewTargzArchiver(fileSystem afero.Fs) Archiver {
	return &tgzArchiver{fileSystem: fileSystem}
}

func (archive *tgzArchiver) Archive(output io.Writer, filePaths []string) error {
	return writeTarGz(archive.fileSystem, filePaths, output)
}

func writeTarGz(fileSystem afero.Fs, filePaths []string, output io.Writer) error {
	gzw := gzip.NewWriter(output)
	defer gzw.Close()

	return writeTar(fileSystem, filePaths, gzw)
}

func (archive *tgzArchiver) Extract(input io.Reader, filter string, destination string) error {
	gzr, err := gzip.NewReader(input)
	if err != nil {
		return fmt.Errorf("error decompressing: %v", err)
	}
	defer gzr.Close()

	return NewTarArchiver(archive.fileSystem).Extract(gzr, filter, destination)
}
