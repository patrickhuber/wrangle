package collections

type queueNode struct {
	item     interface{}
	previous *queueNode
}

type queue struct {
	count int
	tail  *queueNode
	head  *queueNode
}

// Queue defines a queue
type Queue interface {
	Enqueue(item interface{})
	Dequeue() interface{}
	Count() int
	Empty() bool
}

// NewQueue creates a new Queue
func NewQueue() Queue {
	return &queue{
		count: 0,
	}
}

func (q *queue) Enqueue(item interface{}) {
	n := &queueNode{
		item: item,
	}
	if q.count == 0 {
		q.head = n
		q.tail = n
		q.count++
		return
	}
	q.count++
	q.head.previous = n
	q.head = n
}

func (q *queue) Dequeue() interface{} {
	if q.count == 0 {
		return nil
	}
	item := q.tail.item
	q.tail = q.tail.previous
	q.count--
	return item
}

func (q *queue) Count() int {
	return q.count
}

func (q *queue) Empty() bool {
	return q.count == 0
}
