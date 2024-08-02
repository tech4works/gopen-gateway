package log

import (
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"github.com/tech4works/gopen-gateway/internal/domain/todo/interfaces/log"
)

type provider struct {
}

func NewProvider() log.Provider {
	return provider{}
}

func (p provider) Endpoint(endpoint *vo.Endpoint, traceID, clientIP string) log.Console {
	return newEndpointConsole(endpoint, traceID, clientIP)
}

func (p provider) Backend(backend *vo.Backend, httpRequest *vo.HttpRequest) log.Console {
	return newBackendConsole(backend, httpRequest)
}
