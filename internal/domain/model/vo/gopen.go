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
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"time"
)

// Gopen is a struct that represents the configuration for a Gopen server.
type Gopen struct {
	// env is a string field that represents the environment in which the Gopen server is running.
	env string
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
	// timeout represents the timeout duration for a request or operation.
	// It is specified in string format and can be parsed into a time.Duration value.
	// The default value is empty. If not provided, the timeout will be 30s.
	timeout time.Duration
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

// NewGopen creates a new instance of Gopen based on the provided environment and gopenDTO.
// It initializes the fields of Gopen based on values from gopenDTO and sets default values for empty fields.
func NewGopen(env string, gopenDTO *dto.Gopen) *Gopen {
	// damos o parse dos endpoints
	var endpoints []Endpoint
	for _, endpointDTO := range gopenDTO.Endpoints {
		endpoints = append(endpoints, newEndpoint(endpointDTO))
	}

	// damos o parse do timeout
	var timeout time.Duration
	var err error
	if helper.IsNotEmpty(gopenDTO.Timeout) {
		timeout, err = time.ParseDuration(gopenDTO.Timeout)
		if helper.IsNotNil(err) {
			logger.Warning("Parse duration timeout err:", err)
		}
	}

	return &Gopen{
		env:          env,
		version:      gopenDTO.Version,
		port:         gopenDTO.Port,
		hotReload:    gopenDTO.HotReload,
		timeout:      timeout,
		limiter:      newLimiterFromDTO(gopenDTO.Limiter),
		cache:        newCacheFromDTO(gopenDTO.Cache),
		securityCors: newSecurityCors(gopenDTO.SecurityCors),
		middlewares:  newMiddlewares(gopenDTO.Middlewares),
		endpoints:    endpoints,
	}
}

// Port returns the value of the port field in the Gopen struct.
func (g Gopen) Port() int {
	return g.port
}

// HotReload returns the value of the hotReload field in the Gopen struct.
func (g Gopen) HotReload() bool {
	return g.hotReload
}

// Version returns the value of the version field in the Gopen struct.
func (g Gopen) Version() string {
	return g.version
}

// Timeout returns the value of the timeout field in the Gopen struct. If the timeout is greater than 0,
// it returns the timeout value. Otherwise, it returns a default timeout of 30 seconds
func (g Gopen) Timeout() time.Duration {
	if helper.IsGreaterThan(g.timeout, 0) {
		return g.timeout
	}
	return 30 * time.Second
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
func (g Gopen) Middleware(key string) (Backend, bool) {
	return g.middlewares.Get(key)
}

// Middlewares returns the value of the middlewares field in the Gopen struct.
func (g Gopen) Middlewares() Middlewares {
	return g.middlewares
}

// Endpoints returns a slice containing all the endpoints configured in the Gopen struct.
// It iterates over each EndpointVO in the endpoints slice and fills in default values by calling the
// fillDefaultValues method on each EndpointVO, passing the Gopen instance as a parameter.
// The resulting Endpoint slice is returned.
func (g Gopen) Endpoints() []Endpoint {
	endpoints := make([]Endpoint, len(g.endpoints))
	for i, endpointVO := range g.endpoints {
		endpoints[i] = endpointVO.fillDefaultValues(&g)
	}
	return endpoints
}

// PureEndpoints returns a slice containing all the endpoints configured in the Gopen struct.
// No default values are filled in for each EndpointVO, unlike the Endpoints method.
// The resulting Endpoint slice is returned.
func (g Gopen) PureEndpoints() []Endpoint {
	return g.endpoints
}

// CountMiddlewares returns the number of middlewares in the Gopen instance.
func (g Gopen) CountMiddlewares() int {
	return len(g.middlewares)
}

// CountEndpoints returns the number of endpoints in the Gopen struct.
func (g Gopen) CountEndpoints() int {
	return len(g.endpoints)
}

// CountBackends returns the total number of backends present in the `Gopen` struct and its nested `Endpoint` structs.
// It calculates the count by summing the number of middlewares in `Gopen` and recursively iterating through each `Endpoint`
// to count their backends.
// Returns an integer indicating the total count of backends.
func (g Gopen) CountBackends() (count int) {
	count += g.CountMiddlewares()
	for _, endpointVO := range g.endpoints {
		count += endpointVO.CountBackends()
	}
	return count
}

// CountModifiers counts the total number of modifiers in the Gopen struct.
// It iterates through all the middleware backends and endpoint VOs,
// and calls the CountModifiers method on each of them to calculate the count.
// The count is incremented for each modifier found and the final count is returned.
func (g Gopen) CountModifiers() (count int) {
	for _, middlewareBackend := range g.middlewares {
		count += middlewareBackend.CountModifiers()
	}
	for _, endpointDTO := range g.endpoints {
		count += endpointDTO.CountModifiers()
	}
	return count
}
