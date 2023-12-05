package dataptr

import (
	"fmt"
	"strconv"
	"strings"
)

type DataPointer struct {
	Segments []Segment
}

type Segment interface {
	segment()
	fmt.Stringer
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

func (dp DataPointer) String() string {
	builder := strings.Builder{}
	for i, seg := range dp.Segments {
		if i > 0 {
			builder.WriteRune('/')
		}
		builder.WriteString(seg.String())
	}
	return builder.String()
}

func (c Constraint) String() string {
	return fmt.Sprintf("%s=%s", c.Key, c.Value)
}

func (e Element) String() string {
	return e.Name
}

func (i Index) String() string {
	return strconv.Itoa(i.Index)
}
