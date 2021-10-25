package console

import (
	"io"
	"os"
)

type Console interface {
	In() io.Reader
	Out() io.Writer
	Error() io.Writer
}

type osConsole struct {
}

func NewOS() Console {
	return &osConsole{}
}

func (c *osConsole) In() io.Reader {
	return os.Stdin
}

func (c *osConsole) Out() io.Writer {
	return os.Stdout
}

func (c *osConsole) Error() io.Writer {
	return os.Stderr
}
