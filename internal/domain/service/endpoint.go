package service

import (
	"context"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
)

// endpoint is a struct type that represents an `endpoint` service domain in the Gopen server.
//
// It contains a backendService field of type Backend, which encapsulates the functionality
// for interacting with a backend service.
type endpoint struct {
	backendService Backend
}

// Endpoint represents an interface for executing a specific endpoint in the Gopen server.
// It defines the Execute method, which takes a context and an ExecuteEndpoint object as parameters
// and returns a Response object.
type Endpoint interface {
	// Execute executes a specific endpoint in the Gopen server.
	// It takes a context and a vo.ExecuteEndpoint object as parameters and returns a vo.Response object.
	Execute(ctx context.Context, executeData vo.ExecuteEndpoint) vo.Response
}

// NewEndpoint returns a new instance of the endpoint struct with the provided backendService.
func NewEndpoint(backendService Backend) Endpoint {
	return endpoint{
		backendService: backendService,
	}
}

// Execute executes the endpoint with the given executeData and returns the response.
//
// The method initializes the necessary value objects and processes the middlewares and backends
// of the endpoint. It iterates through the middleware keys provided by the endpoint's beforeware
// and afterware fields, executes the configured middleware backends, and updates the request and
// response value objects.
//
// If the response object indicates that the response needs to be aborted, the method returns the
// abort response without further processing.
//
// The method then iterates through the main backends of the endpoint, executes them, and updates
// the request and response value objects. Again, if the response object indicates that the response
// needs to be aborted, the method returns the abort response.
//
// After processing the backends, the method processes the middlewares configured in the
// afterware field, updating the request and response value objects accordingly.
//
// Finally, the method returns the final response value object.
//
// Parameters:
//   - ctx: The context.Context object for the execution.
//   - executeData: The vo.ExecuteEndpoint object containing the necessary data for execution, including
//     the Gopen, Endpoint, and Request value objects.
//
// Returns:
// The vo.Response object representing the response of the executed endpoint.
func (e endpoint) Execute(ctx context.Context, executeData vo.ExecuteEndpoint) vo.Response {
	// instanciamos o objeto gopenVO
	gopenVO := executeData.Gopen()
	// instanciamos o objeto de valor do endpoint
	endpointVO := executeData.Endpoint()
	// instanciamos o objeto de valor da requisição
	requestVO := executeData.Request()
	// inicializamos o objeto de valor de resposta do serviço
	responseVO := vo.NewResponse(endpointVO)

	// iteramos o beforeware, chaves configuradas para middlewares antes das requisições principais
	requestVO, responseVO = e.processMiddlewares(ctx, gopenVO, "beforeware", endpointVO.Beforeware(),
		requestVO, responseVO)
	// verificamos a resposta precisa ser abortada
	if abortResponseVO := responseVO.AbortResponse(); helper.IsNotNil(abortResponseVO) {
		return *abortResponseVO
	}

	// iteramos os backends principais para executa-las
	requestVO, responseVO = e.processBackends(ctx, endpointVO.Backends(), requestVO, responseVO)
	// verificamos a resposta precisa ser abortada
	if abortResponseVO := responseVO.AbortResponse(); helper.IsNotNil(abortResponseVO) {
		return *abortResponseVO
	}

	// iteramos o afterware, chaves configuradas para middlewares depois das requisições principais
	requestVO, responseVO = e.processMiddlewares(ctx, gopenVO, "afterware", endpointVO.Afterware(),
		requestVO, responseVO)
	// verificamos a resposta precisa ser abortada
	if abortResponseVO := responseVO.AbortResponse(); helper.IsNotNil(abortResponseVO) {
		return *abortResponseVO
	}

	// retornamos o objeto de valor de resposta final
	return responseVO
}

// processMiddlewares processes the middleware backends for the given middleware keys.
// It iterates through the middleware keys and checks if each key is configured in the gopenVO
// middlewares field. If configured, it creates an executeBackendVO object and executes the
// backend service using the backendService.Execute method.
// If the response object indicates that the response needs to be aborted, the method breaks
// out of the loop and returns the current request and response value objects.
// The method returns the updated request and response value objects after executing all the
// configured middleware backends.
//
// Parameters:
//   - ctx: The context.Context object for the execution.
//   - gopenVO: The vo.Gopen object containing the middleware configurations.
//   - middlewareType: The type of middleware being processed (beforeware or afterware).
//   - middlewareKeys: The middleware keys to be processed.
//   - requestVO: The vo.Request object representing the current request.
//   - responseVO: The vo.Response object representing the current response.
//
// Returns:
// The vo.Request and vo.Response objects after executing the middleware backends.
func (e endpoint) processMiddlewares(
	ctx context.Context,
	gopenVO vo.Gopen,
	middlewareType string,
	middlewareKeys []string,
	requestVO vo.Request,
	responseVO vo.Response,
) (vo.Request, vo.Response) {
	// iteramos as chaves de middlewares
	for _, middlewareKey := range middlewareKeys {
		// verificamos se essa chave foram configuradas no campo middlewares
		middlewareBackendVO, ok := gopenVO.Middleware(middlewareKey)
		if !ok {
			logger.Warning(middlewareType, middlewareKey, "not configured on middlewares field!")
			continue
		}
		// instanciamos o objeto de valor de execução do backend
		executeBackendVO := vo.NewExecuteBackend(middlewareBackendVO, requestVO, responseVO)
		// processamos o backend do middleware
		requestVO, responseVO = e.backendService.Execute(ctx, executeBackendVO)
		// verificamos a resposta precisa ser abortada
		if responseVO.IsAbortResponse() {
			break
		}
	}
	// retornamos os novos objetos de valor response e request
	return requestVO, responseVO
}

// processBackends iterates through the provided backends and executes each backend.
// It updates the request and response value objects accordingly. If the response object
// indicates that the response needs to be aborted, the iteration stops and the current
// request and response value objects are returned.
//
// Parameters:
//   - ctx: The context.Context object for the execution.
//   - backends: The slice of vo.Backend objects representing the backends to be processed.
//   - requestVO: The vo.Request object representing the current request.
//   - responseVO: The vo.Response object representing the current response.
//
// Returns:
// The updated vo.Request and vo.Response objects after executing the backends.
func (e endpoint) processBackends(ctx context.Context, backends []vo.Backend, requestVO vo.Request, responseVO vo.Response,
) (vo.Request, vo.Response) {
	// iteramos os backends fornecidos
	for _, backendVO := range backends {
		// processamos o backend principal iterado
		requestVO, responseVO = e.backendService.Execute(ctx, vo.NewExecuteBackend(backendVO, requestVO, responseVO))
		// verificamos a resposta precisa ser abortada
		if responseVO.IsAbortResponse() {
			break
		}
	}
	return requestVO, responseVO
}
