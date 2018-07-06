package archiver

import "io"

type Archiver interface {
	Write(output io.Writer, paths []string) error
	Read(input io.Reader, destination string) error
}
