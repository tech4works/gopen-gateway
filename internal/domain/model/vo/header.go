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

package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"net/http"
	"strings"
)

// Header represents a map of string keys to slices of string values.
type Header map[string][]string

// NewHeader creates a new Header object from an existing http.Header object.
func NewHeader(httpHeader http.Header) Header {
	return Header(httpHeader)
}

// newHeaderFailed creates a new Header object with failed status values for consts.XGopenCache, consts.XGopenComplete,
// and consts.XGopenSuccess.
func newHeaderFailed() Header {
	return Header{
		consts.XGopenCache:    {"false"},
		consts.XGopenComplete: {"false"},
		consts.XGopenSuccess:  {"false"},
	}
}

// newResponseHeader creates a new Header object with specific values for consts.XGopenCache, consts.XGopenComplete, consts.XGopenSuccess.
// The complete parameter is used to set the value of consts.XGopenComplete header.
// The success parameter is used to set the value of consts.XGopenSuccess header.
// The returned Header object contains the updated values for consts.XGopenCache, consts.XGopenComplete, consts.XGopenSuccess modifyHeaders.
func newResponseHeader(complete, success bool) Header {
	return Header{
		consts.XGopenCache:    {"false"},
		consts.XGopenComplete: {helper.SimpleConvertToString(complete)},
		consts.XGopenSuccess:  {helper.SimpleConvertToString(success)},
	}
}

// Http converts the Header object to an http.Header object.
// It returns the converted http.Header object.
func (h Header) Http() http.Header {
	return http.Header(h)
}

// AddAll accepts a key (in string format) and an array of values (in string format).
// It adds these values to the existing header under the provided key.
// The function makes a copy of the existing header before performing the operation
// to avoid mutating the original header.
func (h Header) AddAll(key string, values []string) Header {
	newHeader := h.copy()
	newHeader[key] = append(newHeader[key], values...)
	return newHeader
}

// Add is a method on the Header struct.
// It accepts a key and a value, both strings, as parameters.
// It copies the existing Header, appends the provided value to the slice
// of values associated with the provided key in the copied header,
// and then returns the updated header.
func (h Header) Add(key, value string) Header {
	newHeader := h.copy()
	newHeader[key] = append(newHeader[key], value)
	return newHeader
}

// Append appends the provided values to the existing values associated with the provided key in the header.
// If the key does not exist in the header, it returns the original header unchanged.
// It creates a new copy of the header, adds the values to the existing ones using the AddAll method,
// and then returns the modified header copy.
func (h Header) Append(key string, values []string) Header {
	if h.NotExists(key) {
		return h
	}
	return h.AddAll(key, values)
}

// Set is a method on the Header type. It takes a key and a value, both of type string, and returns a Header.
// The method makes a copy of the original Header, sets the value of the given key in the copied Header to a new
// string slice containing the provided value, and then returns the modified Header copy.
func (h Header) Set(key, value string) Header {
	newHeader := h.copy()
	newHeader[key] = []string{value}
	return newHeader
}

// SetAll is a method on the Header struct.
// It accepts a key (in string format) and an array of values (in string format) as
// parameters. It creates a copy of the existing Header, sets the value of the given key
// in the copied Header to the provided array of values, and then returns the modified
// Header copy.
func (h Header) SetAll(key string, values []string) Header {
	newHeader := h.copy()
	newHeader[key] = values
	return newHeader
}

// Replace replaces the values associated with the provided key in the Header object with the given values.
// If the key does not exist in the Header, it returns the original Header object without any changes.
func (h Header) Replace(key string, values []string) Header {
	if h.NotExists(key) {
		return h
	}
	return h.SetAll(key, values)
}

// Rename renames the given oldKey to the newKey in the Header object.
// If the oldKey does not exist in the Header, the method returns the original Header object.
// Otherwise, it creates a copy of the Header object, assigns the value of oldKey to newKey,
// deletes the oldKey, and returns the modified Header object.
func (h Header) Rename(oldKey, newKey string) Header {
	if h.NotExists(oldKey) {
		return h
	}
	newHeader := h.copy()
	newHeader[newKey] = newHeader[oldKey]
	delete(newHeader, oldKey)
	return newHeader
}

// Delete removes the value associated with the given key from the Header h.
// It returns a new Header object with the key removed.
// If the key does not exist in the Header, the returned Header is identical to the original.
func (h Header) Delete(key string) Header {
	newHeader := h.copy()
	delete(newHeader, key)
	return newHeader
}

// Get retrieves the value for a specific key from a Header. If a value exists,
// it concatenates its items with a comma separator and returns them as a string.
// If no value exists for the given key or the value is empty, it returns an empty string.
func (h Header) Get(key string) string {
	values := h[key]
	if helper.IsNotEmpty(values) {
		return strings.Join(values, ", ")
	}
	return ""
}

// Exists checks if a given key exists in the Header object.
// It returns true if the key exists, false otherwise.
func (h Header) Exists(key string) bool {
	_, ok := h[key]
	return ok
}

