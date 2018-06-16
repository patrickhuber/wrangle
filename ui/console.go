package ui

import (
	"bytes"
	"io"
	"os"
)

// Console defines an interface for iteracting with a console
type Console interface {
	Out() io.Writer
	Error() io.Writer
	In() io.Reader
}

type console struct {
	out io.Writer
	err io.Writer
	in  io.Reader
}

// NewOSConsole creates a new console with os.Stdout, os.Stderr and os.Stdin
func NewOSConsole() Console {
	return &console{
		out: os.Stdout,
		err: os.Stderr,
		in:  os.Stdin,
	}
}

// NewMemoryConsole creates a new memory console with bytes.Buffer for each member
func NewMemoryConsole() Console {
	return &console{
		out: &bytes.Buffer{},
		err: &bytes.Buffer{},
		in:  &bytes.Buffer{},
	}
}

func (console *console) Out() io.Writer {
	return console.out
}

func (console *console) Error() io.Writer {
	return console.err
}

func (console *console) In() io.Reader {
	return console.in
}
