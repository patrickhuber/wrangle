package host

import (
	"io"

	"github.com/patrickhuber/go-di"
)

type Host interface {
	io.Closer
	Container() di.Container
}
