package archiver

import "io"

// Archiver defines an archive interface for reading and writing to an archive
type Archiver interface {
	Write(output io.Writer, paths []string) error
	Read(input io.Reader, destination string) error
}
