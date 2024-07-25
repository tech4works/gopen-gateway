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

type Middlewares map[string]Backend

func newMiddlewares(middlewaresJson MiddlewaresJson) (m Middlewares) {
	m = Middlewares{}
	for k, v := range middlewaresJson {
		m[k] = newBackend(&v)
	}
	return m
}

func (m Middlewares) Get(key string) (*Backend, bool) {
	backend, ok := m[key]
	if !ok {
		return nil, false
	}
	return newMiddlewareBackend(&backend), true
}
