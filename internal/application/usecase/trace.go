package usecase

import "github.com/google/uuid"

type trace struct {
}

type Trace interface {
	GenerateTraceId() string
}

func NewTrace() Trace {
	return trace{}
}

func (t trace) GenerateTraceId() string {
	return uuid.NewString()[:8]
}
