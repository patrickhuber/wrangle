package archiver

import "io"

// Archiver defines an archive interface for reading and writing to an archive
type Archiver interface {
	Archive(output io.Writer, paths []string) error
	Extract(input io.Reader, filter string, destination string) error
}
