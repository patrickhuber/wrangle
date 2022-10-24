package ilog

import (
	"bytes"
	"io"
	"log"
)

func Memory(options ...LogOption) Logger {
	var out bytes.Buffer
	return MemoryWith(&out, options...)
}

func MemoryWith(writer io.Writer, options ...LogOption) Logger {
	l := &logger{
		Logger: log.New(writer, "", log.LstdFlags),
	}
	for _, opt := range options {
		opt(l)
	}
	return l
}
