package items

import (
	"io"

	"github.com/patrickhuber/wrangle/store"
)

type Renderer interface {
	RenderItems([]store.Item, io.Writer) error
	Name() string
}
