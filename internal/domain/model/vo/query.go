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

// NewEmptyQuery creates a new empty instance of the Query type.
// It returns an empty Query map with no key-value pairs.
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

// Projection modifies the Query map based on the provided Projection.
// If the projection is nil or empty, the original Query map is returned unchanged.
// If the projection type is ProjectionValueRejection, a new Query map is returned with the keys specified in the
// projection removed. Otherwise, a new Query map is returned with only the keys specified in the projection added.
func (q Query) Projection(projection *Projection) Query {
	if helper.IsNil(projection) || projection.IsEmpty() {
		return q
	} else if helper.Equals(projection.Type(), enum.ProjectionValueRejection) {
		return q.projectionRejection(projection)
	}
	return q.projectionAddition(projection)
}

// Map applies the provided mapper to each key-value pair in the Query map.
// If the mapper is nil or empty, the original Query map is returned unchanged.
// If a key exists in the mapper, the key is replaced with the mapped key in a new Query map,
// while the value remains the same. If a key does not exist in the mapper, the key-value pair is
// added to the new Query map as is.
// The new copy of the Query map with the applied mapper is then returned.
func (q Query) Map(mapper *Mapper) Query {
	if helper.IsNil(mapper) || mapper.IsEmpty() {
		return q
	}

	mappedQuery := Query{}
	for key, value := range q {
		if mapper.Exists(key) {
			mappedQuery[mapper.Get(key)] = value
		} else {
			mappedQuery[key] = value
		}
	}

	return mappedQuery
}

// Modify modifies the Query map based on the provided Modifier.
// It extracts the value as a slice of strings using the ValueAsSliceOfString method from the Modifier.
// It then performs a specific action based on the ActionType of the Modifier.
// The ActionType can be Add, Append, Set, Replace, or Delete.
//   - Add: appends the values to the slice associated with the given key in a new copy of the Query map.
//     If the key does not exist, it creates a new key-value pair with the provided values.
//   - Append: appends the values to the slice associated with the given key in a new copy of the Query map.
//     If the key does not exist, it returns the original Query map without any modifications.
//   - Set: sets the slice of values for the given key in a new copy of the Query map.
//     If the key does not exist, it creates a new key-value pair with the provided values.
//   - Replace: replaces the slice of values associated with the given key in a new copy of the Query map.
//     If the key does not exist, it returns the original Query map without any modifications.
//   - Delete: deletes the slice associated with the given key in a new copy of the Query map.
//     If the key does not exist, it returns the original Query map without any modifications.
//
// The new copy of the Query map with the applied modification is then returned.
func (q Query) Modify(modifier *Modifier, httpRequest *HttpRequest, httpResponse *HttpResponse) Query {
	newValue := modifier.ValueAsSliceOfString(httpRequest, httpResponse)
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

// Encode encodes the Query map into a URL-encoded string.
// If the Query map is empty, it returns an empty string.
// The keys in the Query map are sorted in lexicographical order.
// The values associated with each key are also sorted in lexicographical order.
// The keys and values are URL-encoded using url.QueryEscape.
// The URL-encoded key-value pairs are joined with '&' and returned as a single string.
func (q Query) Encode() string {
	if helper.IsEmpty(q) {
		return ""
	}

	keys := make([]string, 0, len(q))
	for key := range q {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var strBuilder strings.Builder
	for _, key := range keys {
		valuesByKey := q[key]
		sort.Strings(valuesByKey)

		keyEscaped := url.QueryEscape(key)
		for _, value := range valuesByKey {
			if strBuilder.Len() > 0 {
				strBuilder.WriteByte('&')
			}
			strBuilder.WriteString(keyEscaped)
			strBuilder.WriteByte('=')
			strBuilder.WriteString(url.QueryEscape(value))
		}
	}
	return strBuilder.String()
}

func (q Query) String() string {
	if helper.IsEmpty(q) {
		return ""
	}
	mapString := map[string]string{}
	for key, value := range q {
		mapString[key] = strings.Join(value, ",")
	}
	return helper.SimpleConvertToString(mapString)
}

// copy creates a new copy of the Query map by iterating over the key-value pairs
// of the original Query map and assigning them to the newly created Query map.
// The new copy of the Query map is then returned.
func (q Query) copy() Query {
	copiedQuery := Query{}
	for key, value := range q {
		copiedQuery[key] = value
	}
	return copiedQuery
}

// projectionAddition creates a new Query map by adding only the key-value pairs that are specified
// in the projection and exist in the original Query map. Other key-value pairs are not included.
// The new copy of the Query map with the projected key-value pairs is then returned.
func (q Query) projectionAddition(projection *Projection) Query {
	projectedQuery := Query{}
	for key, value := range q {
		if projection.IsAddition(key) {
			projectedQuery[key] = value
		}
	}
	return projectedQuery
}

// projectionRejection removes the key-value pairs from the Query map that are specified
// in the projection and exist in the original Query map. The modified Query map is returned.
func (q Query) projectionRejection(projection *Projection) Query {
	projectedQuery := q.copy()
	for key := range q {
		if projection.Exists(key) {
			delete(projectedQuery, key)
		}
	}
	return projectedQuery
}
