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
	"net/url"
	"sort"
	"strings"
)

// Query represents an HTTP query. It is a type alias for map[string][]string,
// where the key is the parameter name and the value is a slice of parameter values.
// The Query type provides various methods to manipulate and interact with the query
// parameters, such as adding, setting, deleting, filtering, and encoding.
type Query map[string][]string

// NewQuery takes an HTTP query and creates a new instance of Query.
func NewQuery(httpQuery url.Values) Query {
	return Query(httpQuery)
}

// Add appends the values to the slice associated with the given key in a new copy of the Query map.
// If the key does not exist in the original Query map, it creates a new key-value pair with the provided values.
// The new copy of the Query map is then returned.
func (q Query) Add(key string, values []string) (r Query) {
	r = q.copy()
	r[key] = append(r[key], values...)
	return r
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
func (q Query) Set(key string, values []string) (r Query) {
	r = q.copy()
	r[key] = values
	return r
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

// Rename renames a key in a new copy of the Query map by replacing the old key with the new key.
// If the old key does not exist in the original Query map, the original Query map is returned unchanged.
// The new copy of the Query map is then returned.
func (q Query) Rename(oldKey, newKey string) (r Query) {
	if q.NotExists(oldKey) {
		return q
	}
	r = q.copy()
	r[newKey] = r[oldKey]
	delete(r, oldKey)
	return r
}

// Delete deletes the slice associated with the given key in a new copy of the Query map.
// If the key does not exist in the original Query map, the new copy remains unchanged.
// The new copy of the Query map is then returned.
func (q Query) Delete(key string) (r Query) {
	r = q.copy()
	delete(r, key)
	return r
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

func (q Query) Projection(keys []string) (r Query) {
	if helper.IsEmpty(keys) {
		return q
	}
	r = q.copy()
	for key := range q {
		if helper.NotContains(keys, key) {
			r = r.Delete(key)
		}
	}
	return r
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
func (q Query) copy() (r Query) {
	r = Query{}
	for key, value := range q {
		r[key] = value
	}
	return r
}
