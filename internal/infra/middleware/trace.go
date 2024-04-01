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

func NewTrace(traceProvider infra.TraceProvider) Trace {
	return trace{
		traceProvider: traceProvider,
	}
}

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
