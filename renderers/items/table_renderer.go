package items

import (
	"io"

	"github.com/patrickhuber/wrangle/store"
)

type tableRenderer struct {
}

func (r *tableRenderer) RenderItems(itemList []store.Item, writer io.Writer) error {
	return nil
}

func (r *tableRenderer) Name() string {
	return "table"
}

func NewTableRenderer() Renderer {
	return &tableRenderer{}
}
