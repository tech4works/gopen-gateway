package middleware

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
)

type trace struct {
	traceProvider infra.TraceProvider
}

type Trace interface {
	Do(req *api.Request)
}

// NewTrace creates a new Trace instance.
func NewTrace(traceProvider infra.TraceProvider) Trace {
	return trace{
		traceProvider: traceProvider,
	}
}

// Do perform the tracing logic for the request.
// It adds the X-Forwarded-For header to the request with the remote address,
// and sets the X-TraceId header if it is not already specified.
// Then it proceeds to the next function in the request.
func (t trace) Do(req *api.Request) {
	// adicionamos na requisição o X-Forwarded-For
	req.AddHeader(consts.XForwardedFor, req.RemoteAddr())
	// caso não tenha trace id informado, setamos
	if helper.IsEmpty(req.HeaderValue(consts.XTraceId)) {
		req.SetHeader(consts.XTraceId, t.traceProvider.GenerateTraceId())
	}
	// seguimos para a próxima func da requisição
	req.Next()
}
