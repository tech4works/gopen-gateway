package service

import (
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
)

type modifier struct {
}

type Modifier interface {
	Execute(executeData vo.ExecuteModifier) (vo.Request, vo.Response)
}

func NewModifier() Modifier {
	return modifier{}
}

func (m modifier) Execute(executeData vo.ExecuteModifier) (vo.Request, vo.Response) {
	// instanciamos o requestVO para ser modificado ou não
	requestVO := executeData.Request()

	// instanciamos o responseVO para ser modificado ou não
	responseVO := executeData.Response()

	// se tiver informado o status code modificamos
	if executeData.ModifierStatusCode().Valid() &&
		executeData.ModifierStatusCode().EqualsContext(enum.ModifierContextResponse) {
		// executamos a modificação dos statusCodes no objeto de valor instanciado
		modifyVO := vo.NewModifyStatusCodes(executeData.ModifierStatusCode(), requestVO, responseVO)
		requestVO, responseVO = modifyVO.Execute()
	}

	// iteramos para modificar os headers solicitados
	for _, modifierVO := range executeData.ModifierHeader() {
		// verificamos se ele ta no contexto correto
		if modifierVO.NotEqualsContext(executeData.Context()) {
			continue
		}

		// executamos a modificação do cabeçalho no objeto de valor instanciado
		modifyVO := vo.NewModifyHeaders(modifierVO, requestVO, responseVO)
		// retorna os objetos de valor modificado ou não
		requestVO, responseVO = modifyVO.Execute()
	}

	// iteramos para modificar os parâmetros solicitados
	for _, modifierVO := range executeData.ModifierParams() {
		// verificamos se ele ta no contexto correto
		if modifierVO.NotEqualsContext(executeData.Context()) {
			continue
		}

		// executamos a modificação dos parâmetros no objeto de valor instanciado
		modifyVO := vo.NewModifyParams(modifierVO, requestVO, responseVO)
		// retorna os objetos de valor modificado ou não
		requestVO, responseVO = modifyVO.Execute()
	}

	// iteramos para modificar as queries solicitadas
	for _, modifierVO := range executeData.ModifierQuery() {
		// verificamos se ele ta no contexto correto
		if modifierVO.NotEqualsContext(executeData.Context()) {
			continue
		}

		// executamos a modificação das queries no objeto de valor instanciado
		modifyVO := vo.NewModifyQueries(modifierVO, requestVO, responseVO)
		// retorna os objetos de valor modificado ou não
		requestVO, responseVO = modifyVO.Execute()
	}

	// iteramos para modificar os bodies solicitados
	for _, modifierVO := range executeData.ModifierBody() {
		// verificamos se ele ta no contexto correto
		if modifierVO.NotEqualsContext(executeData.Context()) {
			continue
		}

		// executamos a modificação dos bodies no objeto de valor instanciado
		modifyVO := vo.NewModifyBodies(modifierVO, requestVO, responseVO)
		// retorna os objetos de valor modificado ou não
		requestVO, responseVO = modifyVO.Execute()
	}

	// retornamos os objetos de valor manipulados ou não
	return requestVO, responseVO
}
