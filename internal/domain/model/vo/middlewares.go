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

import "github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"

// Middlewares is a type that represents a map of string keys to Backend values.
// It is used to configure and store middleware settings in the Gopen server.
// Each key in the map represents the name of the middleware, and the corresponding Backend
// value defines the properties and behavior of that middleware.
type Middlewares map[string]Backend

// newMiddlewares creates a new Middlewares object based on the provided middlewaresDTO map.
// Each key-value pair in the middlewaresDTO map will be converted to a Backend object and added to the Middlewares object.
// The new Middlewares object will then be returned.
func newMiddlewares(middlewaresDTO map[string]dto.Backend) (m Middlewares) {
	m = Middlewares{}
	for k, v := range middlewaresDTO {
		m[k] = newBackend(v)
	}
	return m
}

// Get retrieves the Backend associated with the specified key from the Middlewares map.
// It returns the Backend if it exists, otherwise it returns an empty Backend and false.
// If the Backend exists, it creates a new modified Backend using the newMiddlewareBackend function,
// passing the existing Backend and a backendExtraConfig with the omitResponse set to true.
// It then returns the modified Backend and true.
func (m Middlewares) Get(key string) (Backend, bool) {
	backend, ok := m[key]
	if !ok {
		return Backend{}, false
	}
	return newMiddlewareBackend(&backend, &BackendExtraConfig{
		omitResponse: true,
	}), true
}