// NotExists checks if a given key does not exist in the Header object.
// It calls the Exists method of the Header object and returns the negation of the result.
// It returns true if the key does not exist, false otherwise.
func (h Header) NotExists(key string) bool {
	return !h.Exists(key)
}

func (h Header) ProjectionToRequest(projectionVO *Projection) Header {
	// se tiver nil ou vazio retornamos ele mesmo
	if helper.IsNil(projectionVO) || projectionVO.IsEmpty() {
		return h
	}
	// retornamos o novo header projetado segundo a config, ignorando os campos de requisição obrigatórios
	return h.projection(projectionVO, []string{consts.XForwardedFor, consts.XTraceId})
}

func (h Header) ProjectionToResponse(projectionVO *Projection) Header {
	// se tiver nil ou vazio retornamos ele mesmo
	if helper.IsNil(projectionVO) || projectionVO.IsEmpty() {
		return h
	}
	// retornamos o novo header projetado segundo a config, ignorando os campos de resposta obrigatórios
	return h.projection(projectionVO, []string{consts.XGopenSuccess, consts.XGopenCache, consts.XGopenCacheTTL,
		consts.XGopenComplete})
}

func (h Header) MapToRequest(mapperVO *Mapper) Header {
	// se o mapper estiver vazio, retornamos o header atual
	if helper.IsNil(mapperVO) || mapperVO.IsEmpty() {
		return h
	}
	// retornamos o novo header mapeado segundo a config, ignorando os campos de resposta obrigatórios
	return h.mapp(mapperVO, []string{consts.XForwardedFor, consts.XTraceId})
}

func (h Header) MapToResponse(mapperVO *Mapper) Header {
	// se o mapper estiver vazio, retornamos o header atual
	if helper.IsNil(mapperVO) || mapperVO.IsEmpty() {
		return h
	}
	// retornamos o novo header mapeado segundo a config, ignorando os campos de resposta obrigatórios
	return h.mapp(mapperVO, []string{consts.XGopenSuccess, consts.XGopenCache, consts.XGopenCacheTTL,
		consts.XGopenComplete})
}

// Aggregate combines the headers of two Header objects.
// It takes anotherHeader as a parameter and returns a new Header object that contains the combined headers.
// The method iterates through each key-value pair in anotherHeader and adds the values to the corresponding key in the new Header object.
// It uses the AddAll method to append the values to the existing ones, creating a new slice.
// The resulting Header object is returned at the end of the method.
func (h Header) Aggregate(anotherHeader Header) Header {
	aggregatedHeader := h.copy()
	for key, values := range anotherHeader {
		aggregatedHeader = aggregatedHeader.AddAll(key, values)
	}
	return aggregatedHeader
}

// copy creates a deep copy of the Header object.
// It returns a new Header object that is a copy of the original Header object.
func (h Header) copy() Header {
	copiedHeader := Header{}
	for key, value := range h {
		copiedHeader[key] = value
	}
	return copiedHeader
}

func (h Header) mapp(mapperVO *Mapper, ignoreKeys []string) Header {
	// inicializamos o novo header a ser retornado
	headerMapped := Header{}
	// iteramos o header atual para preencher o novo header com as chaves mapeadas
	for key, value := range h {
		// caso ele exista obtemos no mapper e não está na lista de chaves a serem ignoradas, adicionamos o novo nome
		if helper.NotContains(ignoreKeys, key) && mapperVO.Exists(key) {
			headerMapped[mapperVO.Get(key)] = value
		} else {
			headerMapped[key] = value
		}
	}
	// retornamos o header mapeado
	return headerMapped
}

func (h Header) projection(projectionVO *Projection, ignoreKeys []string) Header {
	// projetamos com base no tipo de projeção, All, Addition, Rejection
	if helper.Equals(projectionVO.Type(), enum.ProjectionTypeRejection) {
		return h.projectionRejection(projectionVO, ignoreKeys)
	}
	// se não for rejection, ele é Addition ou All, que é a mesma regra
	return h.projectionAddition(projectionVO, ignoreKeys)
}

func (h Header) projectionAddition(projectionVO *Projection, ignoreKeys []string) Header {
	// inicializamos o header
	projectedHeader := Header{}
	// iteramos o header atual
	for key, value := range h {
		// se ele contiver na lista ou na projeção como 1, adicionamos
		if helper.Contains(ignoreKeys, key) || projectionVO.IsAddition(key) {
			projectedHeader[key] = value
		}
	}
	// retornamos o novo header
	return projectedHeader
}

func (h Header) projectionRejection(projectionVO *Projection, ignoreKeys []string) Header {
	// iniciamos o valor do header copiando os valores originais
	projectedHeader := h.copy()
	// iteramos o header atual
	for key := range h {
		// se ele não contiver na lista a ser ignorada e estiver na projeção, removemos
		if helper.NotContains(ignoreKeys, key) && projectionVO.Exists(key) {
			delete(projectedHeader, key)
		}
	}
	// retornamos o novo header
	return projectedHeader
}
