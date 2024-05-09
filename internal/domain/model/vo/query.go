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
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"net/url"
	"sort"
	"strings"
)

// Query represents an HTTP query. It is a type alias for map[string][]string,
// where the key is the parameter name and the value is a slice of parameter values.
// The Query type provides various methods to manipulate and interact with the query
// parameters, such as adding, setting, deleting, filtering, and encoding.
type Query map[string][]string

func NewEmptyQuery() Query {
	return Query{}
}

// NewQuery takes an HTTP query and creates a new instance of Query.
func NewQuery(httpQuery url.Values) Query {
	return Query(httpQuery)
}

// Add appends the values to the slice associated with the given key in a new copy of the Query map.
// If the key does not exist in the original Query map, it creates a new key-value pair with the provided values.
// The new copy of the Query map is then returned.
func (q Query) Add(key string, values []string) Query {
	newQuery := q.copy()
	newQuery[key] = append(newQuery[key], values...)
	return newQuery
}

// Append appends the values to the slice associated with the given key in a new copy of the Query map.
// If the key does not exist in the original Query map, it returns the original Query map without any modifications.
// The new copy of the Query map is then returned.
func (q Query) Append(key string, values []string) Query {
	if q.NotExists(key) {
		return q
	}
	return q.Add(key, values)
}

// Set sets the slice of values for the given key in a new copy of the Query map.
// If the key does not exist in the original Query map, it creates a new key-value pair with the provided values.
// The new copy of the Query map is then returned.
func (q Query) Set(key string, values []string) Query {
	newQuery := q.copy()
	newQuery[key] = values
	return newQuery
}

// Replace replaces the slice of values associated with the given key in a new copy of the Query map.
// If the key does not exist in the original Query map, it returns the original Query map without any modifications.
// The new copy of the Query map is then returned.
func (q Query) Replace(key string, values []string) Query {
	if _, ok := q[key]; !ok {
		return q
	}
	return q.Set(key, values)
}

// Delete deletes the slice associated with the given key in a new copy of the Query map.
// If the key does not exist in the original Query map, the new copy remains unchanged.
// The new copy of the Query map is then returned.
func (q Query) Delete(key string) Query {
	newQuery := q.copy()
	delete(newQuery, key)
	return newQuery
}

// Exists checks if the given key exists in the Query map.
// It returns true if the key exists, and false otherwise.
func (q Query) Exists(key string) bool {
	_, ok := q[key]
	return ok
}

// NotExists checks if the given key does not exist in the Query map.
// It returns true if the key does not exist, and false if the key exists.
func (q Query) NotExists(key string) bool {
	return !q.Exists(key)
}

func (q Query) Projection(projection *Projection) Query {
	// se o objeto de valor estiver nil ou vazio, retornamos
	if helper.IsNil(projection) || projection.IsEmpty() {
		return q
	}
	// projetamos com base no tipo de projeção, All, Addition, Rejection
	if helper.Equals(projection.Type(), enum.ProjectionValueRejection) {
		return q.projectionRejection(projection)
	}
	// se não for rejection, ele é Addition ou All, que é a mesma regra
	return q.projectionAddition(projection)
}

func (q Query) Map(mapper *Mapper) Query {
	// se tiver nil ou vazio, retornamos o query atual
	if helper.IsNil(mapper) || mapper.IsEmpty() {
		return q
	}
	// instanciamos o novo query a ser mapeado
	mappedQuery := Query{}
	for key, value := range q {
		if mapper.Exists(key) {
			mappedQuery[mapper.Get(key)] = value
		} else {
			mappedQuery[key] = value
		}
	}
	// retornamos o query mapeado
	return mappedQuery
}

func (q Query) Modify(modifier *Modifier, httpRequest *HttpRequest, httpResponse *HttpResponse) Query {
	// instanciamos o valor a ser usado para modificar
	newValue := modifier.ValueAsSliceOfString(httpRequest, httpResponse)
	// modificamos com base no action
	switch modifier.Action() {
	case enum.ModifierActionAdd:
		return q.Add(modifier.Key(), newValue)
	case enum.ModifierActionApd:
		return q.Append(modifier.Key(), newValue)
	case enum.ModifierActionSet:
		return q.Set(modifier.Key(), newValue)
	case enum.ModifierActionRpl:
		return q.Replace(modifier.Key(), newValue)
	case enum.ModifierActionDel:
		return q.Delete(modifier.Key())
	default:
		return q
	}
}

// Encode encodes the values into “URL encoded” form
// ("bar=baz&foo=qux") sorted by key.
func (q Query) Encode() string {
	// se for vazio retornamos a string vazia
	if helper.IsEmpty(q) {
		return ""
	}

	// instanciamos o valor string a ser usado para adicionar os valores
	var buf strings.Builder

	// obtemos as keys
	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
	}
	// fazemos o sort
	sort.Strings(keys)

	// iteramos as chaves ordenadas
	for _, k := range keys {
		// obtemos o valor da chave
		vs := q[k]
		// fazemos o sort dos valores
		sort.Strings(vs)

		// escapamos a chave
		keyEscaped := url.QueryEscape(k)

		// iteramos sobre os valores pela chave ja ordenados
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
		}
	}
	// retornamos o valor da query como string
	return buf.String()
}

// copy creates a shallow copy of the Query map.
// It iterates over each key-value pair in the original Query map and assigns them to the new copy.
// The copied Query map is then returned.
func (q Query) copy() Query {
	copiedQuery := Query{}
	for key, value := range q {
		copiedQuery[key] = value
	}
	return copiedQuery
}

func (q Query) projectionAddition(projection *Projection) Query {
	// inicializamos o query vazio
	projectedQuery := Query{}
	// iteramos o query atual
	for key, value := range q {
		// se a key atual conter na projeção com o valor 1, adicionamos
		if projection.IsAddition(key) {
			projectedQuery[key] = value
		}
	}
	// retornamos o novo query
	return projectedQuery
}

func (q Query) projectionRejection(projection *Projection) Query {
	// iniciamos o valor do query copiando os valores originais
	projectedQuery := q.copy()
	// iteramos o query atual
	for key := range q {
		// se ele estiver na projeção, removemos
		if projection.Exists(key) {
			delete(projectedQuery, key)
		}
	}
	// retornamos o novo query
	return projectedQuery
}
