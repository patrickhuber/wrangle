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
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Debugln(v ...interface{})
	Trace(v ...interface{})
	Tracef(format string, v ...interface{})
	Traceln(v ...interface{})
	Level() Level
}

// Default returns the default platform logger
func Default(options ...LogOption) Logger {
	l := &logger{
		Logger: log.Default(),
		level:  ErrorLevel,
	}
	for _, opt := range options {
		opt(l)
	}
	return l
}

type logger struct {
	*log.Logger
	level Level
}

func shouldLog(level Level, target Level) bool {
	return level >= target
}

func (l *logger) Debug(v ...interface{}) {
	if !shouldLog(l.level, DebugLevel) {
		return
	}
	l.Print(v...)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	if !shouldLog(l.level, DebugLevel) {
		return
	}
	l.Printf(format, v...)
}

func (l *logger) Debugln(v ...interface{}) {
	if !shouldLog(l.level, DebugLevel) {
		return
	}
	l.Println(v...)
}

func (l *logger) Trace(v ...interface{}) {
	if !shouldLog(l.level, TraceLevel) {
		return
	}
	l.Print(v...)
}

func (l *logger) Tracef(format string, v ...interface{}) {
	if !shouldLog(l.level, TraceLevel) {
		return
	}
	l.Printf(format, v...)
}

func (l *logger) Traceln(v ...interface{}) {
	if !shouldLog(l.level, TraceLevel) {
		return
	}
	l.Println(v...)
}

func (l *logger) Level() Level {
	return l.level
}

type LogOption func(l *logger)

func SetLevel(level Level) LogOption {
	return func(l *logger) {
		l.level = level
	}
}
