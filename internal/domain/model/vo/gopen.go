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
	"time"
)

// Gopen is a struct that represents the configuration for a Gopen server.
type Gopen struct {
	// version is a string field that represents the version of the Gopen server configured in the configuration json.
	version string
	// port represents the port number on which the Gopen application will listen for incoming requests.
	// It is an integer value and can be specified in the Gopen configuration JSON file.
	port int
	// hotReload represents a boolean flag indicating whether hot-reloading is enabled or not.
	// It is a field in the Gopen struct and is specified in the Gopen configuration JSON file.
	// It is used to control whether the Gopen application will automatically reload the configuration file
	// and apply the changes and restart the server.
	// If the value is true, hot-reloading is enabled. If the value is false, hot-reloading is disabled.
	// By default, hot-reloading is disabled, so if the field is not specified in the JSON file, it will be set to false.
	hotReload bool
	// timeout represents the timeout duration for a httpRequest or operation.
	// It is specified in string format and can be parsed into a time.Duration value.
	// The default value is empty. If not provided, the timeout will be 30s.
	timeout Duration
	// limiter represents the configuration for rate limiting.
	// It specifies the maximum header size, maximum body size, maximum multipart memory size, and the rate of allowed requests.
	limiter Limiter
	// cache is a struct representing the `cache` configuration in the Gopen struct. It contains the following fields:
	// - Duration: a string representing the duration of the `cache` in a format compatible with Go's time.ParseDuration
	// function. It defaults to an empty string. If not provided, the duration will be 30s.
	// - StrategyHeaders: a slice of strings representing the modifyHeaders used to determine the `cache` strategy. It defaults
	// to an empty slice.
	// - OnlyIfStatusCodes: A slice of integers representing the HTTP status codes for which the `cache` should be used.
	// Default is an empty slice. If not provided, the default value is 2xx success HTTP status codes
	// - OnlyIfMethods: a slice of strings representing the HTTP methods for which the `cache` should be used. The default
	// is an empty slice. If not provided by default, we only consider the http GET method.
	//- AllowCacheControl: a pointer to a boolean indicating whether the `cache` should honor the Cache-Control header.
	// It defaults to empty.
	cache *Cache
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
func NewGopen(gopenJsonVO *GopenJson) *Gopen {
	// montamos o VO
	gopenVO := &Gopen{
		version:      gopenJsonVO.Version,
		port:         gopenJsonVO.Port,
		hotReload:    gopenJsonVO.HotReload,
		timeout:      gopenJsonVO.Timeout,
		limiter:      newLimiter(gopenJsonVO.Limiter),
		cache:        newCache(gopenJsonVO.Cache),
		securityCors: newSecurityCors(gopenJsonVO.SecurityCors),
		middlewares:  newMiddlewares(gopenJsonVO.Middlewares),
	}

	// damos o parse dos endpoints de VO json para o VO
	var endpoints []Endpoint
	for _, endpointJsonVO := range gopenJsonVO.Endpoints {
		endpoints = append(endpoints, newEndpoint(gopenVO, &endpointJsonVO))
	}
	gopenVO.endpoints = endpoints

	// retornamos o novo objeto de valor
	return gopenVO
}

// Port returns the value of the port field in the Gopen struct.
func (g Gopen) Port() int {
	return g.port
}

// Timeout returns the value of the timeout field in the Gopen struct. If the timeout is greater than 0,
// it returns the timeout value. Otherwise, it returns a default timeout of 30 seconds
func (g Gopen) Timeout() Duration {
	if helper.IsGreaterThan(g.timeout, 0) {
		return g.timeout
	}
	return Duration(30 * time.Second)
}

// Cache returns the value of the cache field in the Gopen struct.
func (g Gopen) Cache() *Cache {
	return g.cache
}

// Limiter returns the value of the limiter field in the Gopen struct.
func (g Gopen) Limiter() Limiter {
	return g.limiter
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
