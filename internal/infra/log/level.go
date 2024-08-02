package log

import (
	"fmt"
	"github.com/GabrielHCataldo/go-logger/logger"
)

type level string

const (
	InfoLevel  level = "INF"
	DebugLevel level = "DBG"
	WarnLevel  level = "WRN"
	ErrorLevel level = "ERR"
)

func (l level) String() string {
	return fmt.Sprint(l.color(), string(l), logger.StyleReset)
}

func (l level) color() string {
	switch l {
	case InfoLevel:
		return logger.ForegroundBlue
	case DebugLevel:
		return logger.ForegroundCyan
	case WarnLevel:
		return logger.ForegroundYellow
	case ErrorLevel:
		return logger.ForegroundRed
	default:
		return logger.StyleReset
	}
}
