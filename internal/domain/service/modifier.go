package service

import (
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
)

type modifier struct {
}

type Modifier interface {
	ExecuteInRequestContext(executeData vo.ExecuteModifierInRequestContext) vo.Request
	ExecuteInResponseContext(executeData vo.ExecuteModifierInResponseContext) vo.Response
}

func NewModifier() Modifier {
	return modifier{}
}

func (m modifier) ExecuteInRequestContext(executeData vo.ExecuteModifierInRequestContext) vo.Request {
	// instanciamos o backendRequestVO atual
	requestVO := executeData.Request()
	currentBackendRequestVO := requestVO.CurrentBackendRequest()

	// instanciamos o responseVO para ser usado como parâmetro nos modifies
	responseVO := executeData.Response()

	// inicializamos o header a ser modificado ou não, mas para construir o novo VO
	globalHeader := requestVO.Header()
	localHeader := currentBackendRequestVO.Header()

	// iteramos para modificar os headers solicitados
	for _, modifierVO := range executeData.ModifierHeader() {
		// executamos a modificação do cabeçalho no objeto de valor instanciado
		modifyVO := vo.NewModifyHeaders(modifierVO, &globalHeader, &localHeader, requestVO, responseVO)
		modifyVO.Execute()
	}

	// inicializamos a url a ser modificado ou não, mas para construir o novo VO
	localPath := currentBackendRequestVO.Path()
	// inicializamos os parâmetros a ser modificado ou não, mas para construir o novo VO
	globalParams := requestVO.Params()
	localParams := currentBackendRequestVO.Params()

	// iteramos para modificar os parâmetros solicitados
	for _, modifierVO := range executeData.ModifierParams() {
		// executamos a modificação dos parâmetros no objeto de valor instanciado
		modifyVO := vo.NewModifyParams(modifierVO, &globalParams, &localParams, &localPath, requestVO, responseVO)
		modifyVO.Execute()
	}

	// inicializamos as queries a ser modificado ou não, mas para construir o novo VO
	globalQuery := requestVO.Query()
	localQuery := currentBackendRequestVO.Query()

	// iteramos para modificar as queries solicitadas
	for _, modifierVO := range executeData.ModifierQuery() {
		// executamos a modificação das queries no objeto de valor instanciado
		modifyVO := vo.NewModifyQueries(modifierVO, &globalQuery, &localQuery, requestVO, responseVO)
		modifyVO.Execute()
	}

	// inicializamos os bodies a ser modificado ou não, mas para construir o novo VO
	globalBody := requestVO.Body()
	localBody := currentBackendRequestVO.Body()

	// iteramos para modificar os bodies solicitados
	for _, bodyModifierVO := range executeData.ModifierBody() {
		// executamos a modificação dos bodies no objeto de valor instanciado
		modifyVO := vo.NewModifyBodies(bodyModifierVO, &globalBody, &localBody, requestVO, responseVO)
		modifyVO.Execute()
	}

	// criamos o novo objeto de valor da solicitação backend
	backendRequestVO := currentBackendRequestVO.Modify(localPath, localHeader, localParams, localQuery, localBody)

	// criamos o novo objeto de valor da solicitação com os novos valores e com backendRequest atualizado
	return requestVO.Modify(globalHeader, globalParams, globalQuery, globalBody, backendRequestVO)
}

func (m modifier) ExecuteInResponseContext(executeData vo.ExecuteModifierInResponseContext) vo.Response {
	// instanciamos o vo de resposta
	responseVO := executeData.Response()
	// instanciamos o vo de response atual
	lastBackendResponseVO := responseVO.LastBackendResponse()

	// instanciamos o vo de request usado para passar como parâmetro
	requestVO := executeData.Request()

	// inicializamos o header a ser modificado ou não, mas para construir o novo VO
	globalHeader := responseVO.Header()
	localHeader := lastBackendResponseVO.Header()

	// iteramos para modificar os headers solicitados
	for _, modifierVO := range executeData.ModifierHeader() {
		// executamos a modificação do cabeçalho no objeto de valor instanciado
		modifyVO := vo.NewModifyHeaders(modifierVO, &globalHeader, &localHeader, requestVO, responseVO)
		modifyVO.Execute()
	}

	// inicializamos os bodies a ser modificado ou não, mas para construir o novo VO
	globalBody := responseVO.Body()
	localBody := lastBackendResponseVO.Body()

	// iteramos para modificar os bodies solicitados
	for _, modifierVO := range executeData.ModifierBody() {
		// executamos a modificação dos bodies no objeto de valor instanciado
		modifyVO := vo.NewModifyBodies(modifierVO, &globalBody, &localBody, requestVO, responseVO)
		modifyVO.Execute()
	}

	// inicializamos o status code a ser modificado ou não
	globalStatusCode := responseVO.StatusCode()
	localStatusCode := lastBackendResponseVO.StatusCode()

	if executeData.ModifierStatusCode().Valid() {
		modifierVO := executeData.ModifierStatusCode()

		// executamos a modificação dos statusCodes no objeto de valor instanciado
		modifyVO := vo.NewModifyStatusCodes(modifierVO, &globalStatusCode, &localStatusCode, requestVO, responseVO)
		modifyVO.Execute()
	}

	// construímos o novo backendResponseVO modificado
	backendResponseVO := lastBackendResponseVO.Modify(localStatusCode, localHeader, localBody)

	// construímos o novo responseVO modificado
	return responseVO.Modify(globalStatusCode, globalHeader, globalBody, backendResponseVO)
}
