package store_test

import (
	patch "github.com/cppforlife/go-patch/patch"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("", func() {
	It("can find key value", func() {
		pointer, err := patch.NewPointerFromString("/key1")
		Expect(err).To(BeNil())
		doc := map[interface{}]interface{}{
			"key1": "abc",
			"key2": "xyz",
		}
		response, err := patch.FindOp{Path: pointer}.Apply(doc)
		Expect(err).To(BeNil())
		Expect(response).To(Equal("abc"))
	})
	It("can create pointer", func() {
		ptr, err := patch.NewPointerFromString("/some/path")
		Expect(err).To(BeNil())
		Expect(len(ptr.Tokens())).To(Equal(3))
	})
})
