package infra

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type trace struct {
}

type TraceProvider interface {
	GenerateTraceId() string
}

func NewTraceProvider() TraceProvider {
	return trace{}
}

func (t trace) GenerateTraceId() string {
	u := uuid.New().String()
	unixNano := time.Now().UnixNano()

	return fmt.Sprintf("%s%d", u[:8], unixNano)[:16]
}
