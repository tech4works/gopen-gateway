package log

import (
	"fmt"
)

type level string

const (
	InfoLevel  level = "INF"
	DebugLevel level = "DBG"
	WarnLevel  level = "WRN"
	ErrorLevel level = "ERR"
)

func (l level) String() string {
	return fmt.Sprint(l.color(), string(l), "\x1b[0m")
}

func (l level) color() string {
	switch l {
	case InfoLevel:
		return "\x1b[34m"
	case DebugLevel:
		return "\x1b[36m"
	case WarnLevel:
		return "\x1b[33m"
	case ErrorLevel:
		return "\x1b[31m"
	default:
		return "\x1b[0m"
	}
}
