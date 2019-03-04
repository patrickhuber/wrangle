package items

import (
	"io"

	"github.com/patrickhuber/wrangle/store"
)

type jsonRenderer struct {
}

func (r *jsonRenderer) RenderItems(itemList []store.Item, writer io.Writer) error {
	return nil
}

func (r *jsonRenderer) Name() string {
	return "json"
}

func NewJsonRenderer() Renderer {
	return &jsonRenderer{}
}
