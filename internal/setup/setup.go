package setup

import (
	"io"

	"github.com/patrickhuber/wrangle/pkg/di"
)

type Setup interface {
	io.Closer
	Container() di.Container
}
