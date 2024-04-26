/*
 * Copyright 2024 Gabriel Cataldo
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"context"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/interfaces"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"net/http"
)

// backend represents a type that encapsulates the functionality for interacting with a backend service.
// It provides methods for executing backend requests and handling backend responses.
type backend struct {
	modifierService Modifier
	restTemplate    interfaces.RestTemplate
}

// Backend represents a type that encapsulates the functionality for interacting with a backend service.
// It provides a method for executing backend requests and handling backend responses.
type Backend interface {
	// Execute is a method that is implemented by a Backend type. It represents the execution of a backend server request and
	// response. The method takes a context.Context object and a vo.ExecuteBackend object as parameters.
	// It returns the vo.Request and vo.Response objects associated with the ExecuteBackend object.
	//
	// Parameters:
	// ctx - The context.Context object that provides a context for the execution of the request.
	// executeData - The vo.ExecuteBackend object that encapsulates the data for the execution of the backend request.
	//
	// Returns:
	// vo.Request - The vo.Request associated with vo.ExecuteBackend and can be modified by the backend execution process.
	// vo.Response - The vo.Response associated with vo.ExecuteBackend and can be modified by the backend execution process.
	Execute(ctx context.Context, executeData *vo.ExecuteBackend) (*vo.Request, *vo.Response)
}

// NewBackend initializes and returns a new Backend instance.
//
// This function serves as a factory, accepting implementations of Modifier and interfaces.RestTemplate
// as arguments, creating and returning an instance of the backend type that satisfies the Backend interface.
//
// Parameters:
// modifierService: Provides the service for modifying backend information. Must conform to the Modifier interface.
// restTemplate: Provides the functionality for conducting RESTful operations. Must conform to the RestTemplate
// interface from interfaces package.
//
// Returns:
// A Backend instance with modifierService and restTemplate composed in.
func NewBackend(modifierService Modifier, restTemplate interfaces.RestTemplate) Backend {
	return backend{
		modifierService: modifierService,
		restTemplate:    restTemplate,
	}
}

// Execute sends an HTTP request to a backend based on the provided executeData.
// The function's steps are as follows:
//
//  1. The function constructs the backend request. This also includes a potential response modification.
//  2. The backend request gets converted to an HTTP request. If this operation fails then the error will be returned in
//     the response object.
//  3. The function performs an HTTP request by calling a REST client. If this operation fails,
//     then the error will be returned in the response object with its abort flag set to true.
//  4. Finally, the function creates a backend response object from the returned HTTP response.
//     This response is again able to include a request modification.
//
// The function returns the updated request and response value objects.
// An already constructed response object is returned if any error occurs during the function execution.
//
// If the function executes successfully, it ensures that the HTTP response body is closed.
//
// Parameters:
// ctx: the execution context.
// executeData: contains the information necessary to execute a backend request.
//
// Returns:
// The function returns two value objects:
// requestVO: the potentially modified backend request.
// responseVO: the backend response. If an error occurred, it contains the error information.
func (b backend) Execute(ctx context.Context, executeData *vo.ExecuteBackend) (*vo.Request, *vo.Response) {
	// construímos o backend request, junto pode vir uma possível alteração no response pelo modifier
	requestVO, responseVO := b.buildBackendRequest(executeData)

	// locamos o objeto de valor
	backendRequestVO := requestVO.CurrentBackendRequest()

	// montamos o http request com o context
	httpRequest, err := backendRequestVO.Http(ctx)
	// caso ocorra um erro na montagem, retornamos
	if helper.IsNotNil(err) {
		return requestVO, responseVO.Error(executeData.Endpoint().Path(), err)
	}

	// todo: adicionar log de execução do backend em questão

	// chamamos a interface de infra para chamar a conexão http e tratar a resposta
	httpResponse, err := b.restTemplate.MakeRequest(httpRequest)
	// caso ocorra um erro, retornamos o response como abort = true e a resposta formatada
	if helper.IsNotNil(err) {
		return requestVO, responseVO.Error(executeData.Endpoint().Path(), err)
	}
	// chamamos para fechar o body assim que possível
	defer b.closeBodyResponse(httpResponse)

	// construímos o objeto de valor de resposta do backend, junto pode vir uma possível alteração no request pelo modifier
	return b.buildBackendResponse(executeData.Backend(), requestVO, responseVO, httpResponse)
}

// buildBackendRequest is a method in the backend framework that uses executeData of type vo.ExecuteBackend.
// 1. Instantiate a request value object from executeData
// 2. Instantiate a backend value object from executeData
// 3. It retrieves the balanced host from the backendVO (todo: possibly using a subdomain balancer)
// 4. It constructs a new backendRequestVO object using backendVO, balanceHost and the initial request
// 5. Replaces the initial requestVO with a new version that includes the backendRequestVO
// 6. It invokes the Execute method of the modifierService to change the backend request and response and the actual request and response.
// The method returns a request value object and a possibly changed response value object.
func (b backend) buildBackendRequest(executeData *vo.ExecuteBackend) (*vo.Request, *vo.Response) {
	// instanciamos o objeto de valor de request
	requestVO := executeData.Request()

	// instanciamos o objeto de valor backend
	backendVO := executeData.Backend()

	// obtemos o host do backend
	balancedHost := backendVO.BalancedHost()

	// montamos o objeto de valor com os dados montados no meu serviço de domínio
	backendRequestVO := vo.NewBackendRequest(backendVO, balancedHost, executeData.Request())

	// criamos um novo objeto de valor de solicitação com o novo backendRequestVO e substituímos a request vo atual
	requestVO = requestVO.Append(backendRequestVO)

	// chamamos o sub-dominio para modificar as requisições tanto de backend como a própria request e a resposta
	// do backend e da propria response
	return b.modifierService.Execute(vo.NewExecuteRequestModifier(executeData.Backend(), requestVO, executeData.Response()))
}

// closeBodyResponse closes the HTTP response body.
// If there is an error while closing the body, a warning message will be logged.
func (b backend) closeBodyResponse(response *http.Response) {
	err := response.Body.Close()
	if helper.IsNotNil(err) {
		logger.WarningSkipCaller(2, "Error close http response:", err)
	}
}

// buildBackendResponse is a method in the backend framework that creates a new backend response object based on
// the provided parameters:
//
// Parameters:
//
//	backendVO: the backend value object.
//	requestVO: the request value object.
//	responseVO: the response value object.
//	httpResponse: the HTTP response object.
//
// Steps:
// 1. Constructs a new backend response value object using backendVO and httpResponse.
// 2. Appends the new backend request to the response value object.
// 3. Calls the modifierService's Execute method to modify the backend response.
// 4. If the response indicates abort, returns the request and an abort response.
// 5. If all steps are successful, returns the modified request and response.
//
// Returns:
//
//	The function returns two value objects:
//	- requestVO: the potentially modified backend request.
//	- responseVO: the backend response. If an error occurred, it contains the error information.
func (b backend) buildBackendResponse(
	backendVO *vo.Backend,
	requestVO *vo.Request,
	responseVO *vo.Response,
	httpResponse *http.Response,
) (*vo.Request, *vo.Response) {
	// construímos o novo objeto de valor da resposta do backend
	backendResponseVO := vo.NewBackendResponse(backendVO, httpResponse)

	// todo: adicionar log de resposta do backend em questão

	// adicionamos o novo backend request no objeto de valor de resposta
	responseVO = responseVO.Append(backendResponseVO)

	// chamamos o sub-dominio para modificar a resposta do backend
	requestVO, responseVO = b.modifierService.Execute(vo.NewExecuteResponseModifier(backendVO, requestVO, responseVO))

	// se resposta é para abortar retornamos
	if responseVO.Abort() {
		return requestVO, responseVO.AbortResponse()
	}

	// se tudo ocorrer bem retornamos a requisição e o response resultante
	return requestVO, responseVO
}
