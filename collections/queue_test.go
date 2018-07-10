package collections

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	t.Run("CanEnqueueOneItem", func(t *testing.T) {
		r := require.New(t)
		q := NewQueue()
		q.Enqueue(1)
		r.Equal(1, q.Count())
	})
	t.Run("CanDequeueEmpty", func(t *testing.T) {
		r := require.New(t)
		q := NewQueue()
		item := q.Dequeue()
		r.Nil(item)
	})
	t.Run("ItemsAreInFirstInFirstOutOrder", func(t *testing.T) {
		r := require.New(t)
		q := NewQueue()
		q.Enqueue(1)
		q.Enqueue(2)
		q.Enqueue(3)
		r.Equal(3, q.Count())
		r.Equal(1, q.Dequeue())
		r.Equal(2, q.Dequeue())
		r.Equal(3, q.Dequeue())
	})
	t.Run("EmptyReturnsCorrectValue", func(t *testing.T) {
		r := require.New(t)
		q := NewQueue()
		r.True(q.Empty())
		q.Enqueue(1)
		r.False(q.Empty())
		q.Dequeue()
		r.True(q.Empty())
	})
}
