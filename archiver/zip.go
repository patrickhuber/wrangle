package archiver

import "io"

type zipArchive struct {
}

func (archive *zipArchive) Write(output io.Writer, filePaths []string) error {
	return nil
}

func (archive *zipArchive) Read(input io.Reader, destination string) error {
	return nil
}
