package ilog

import (
	"bytes"
	"io"
	"log"
)

func Memory() Logger {
	var out bytes.Buffer
	return MemoryWith(&out)
}

func MemoryWith(writer io.Writer) Logger {
	return log.New(writer, "", log.LstdFlags)
}
