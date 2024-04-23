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
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
)

type modifier struct {
}

// Modifier is an interface that represents a modifier which can be executed on a backend request and response.
type Modifier interface {
	// Execute is a method that executes a modifier on a backend request and response.
	// It takes a vo.ExecuteModifier as a parameter, which contains the modifier context,
	// backend modifiers, request, and response. It returns a vo.Request object and a vo.Response object.
	// This method allows modification of the request and response based on the provided modifier.
	Execute(executeData *vo.ExecuteModifier) (*vo.Request, *vo.Response)
}

// NewModifier creates and returns a new Modifier instance.
func NewModifier() Modifier {
	return modifier{}
}

// Execute executes the given modifier on the provided backend request and response.
// The modifier is only applied if it is valid and matches the given context.
// If the modifier is applicable, a Modification strategy is created and executed,
// potentially modifying the request and response.
//
// Parameters:
//   - executeData: An ExecuteModifier object containing the modifier context, backend modifiers,
//     request, and response to be modified.
//
// Returns:
// - *vo.Request: The potentially altered request object.
// - *vo.Response: The potentially altered response object.
func (m modifier) Execute(executeData *vo.ExecuteModifier) (*vo.Request, *vo.Response) {
	// checamos se o backendModifier veio nil e ja retornamos
	if helper.IsNil(executeData.BackendModifiers()) {
		return executeData.Request(), executeData.Response()
	}

	// instanciamos o requestVO para ser modificado ou não
	requestVO := executeData.Request()
	// instanciamos o responseVO para ser modificado ou não
	responseVO := executeData.Response()

	// executamos o modificador de código de status
	if helper.IsNotEmpty(executeData.ModifierStatusCode()) &&
		helper.Equals(enum.ModifierContextResponse, executeData.Context()) {
		modifyVO := vo.NewModifyStatusCodes(executeData.ModifierStatusCode(), requestVO, responseVO)
		requestVO, responseVO = modifyVO.Execute()
	}

	// executamos os modificadores de cabeçalho
	modifierHeader := executeData.ModifierHeader()
	requestVO, responseVO = m.modify(modifierHeader, executeData.Context(), requestVO, responseVO, vo.NewHeaders)

	// executamos os modificadores de parâmetros
	modifierParam := executeData.ModifierParam()
	requestVO, responseVO = m.modify(modifierParam, executeData.Context(), requestVO, responseVO, vo.NewModifyParam)

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
func (m modifier) modify(modifiers []vo.Modifier, context enum.ModifierContext, requestVO *vo.Request,
	responseVO *vo.Response, newModifyVO vo.NewModifyVOFunc) (*vo.Request, *vo.Response) {
	// iteramos os modificadores
	for _, modifierVO := range modifiers {
		// caso ele seja invalido ou não tiver no context vamos para o próximo
		if modifierVO.Invalid() || modifierVO.NotEqualsContext(context) {
			continue
		}
		// damos o new modify vo para instanciar a estratégia
		strategy := newModifyVO(&modifierVO, requestVO, responseVO)
		// executamos a estrátegia, substituímos os objetos de valor modificados, ou não
		requestVO, responseVO = strategy.Execute()
	}
	// retornamos os objetos de valor modificados ou não
	return requestVO, responseVO
}
