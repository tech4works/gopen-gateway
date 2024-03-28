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

func (e endpoint) Execute(ctx context.Context, executeData vo.ExecuteEndpoint) vo.Response {
	// instanciamos o objeto de valor do endpoint
	endpointVO := executeData.Endpoint()
	// instanciamos o objeto de valor da requisição
	requestVO := executeData.Request()
	// inicializamos o objeto de valor de resposta do serviço
	responseVO := vo.NewResponse(endpointVO)

	// iteramos o beforeware, chaves configuradas para middlewares antes das requisições principais
	for _, beforewareKey := range endpointVO.Beforeware() {
		// verificamos se essa chave foram configuradas no campo middlewares
		beforewareVO, ok := executeData.Gopen().Middleware(beforewareKey)
		if !ok {
			logger.Warning("beforeware", beforewareKey, "not configured on middlewares field!")
			continue
		}

		// processamos o backend de beforeware
		requestVO, responseVO = e.backendService.Execute(ctx, vo.NewExecuteBackend(beforewareVO, requestVO, responseVO))

		// verificamos a resposta precisa ser abortada
		if abortResponseVO := responseVO.AbortResponse(); helper.IsNotNil(abortResponseVO) {
			return *abortResponseVO
		}
	}

	// iteramos os backends principais para executa-las
	for _, backendVO := range endpointVO.Backends() {
		// processamos o backend principal iterado
		requestVO, responseVO = e.backendService.Execute(ctx, vo.NewExecuteBackend(backendVO, requestVO, responseVO))

		// verificamos a resposta precisa ser abortada
		if abortResponseVO := responseVO.AbortResponse(); helper.IsNotNil(abortResponseVO) {
			return *abortResponseVO
		}
	}

	// iteramos o afterware, chaves configuradas para middlewares depois das requisições principais
	for _, afterwareKey := range endpointVO.Afterware() {
		// verificamos se essa chave foram configuradas no campo middlewares
		afterwareVO, ok := executeData.Gopen().Middleware(afterwareKey)
		if !ok {
			logger.Warning("afterware", afterwareKey, "not configured on middlewares field!")
			continue
		}

		// processamos o backend de afterware
		requestVO, responseVO = e.backendService.Execute(ctx, vo.NewExecuteBackend(afterwareVO, requestVO, responseVO))

		// verificamos a resposta precisa ser abortada
		if abortResponseVO := responseVO.AbortResponse(); helper.IsNotNil(abortResponseVO) {
			return *abortResponseVO
		}
	}

	return responseVO
}
