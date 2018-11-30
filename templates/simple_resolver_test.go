package templates

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SimpleResolver", func() {
	It("can create", func() {

		resolver, err := newSimpleResolver("key", "value", "key1", "value1")
		Expect(err).To(BeNil())
		Expect(resolver).ToNot(BeNil())

		value, err := resolver.Get("key")
		Expect(err).To(BeNil())
		Expect(value).ToNot(BeNil())
		Expect(value).To(Equal("value"))

		value, err = resolver.Get("key1")
		Expect(err).To(BeNil())
		Expect(value).ToNot(BeNil())
		Expect(value).To(Equal("value1"))
	})
})
