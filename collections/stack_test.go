package collections_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/collections"
)

var _ = Describe("Stack", func() {
	Describe("Push", func() {
		It("can push onto stack", func() {
			stack := collections.NewStack()
			stack.Push(1)
			Expect(stack.Count()).To(Equal(1))
		})
		It("can push multiple items", func() {
			stack := collections.NewStack()
			stack.Push(1)
			Expect(stack.Count()).To(Equal(1))
			stack.Push(1)
			Expect(stack.Count()).To(Equal(2))
			stack.Push(1)
			Expect(stack.Count()).To(Equal(3))
		})
	})
	Describe("Pop", func() {
		It("can pop", func() {
			stack := collections.NewStack()
			stack.Push(1)
			stack.Push(2)
			stack.Push(3)
			Expect(stack.Count()).To(Equal(3))
			Expect(stack.Pop()).To(Equal(3))
			Expect(stack.Pop()).To(Equal(2))
			Expect(stack.Pop()).To(Equal(1))
		})
		It("returns null when empty", func() {
			stack := collections.NewStack()
			stack.Push(1)
			Expect(stack.Pop()).To(Equal(1))
			Expect(stack.Pop()).To(BeNil())
		})
	})
})
