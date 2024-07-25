package log

import (
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/todo/interfaces/log"
)

type noop struct {
}

func NewNoopConsole() log.Console {
	return noop{}
}

func (n noop) Infof(_ string, _ ...any) {
}

func (n noop) Info(_ ...any) {
}

func (n noop) Warnf(_ string, _ ...any) {
}

func (n noop) Warn(_ ...any) {
}

func (n noop) Errorf(_ string, _ ...any) {
}

func (n noop) Error(_ ...any) {
}
