package ilog

import (
	"fmt"
	"strings"
)

type Level uint32

const (
	FatalLevel Level = iota
	ErrorLevel
	WarningLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

func (l Level) String() string {
	switch l {
	case FatalLevel:
		return "fatal"
	case ErrorLevel:
		return "error"
	case WarningLevel:
		return "warn"
	case InfoLevel:
		return "info"
	case DebugLevel:
		return "debug"
	case TraceLevel:
		return "trace"
	}
	return "unknown"
}

func ParseLevel(level string) (Level, error) {
	level = strings.TrimSpace(strings.ToLower(level))
	switch level {
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn":
		return WarningLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "trace":
		return TraceLevel, nil
	}
	return FatalLevel, fmt.Errorf("unknown level '%s'", level)
}
