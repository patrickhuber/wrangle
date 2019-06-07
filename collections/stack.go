package collections

type stackNode struct {
	item interface{}
	next *stackNode
}

type stack struct {
	count int
	head  *stackNode
}

type Stack interface {
	Push(item interface{})
	Pop() interface{}
	Count() int
}

func NewStack() Stack {
	return &stack{
		count: 0,
	}
}

func (s *stack) Push(item interface{}) {
	n := &stackNode{
		item: item,
	}
	if s.count == 0 {
		s.head = n
		s.count++
		return
	}
	n.next = s.head
	s.head = n
	s.count++
}

func (s *stack) Pop() interface{} {
	if s.count == 0 {
		return nil
	}

	item := s.head.item
	s.count--
	s.head = s.head.next
	return item
}

func (s *stack) Peek() interface{} {
	if s.count == 0 {
		return nil
	}
	return s.head.item
}

func (s *stack) Count() int {
	return s.count
}
