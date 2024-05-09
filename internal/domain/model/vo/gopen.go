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

// NewGopen initializes a new Gopen struct based on the provided GopenJson object.
// It populates the fields of the Gopen struct with the corresponding values from the GopenJson object.
// It also populates the endpoints field by iterating over the EndpointJson objects in the Endpoints slice of the GopenJson object,
// and converting each EndpointJson object to an Endpoint object using the newEndpoint function.
// The newly created Gopen struct is returned as a pointer.
func NewGopen(gopenJson *GopenJson) *Gopen {
	// damos o parse dos endpoints de json para o VO
	var endpoints []Endpoint
	for _, endpointJson := range gopenJson.Endpoints {
		endpoints = append(endpoints, newEndpoint(gopenJson, &endpointJson))
	}

	// montamos o gopen VO
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

// Endpoints returns a slice containing all the endpoints configured in the Gopen struct.
func (g Gopen) Endpoints() []Endpoint {
	return g.endpoints
}

func (g Gopen) CountEndpoints() int {
	return len(g.endpoints)
}

func (g Gopen) CountMiddlewares() int {
	if helper.IsNotNil(g.middlewares) {
		return len(g.middlewares)
	}
	return 0
}

func (g Gopen) CountBackends() (count int) {
	for _, endpoint := range g.Endpoints() {
		count += endpoint.CountBeforewares()
		count += endpoint.CountBackends()
		count += endpoint.CountAfterwares()
	}
	return count
}

func (g Gopen) CountAllDataTransforms() (count int) {
	for _, endpoint := range g.Endpoints() {
		count += endpoint.CountAllDataTransforms()
	}
	return count
}
