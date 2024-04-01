package controller

import (
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/service"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
)

type endpoint struct {
	endpointService service.Endpoint
}

type Endpoint interface {
	Execute(req *api.Request)
}

func NewEndpoint(endpointService service.Endpoint) Endpoint {
	return endpoint{
		endpointService: endpointService,
	}
}

// Execute executes the service to process the endpoint.
// It takes a pointer to api.Request as an argument.
// It builds the service parameters using mapper.BuildExecuteServiceParams.
// It invokes the endpointService.Execute method passing the built parameters.
// It writes the response to the request using req.Write method.
func (e endpoint) Execute(req *api.Request) {
	// executamos o serviço de dominío para processar o endpoint
	responseVO := e.endpointService.Execute(mapper.BuildExecuteServiceParams(req))
	// respondemos a requisição a partir do objeto de valor recebido
	req.Write(responseVO)
}
