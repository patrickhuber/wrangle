package console

import (
	"bytes"
	"io"
)

type memory struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

func NewMemory() Console {
	return &memory{
		in:  &bytes.Buffer{},
		out: &bytes.Buffer{},
		err: &bytes.Buffer{},
	}
}

func (c *memory) In() io.Reader {
	return c.in
}

func (c *memory) Out() io.Writer {
	return c.out
}

func (c *memory) Error() io.Writer {
	return c.err
}
