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
	"fmt"
	"github.com/tech4works/checker"
	"strings"
)

type URLPath struct {
	value  string
	params Params
}

func NewURLPath(path string, paramValues map[string]string) URLPath {
	filteredParams := map[string]string{}
	for key, value := range paramValues {
		if checker.Contains(path, patternParamPathKey(key)) {
			filteredParams[key] = value
		}
	}
	return URLPath{
		value:  path,
		params: NewParams(filteredParams),
	}
}

func (u URLPath) Raw() string {
	return u.value
}

func (u URLPath) String() string {
	urlPathStr := u.value
	for _, key := range u.params.Keys() {
		pathKey := patternParamPathKey(key)
		urlPathStr = strings.ReplaceAll(urlPathStr, pathKey, u.params.Get(key))
	}
	return urlPathStr
}

func (u URLPath) Params() Params {
	return u.params
}

func (u URLPath) Exists(key string) bool {
	return checker.Contains(u.value, patternParamPathKey(key)) && u.params.Exists(key)
}

func (u URLPath) NotExists(key string) bool {
	return !u.Exists(key)
}

func patternParamPathKey(key string) string {
	return fmt.Sprintf(":%s", key)
}
