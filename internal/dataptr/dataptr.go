package dataptr

type DataPointer struct {
	Segments []Segment
}

type Segment interface {
	segment()
}

type Constraint struct {
	Segment
	Value string
	Key   string
}

type Element struct {
	Segment
	Name string
}

type Index struct {
	Segment
	Index int
}
