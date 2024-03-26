package infra

import (
	"fmt"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/external"
	"github.com/google/uuid"
	"time"
)

type trace struct {
}

func NewTraceProvider() external.TraceProvider {
	return trace{}
}

func (t trace) GenerateTraceId() string {
	u := uuid.New().String()
	unixNano := time.Now().UnixNano()

	return fmt.Sprintf("%s%d", u[:8], unixNano)[:16]
}
