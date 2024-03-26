package service

import (
	"context"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"net/http"
)

type backend struct {
	modifierService Modifier
	restTemplate    RestTemplate
}

type Backend interface {
	Execute(ctx context.Context, executeData vo.ExecuteBackend) (vo.Request, vo.Response)
}

func NewBackend(modifierService Modifier, restTemplate RestTemplate) Backend {
	return backend{
		modifierService: modifierService,
		restTemplate:    restTemplate,
	}
}

func (b backend) Execute(ctx context.Context, executeData vo.ExecuteBackend) (vo.Request, vo.Response) {
	// instanciamos o objeto de valor de resposta
	responseVO := executeData.Response()

	// construímos o backend request
	requestVO := b.buildBackendRequest(executeData)

	// locamos o objeto de valor
	backendRequestVO := executeData.Request().CurrentBackendRequest()

	// montamos o http request com o context
	httpRequest, err := backendRequestVO.Http(ctx)
	// caso ocorra um erro na montagem, retornamos
	if helper.IsNotNil(err) {
		return requestVO, executeData.Response().Err(requestVO.Url(), err)
	}

	// chamamos a interface de infra para chamar a conexão http e tratar a resposta
	httpResponse, err := b.restTemplate.MakeRequest(httpRequest)
	// caso ocorra um erro, retornamos o response como abort = true e a resposta formatada
	if helper.IsNotNil(err) {
		return requestVO, responseVO.Err(requestVO.Url(), err)
	}
	// chamamos para fechar o body assim que possível
	defer b.closeBodyResponse(httpResponse)

	// construímos o objeto de valor de resposta do backend
	responseVO = b.buildBackendResponse(executeData, httpResponse)

	// retornamos o requestVO e responseVO gerados e utilizados
	return requestVO, responseVO
}

func (b backend) buildBackendRequest(executeData vo.ExecuteBackend) vo.Request {
	// instanciamos o objeto de valor de request
	requestVO := executeData.Request()

	// instanciamos o objeto de valor backend
	backendVO := executeData.Backend()

	// obtemos o host do backend todo: ter um sub-dominio de balancer
	balancedHost := backendVO.Host()

	// montamos o objeto de valor com os dados montados no meu serviço de domínio
	backendRequestVO := vo.NewBackendRequest(backendVO, balancedHost, executeData.Request())

	// criamos um novo objeto de valor de solicitação com o novo backendRequestVO e substituímos a request vo atual
	requestVO = requestVO.Append(backendRequestVO)

	// chamamos o sub-dominio para modificar as requisições tanto de backend como a request global
	return b.modifierService.ExecuteInRequestContext(vo.NewExecuteModifierInRequestContext(executeData, requestVO))
}

func (b backend) closeBodyResponse(response *http.Response) {
	err := response.Body.Close()
	if helper.IsNotNil(err) {
		logger.WarningSkipCaller(2, "Error close http response:", err)
	}
}

func (b backend) buildBackendResponse(executeData vo.ExecuteBackend, httpResponse *http.Response) vo.Response {
	// construímos o novo objeto de valor da resposta do backend
	backendResponseVO := vo.NewBackendResponse(executeData.Backend(), httpResponse)

	// adicionamos o novo backend request no objeto de valor de resposta
	responseVO := executeData.Response().Append(backendResponseVO)

	// chamamos o sub-dominio para modificar a resposta do backend
	return b.modifierService.ExecuteInResponseContext(vo.NewExecuteModifierInResponseContext(executeData, responseVO))
}
