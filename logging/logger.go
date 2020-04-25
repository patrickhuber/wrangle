package logging

import (
	"log"
	"os"
)

type logger struct {
	err *log.Logger
	out *log.Logger
}

// Logger creates a logger package that writes output to error or output depending on the log type
type Logger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
}

// Default returns the default logger which logs to stdout and stderr
func Default() Logger {
	return &logger{
		err: log.New(os.Stderr, "err", 0),
		out: log.New(os.Stdout, "out", 0),
	}
}

// With creates a logger with the given out and error loggers
func With(out *log.Logger, err *log.Logger) Logger {
	return &logger{
		err: err,
		out: out,
	}
}

func (l *logger) Fatal(v ...interface{}) {
	l.err.Fatal(v...)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.err.Fatalf(format, v...)
}

func (l *logger) Fatalln(v ...interface{}) {
	l.err.Fatalln(v...)
}

func (l *logger) Print(v ...interface{}) {
	l.out.Print(v...)
}

func (l *logger) Printf(format string, v ...interface{}) {
	l.out.Printf(format, v...)
}

func (l *logger) Println(v ...interface{}) {
	l.out.Println(v...)
}
