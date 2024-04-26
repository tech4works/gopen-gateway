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

// MiddlewaresJson represents a map of string keys to BackendJson values. It is used to configure and store middleware
// settings in the Gopen application.
// Each key in the map represents the name of the middleware, and the corresponding BackendJson value defines the
// properties and behavior of that middleware.
type MiddlewaresJson map[string]BackendJson

// Middlewares is a type that represents a map of string keys to Backend values.
// It is used to configure and store middleware settings in the Gopen server.
// Each key in the map represents the name of the middleware, and the corresponding Backend
// value defines the properties and behavior of that middleware.
type Middlewares map[string]Backend

// newMiddlewares creates a new instance of Middlewares based on the provided MiddlewaresJson parameter.
// It iterates over the keys and values of the MiddlewaresJson map and creates a new Backend instance for each value,
// using the newBackend function.
// The newly created Backends are assigned to the corresponding keys in the Middlewares map.
// The Middlewares map is then returned.
func newMiddlewares(middlewaresJsonVO MiddlewaresJson) (m Middlewares) {
	m = Middlewares{}
	for k, v := range middlewaresJsonVO {
		m[k] = newBackend(&v)
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
