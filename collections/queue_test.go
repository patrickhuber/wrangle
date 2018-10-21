package collections

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Queue", func() {
	It("can enqueue one item", func() {
		q := NewQueue()
		q.Enqueue(1)
		Expect(q.Count()).To(Equal(1))
	})
	It("can dequeue empty", func() {
		q := NewQueue()
		item := q.Dequeue()
		Expect(item).To(BeNil())
	})
	It("dequeues items first in first out", func() {
		q := NewQueue()
		q.Enqueue(1)
		q.Enqueue(2)
		q.Enqueue(3)
		Expect(q.Count()).To(Equal(3))
		Expect(q.Dequeue()).To(Equal(1))
		Expect(q.Dequeue()).To(Equal(2))
		Expect(q.Dequeue()).To(Equal(3))
	})
	Describe("Empty", func() {
		When("has items", func() {
			It("returns false", func() {
				q := NewQueue()
				Expect(q.Empty()).To(BeTrue())
				q.Enqueue(1)
				Expect(q.Empty()).To(BeFalse())
				q.Dequeue()
				Expect(q.Empty()).To(BeTrue())
			})
		})
		When("has no items", func() {
			It("returns true", func() {
				q := NewQueue()
				Expect(q.Empty()).To(BeTrue())
			})
		})
	})

})
