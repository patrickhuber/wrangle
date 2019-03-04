package items

import (
	"io"

	"github.com/patrickhuber/wrangle/store"
)

type treeRenderer struct {
}

func (r *treeRenderer) RenderItems(itemList []store.Item, writer io.Writer) error {
	return nil
}

func (r *treeRenderer) Name() string {
	return "tree"
}

func NewTreeRenderer() Renderer {
	return &treeRenderer{}
}
