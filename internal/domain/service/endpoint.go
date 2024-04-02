package service

import (
	"context"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
)

type endpoint struct {
	backendService Backend
}

type Endpoint interface {
	Execute(ctx context.Context, executeData vo.ExecuteEndpoint) vo.Response
}

func NewEndpoint(backendService Backend) Endpoint {
	return endpoint{
		backendService: backendService,
	}
}

// Execute executes the endpoint operation with the given executeData and returns the response.
// It processes the beforeware, main backends, and afterware in order.
// The beforeware and afterware are configured as middleware keys in the endpointVO.
// If the response needs to be aborted, it returns the abortResponseVO.
// Otherwise, it returns the final responseVO.
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
