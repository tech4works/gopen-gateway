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
)

// Gopen is a struct that represents the configuration for a Gopen server.
type Gopen struct {
	// port represents the port number on which the Gopen application will listen for incoming requests.
	// It is an integer value and can be specified in the Gopen configuration JSON file.
	port int
	// securityCors represents the configuration options for Cross-Origin Resource Sharing (CORS) settings in Gopen.
	securityCors *SecurityCors
	// middlewares is a map that represents the middleware configuration in Gopen.
	// The keys of the map are the names of the middlewares, and the values are
	// Backend objects that define the properties of each middleware.
	// The Backend struct contains fields like name, hosts, path, method, forwardHeaders,
	// forwardQueries, modifiers, and extraConfig, which specify the behavior
	// and settings of the middleware.
	middlewares Middlewares
	// endpoints is a field in the Gopen struct that represents a slice of Endpoint objects.
	// Each Endpoint object defines a specific API endpoint with its corresponding settings such as path, method,
	// timeout, limiter, cache, etc.
	endpoints []Endpoint
}

// NewGopen returns a new Gopen object based on the provided GopenJson configuration.
// It initializes the fields of the Gopen object using the values from the GopenJson object.
// The endpoints are created by iterating over the EndpointJson objects in the GopenJson endpoints slice,
// and calling the newEndpoint function to create an Endpoint object for each EndpointJson object.
// The Gopen object is then returned as a pointer.
func NewGopen(gopenJson *GopenJson) *Gopen {
	var endpoints []Endpoint
	for _, endpointJson := range gopenJson.Endpoints {
		endpoints = append(endpoints, newEndpoint(gopenJson, &endpointJson))
	}
	return &Gopen{
		port:         gopenJson.Port,
		securityCors: newSecurityCors(gopenJson.SecurityCors),
		middlewares:  newMiddlewares(gopenJson.Middlewares),
		endpoints:    endpoints,
	}
}

// Port returns the value of the port field in the Gopen struct.
func (g Gopen) Port() int {
	return g.port
}

// SecurityCors returns the value of the securityCors field in the Gopen struct.
func (g Gopen) SecurityCors() *SecurityCors {
	return g.securityCors
}

// Middleware retrieves a backend from the middlewares map based on the given key and returns it with a boolean
// indicating whether it exists or not. The returned backend is wrapped in a new middleware backend with the omitResponse
// field set to true.
func (g Gopen) Middleware(key string) (*Backend, bool) {
	return g.middlewares.Get(key)
}

// Endpoints returns the slice of Endpoint objects in the Gopen struct.
func (g Gopen) Endpoints() []Endpoint {
	return g.endpoints
}

// CountEndpoints returns the number of endpoints in the Gopen struct.
func (g Gopen) CountEndpoints() int {
	return len(g.endpoints)
}

// CountMiddlewares returns the number of middlewares in the Gopen struct. If the middlewares field is not nil,
// it returns the length of the middlewares map. Otherwise, it returns 0.
func (g Gopen) CountMiddlewares() int {
	if helper.IsNotNil(g.middlewares) {
		return len(g.middlewares)
	}
	return 0
}

// CountBackends returns the total number of backends across all endpoints in the Gopen struct.
// It iterates over each endpoint and gets the count of beforewares, backends, and afterwares.
// The sum of these counts is returned as the total count.
// This method is used to determine the total number of backends configured in the Gopen server.
// The count includes all backends across all endpoints, regardless of their status or visibility.
// It does not take into account any filtering or other conditions.
func (g Gopen) CountBackends() (count int) {
	for _, endpoint := range g.Endpoints() {
		count += endpoint.CountBeforewares()
		count += endpoint.CountBackends()
		count += endpoint.CountAfterwares()
	}
	return count
}

// CountAllDataTransforms recursively counts the number of data transforms in the Endpoint struct
// by adding the count of transforms in the Endpoint's response and the count of transforms in each backend.
// It returns the total count of transforms.
func (g Gopen) CountAllDataTransforms() (count int) {
	for _, endpoint := range g.Endpoints() {
		count += endpoint.CountAllDataTransforms()
	}
	return count
}
