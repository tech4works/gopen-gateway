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

// NewModifier creates and returns a new Modifier instance.
func NewModifier() Modifier {
	return modifier{}
}

// Execute applies modifications to the request and response objects based on the given executeData.
//
// It iterates through the executeData's modifier collections and performs the corresponding modifications
// on the request and response objects, if the modifiers' contexts are valid. The method returns the modified
// request and response objects.
//
// Parameters:
//   - executeData : vo.ExecuteModifier
//     The executeData object containing the request, response, and modifier collections.
//
// Returns:
// - vo.Request: The potentially altered request value object
// - vo.Response: The potentially altered response value object
func (m modifier) Execute(executeData vo.ExecuteModifier) (vo.Request, vo.Response) {
	// instanciamos o requestVO para ser modificado ou não
	requestVO := executeData.Request()
	// instanciamos o responseVO para ser modificado ou não
	responseVO := executeData.Response()

	// executamos o modificador de código de status
	if executeData.ModifierStatusCode().Valid() &&
		executeData.ModifierStatusCode().EqualsContext(enum.ModifierContextResponse) {
		modifyVO := vo.NewModifyStatusCodes(executeData.ModifierStatusCode(), requestVO, responseVO)
		requestVO, responseVO = modifyVO.Execute()
	}

	// executamos os modificadores de cabeçalho
	modifierHeader := executeData.ModifierHeader()
	requestVO, responseVO = m.modify(modifierHeader, executeData.Context(), requestVO, responseVO, vo.NewHeaders)

	// executamos os modificadores de parâmetros
	modifierParams := executeData.ModifierParams()
	requestVO, responseVO = m.modify(modifierParams, executeData.Context(), requestVO, responseVO, vo.NewModifyParams)

	// executamos os modificadores de queries
	modifierQuery := executeData.ModifierQuery()
	requestVO, responseVO = m.modify(modifierQuery, executeData.Context(), requestVO, responseVO, vo.NewModifyQueries)

	// executamos os modificadores de body
	modifierBody := executeData.ModifierBody()
	requestVO, responseVO = m.modify(modifierBody, executeData.Context(), requestVO, responseVO, vo.NewModifyBodies)

	// retornamos os objetos de valore
	return requestVO, responseVO
}

// The method modify iterates over a list of provided modifiers and applies them to the request
// and response value objects if the modifier is valid, and it matches the given context.
// A Modification strategy is created for each individual valid and matching modifier
// Then, this strategy is executed, potentially altering the provided request and response value objects.
//
// Parameters:
// - modifiers: A slice of Modifier value objects to iterate over and potentially apply
// - context: The current context that incoming modifiers must match to be applied
// - requestVO: Request value object that may be modified by the execution of a modifier strategy
// - responseVO: Response value object that may be modified by the execution of a modifier strategy
// - newModifyVO: function to create a new ModifyLastBackendResponse value object (a modification strategy)
//
// Returns:
// - vo.Request: The potentially altered request value object
// - vo.Response: The potentially altered response value object
func (m modifier) modify(modifiers []vo.Modifier, context enum.ModifierContext, requestVO vo.Request,
	responseVO vo.Response, newModifyVO vo.NewModifyVOFunc) (vo.Request, vo.Response) {
	// iteramos os modificadores
	for _, modifierVO := range modifiers {
		// caso ele seja invalido ou não tiver no context vamos para o próximo
		if modifierVO.Invalid() || modifierVO.NotEqualsContext(context) {
			continue
		}
		// damos o new modify vo para instanciar a estratégia
		strategy := newModifyVO(modifierVO, requestVO, responseVO)
		// executamos a estrátegia, substituímos os objetos de valor modificados, ou não
		requestVO, responseVO = strategy.Execute()
	}
	// retornamos os objetos de valor modificados ou não
	return requestVO, responseVO
}
