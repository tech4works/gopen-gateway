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

// NewTraceProvider creates a new instance of trace that implements the TraceProvider interface.
func NewTraceProvider() TraceProvider {
	return trace{}
}

// GenerateTraceId generates a unique trace ID by combining a UUID and the current UnixNano timestamp.
// It uses the UUID package to generate a random UUID. It then retrieves the current UnixNano timestamp
// using the time package. The UUID string is concatenated with the timestamp and truncated to 16 characters
// using the fmt package.
// The resulting trace ID is returned as a string.
// Example:
//
//	t := trace{}
//	traceId := t.GenerateTraceId()
//	fmt.Println(traceId)
//
// Output:
//
//	4ae6c92d16089e521626
func (t trace) GenerateTraceId() string {
	u := uuid.New().String()
	unixNano := time.Now().UnixNano()
	return fmt.Sprintf("%s%d", u[:8], unixNano)[:16]
}
