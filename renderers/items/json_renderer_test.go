package items_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/renderers/items"
	"github.com/patrickhuber/wrangle/store"
)

var _ = Describe("JsonRenderer", func() {
	It("should render json", func() {
		itemList := []store.Item{
			store.NewValueItem("hi", "hi"),
			store.NewValueItem("test", "test"),
		}
		r := items.NewJsonRenderer()
		buffer := bytes.Buffer{}
		err := r.RenderItems(itemList, &buffer)
		Expect(err).To(BeNil())
		Expect(buffer.String()).To(Equal(""))
	})
})
