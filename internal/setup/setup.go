package setup

import (
	"io"

	"github.com/patrickhuber/di"
)

type Setup interface {
	io.Closer
	Container() di.Container
}
