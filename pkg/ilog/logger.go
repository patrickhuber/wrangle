package ilog

import (
	"log"
)

// Logger defines a logging interface
type Logger interface {
	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalln(v ...interface{})
}

// Default returns the default platform logger
func Default() Logger {
	return log.Default()
}
