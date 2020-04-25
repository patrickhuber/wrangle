package packages

import "io"

// InterfaceReader reads from the given reader and returns an interface struct
type InterfaceReader interface {
	Read(reader io.Reader) (interface{}, error)
}
