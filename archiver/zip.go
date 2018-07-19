package archiver

import (
	"io"

	"github.com/spf13/afero"
)

type zipArchive struct {
}

// NewZipArchiver creates a new zip archiver
func NewZipArchiver(fs afero.Fs) Archiver {
	return &zipArchive{}
}

func (archive *zipArchive) Write(output io.Writer, filePaths []string) error {
	return nil
}

func (archive *zipArchive) Read(input io.Reader, destination string) error {
	return nil
}
