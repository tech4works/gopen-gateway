/*
 * Copyright 2024 Tech4Works
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
	"net/url"
	"sort"
	"strings"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
)

type Query struct {
	values map[string][]string
}

func NewQuery(values map[string][]string) Query {
	return Query{
		values: values,
	}
}

func NewEmptyQuery() Query {
	return Query{}
}

func (q Query) GetAll(key string) []string {
	return q.values[key]
}

func (q Query) Exists(key string) bool {
	_, ok := q.values[key]
	return ok
}

func (q Query) NotExists(key string) bool {
	return !q.Exists(key)
}

func (q Query) Length() int {
	return len(q.values)
}

func (q Query) IsEmpty() bool {
	return q.Length() == 0
}

func (q Query) Keys() (ss []string) {
	for key := range q.values {
		ss = append(ss, key)
	}
	return ss
}

func (q Query) Encode() string {
	if q.IsEmpty() {
		return ""
	}

	orderedKeys := q.Keys()
	sort.Strings(orderedKeys)

	var strBuilder strings.Builder
	for _, key := range orderedKeys {
		valueByKey := q.GetAll(key)
		sort.Strings(valueByKey)

		keyEscaped := url.QueryEscape(key)
		for _, value := range valueByKey {
			if checker.IsGreaterThan(strBuilder.Len(), 0) {
				strBuilder.WriteByte('&')
			}
			strBuilder.WriteString(keyEscaped)
			strBuilder.WriteByte('=')
			strBuilder.WriteString(url.QueryEscape(value))
		}
	}
	return strBuilder.String()
}

func (q Query) Map() any {
	return q.Copy()
}

func (q Query) Copy() map[string][]string {
	copiedMap := map[string][]string{}
	for key, value := range q.values {
		copiedMap[key] = value
	}
	return copiedMap
}

func (q Query) String() string {
	return converter.ToString(q.values)
}
