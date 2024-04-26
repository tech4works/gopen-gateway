package vo

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

import (
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
)

// GopenJson represents the configuration json for the Gopen application.
type GopenJson struct {
	// Comment represents a string comment field in the GopenJson struct.
	Comment string `json:"@comment,omitempty"`
	// Version represents the version of the Gopen configuration.
	Version string `json:"version,omitempty"`
	// Port represents the port number on which the Gopen application will listen for incoming requests.
	// It is an integer value and can be specified in the Gopen configuration JSON file.
	Port int `json:"port,omitempty" validate:"required,min=1,max=65535"`
	// HotReload represents a boolean flag indicating whether hot-reloading is enabled or not.
	// It is a field in the Gopen struct and is specified in the Gopen configuration JSON file.
	// It is used to control whether the Gopen application will automatically reload the configuration file
	// and apply the changes and restart the server.
	// If the value is true, hot-reloading is enabled. If the value is false, hot-reloading is disabled.
	// By default, hot-reloading is disabled, so if the field is not specified in the JSON file, it will be set to false.
	HotReload bool `json:"hot-reload,omitempty"`
	// Store represents the store configuration for the Gopen application.
	// It contains the Redis configuration.
	Store *StoreJson `json:"store,omitempty"`
	// Timeout represents the timeout duration for a request or operation.
	// It is specified in string format and can be parsed into a time.Duration value.
	// The default value is empty. If not provided, the timeout will be 30s.
	Timeout Duration `json:"timeout,omitempty"`
	// Cache is a struct representing the cache configuration in the Gopen struct. It contains the following fields:
	// - Duration: a string representing the duration of the cache in a format compatible with Go's time.ParseDuration
	// function. It defaults to an empty string. If not provided, the duration will be 30s.
	// - StrategyHeaders: a slice of strings representing the headers used to determine the cache strategy. It defaults
	// to an empty slice.
	// - OnlyIfStatusCodes: A slice of integers representing the HTTP status codes for which the cache should be used.
	// Default is an empty slice. If not provided, the default value is 2xx success HTTP status codes
	// - OnlyIfMethods: a slice of strings representing the HTTP methods for which the cache should be used. The default
	// is an empty slice. If not provided by default, we only consider the http GET method.
	//- AllowCacheControl: a pointer to a boolean indicating whether the cache should honor the Cache-Control header.
	// It defaults to nil.
	Cache *CacheJson `json:"cache,omitempty"`
	// Limiter represents the configuration for rate limiting.
	// It specifies the maximum header size, maximum body size, maximum multipart memory size, and the rate of allowed requests.
	Limiter *LimiterJson `json:"limiter,omitempty"`
	// SecurityCors represents the configuration options for Cross-Origin Resource Sharing (CORS) settings in Gopen.
	SecurityCors *SecurityCorsJson `json:"security-cors,omitempty"`
	// Middlewares is a map that represents the middleware configuration in Gopen.
	// The keys of the map are the names of the middlewares, and the values are
	// Backend objects that define the properties of each middleware.
	// The Backend struct contains fields like Name, Hosts, Path, Method, ForwardHeaders,
	// ForwardQueries, Modifiers, and ExtraConfig, which specify the behavior
	// and settings of the middleware.
	Middlewares map[string]BackendJson `json:"middlewares,omitempty"`
	// Endpoints is a field in the Gopen struct that represents a slice of Endpoint objects.
	// Each Endpoint object defines a specific API endpoint with its corresponding settings such as path, method,
	// timeout, limiter, cache, etc.
	Endpoints []EndpointJson `json:"endpoints,omitempty"`
}

// StoreJson represents the store configuration json for the Gopen application.
// It contains the Redis configuration.
type StoreJson struct {
	// Redis represents the Redis configuration for the Gopen application.
	Redis *Redis `json:"redis,omitempty"`
}

// Redis represents the configuration for connecting to a Redis server.
// It contains the following fields:
// - Address: a string representing the address of the Redis server. It defaults to an empty string.
// - Password: a string representing the password to authenticate with the Redis server. It defaults to an empty string.
type Redis struct {
	Address  string `json:"address,omitempty" validate:"required,url"`
	Password string `json:"password,omitempty"`
}

// CacheJson represents the cache configuration json in the GopenJson struct.
type CacheJson struct {
	// Duration is a string field in the Cache struct that represents the duration of the cache.
	// It is specified in a format compatible with Go's time.ParseDuration function.
	// The default value is an empty string. If not provided, the duration will be 30s.
	Duration Duration `json:"duration,omitempty" validate:"required,gt=0"`
	// StrategyHeaders represents a slice of strings that contains the headers used to determine the cache strategy key.
	StrategyHeaders []string `json:"strategy-headers,omitempty" validate:"dive,required"`
	// OnlyIfStatusCodes represents a slice of integers representing the HTTP status codes
	// for which the cache should be used. Default is an empty slice. If not provided,
	// the default value is 2xx success HTTP status codes.
	OnlyIfStatusCodes []int `json:"only-if-status-codes,omitempty" validate:"dive,gte=100,lte=599"`
	// OnlyIfMethods represents a slice of strings that contains the HTTP methods for which the cache should be used.
	// The default value is an empty slice. If not provided, the default value is an empty slice.
	// Example: []string{"GET", "POST"}
	OnlyIfMethods []string `json:"only-if-methods,omitempty" validate:"dive,required,http_method"`
	// AllowCacheControl represents a pointer to a boolean indicating whether the cache should
	// honor the Cache-Control header. It defaults to nil. If not provided, the default value is false.
	AllowCacheControl *bool `json:"allow-cache-control,omitempty"`
}

// EndpointCacheJson represents the cache configuration json for an endpoint.
type EndpointCacheJson struct {
	// Enabled represents a boolean indicating whether caching is enabled for an endpoint.
	Enabled bool `json:"enabled"`
	// IgnoreQuery represents a boolean indicating whether to ignore query parameters when caching.
	IgnoreQuery bool `json:"ignore-query,omitempty"`
	// Duration represents the duration configuration for caching an endpoint response.
	Duration Duration `json:"duration,omitempty" validate:"omitempty,gt=0"`
	// StrategyHeaders represents a slice of strings for strategy headers
	StrategyHeaders []string `json:"strategy-headers,omitempty" validate:"dive,required"`
	// OnlyIfStatusCodes represents the status codes that the cache should be applied to.
	OnlyIfStatusCodes []int `json:"only-if-status-codes,omitempty" validate:"dive,gte=100,lte=599"`
	// AllowCacheControl represents a boolean value indicating whether the cache control header is allowed for the endpoint cache.
	AllowCacheControl *bool `json:"allow-cache-control,omitempty"`
}

// LimiterJson represents the configuration for limiter json in the Gopen application.
type LimiterJson struct {
	// MaxHeaderSize represents the maximum size of the header in bytes for rate limiting.
	MaxHeaderSize Bytes `json:"max-header-size,omitempty"`
	// MaxBodySize represents the maximum size of the body in bytes for rate limiting.
	MaxBodySize Bytes `json:"max-body-size,omitempty"`
	// MaxMultipartMemorySize represents the maximum memory size for multipart request bodies.
	MaxMultipartMemorySize Bytes `json:"max-multipart-memory-size,omitempty"`
	// Rate represents the configuration for rate limiting in the Limiter struct. It specifies the capacity and
	// frequency of allowed requests.
	Rate RateJson `json:"rate,omitempty"`
}

// RateJson represents the configuration for rate limiting. It specifies the capacity
// and frequency of allowed requests.
type RateJson struct {
	// Capacity represents the maximum number of allowed requests within a given time period.
	Capacity int `json:"capacity,omitempty" validate:"required,gt=0"`
	// Every represents the frequency of allowed requests in the Rate configuration for rate limiting.
	Every Duration `json:"every,omitempty" validate:"omitempty,gt=0"`
}

// EndpointLimiterJson represents the configuration for endpoint limiter json in the Gopen application.
type EndpointLimiterJson struct {
	// MaxHeaderSize represents the maximum size of the header in bytes for rate limiting.
	MaxHeaderSize Bytes `json:"max-header-size,omitempty"`
	// MaxBodySize represents the maximum size of the body in bytes for rate limiting.
	MaxBodySize Bytes `json:"max-body-size,omitempty"`
	// MaxMultipartMemorySize represents the maximum memory size for multipart request bodies.
	MaxMultipartMemorySize Bytes `json:"max-multipart-memory-size,omitempty"`
	// Rate represents the configuration for rate limiting in the Limiter struct. It specifies the capacity and
	// frequency of allowed requests.
	Rate *RateJson `json:"rate,omitempty"`
}

// SecurityCorsJson represents the configuration options for Cross-Origin Resource Sharing (CORS) settings in Gopen.
type SecurityCorsJson struct {
	// AllowOrigins represents the allowed origins for Cross-Origin Resource Sharing (CORS) settings in Gopen.
	AllowOrigins []string `json:"allow-origins" validate:"dive,required,url"`
	// AllowMethods represents the allowed HTTP methods for Cross-Origin Resource Sharing (CORS) settings in Gopen.
	AllowMethods []string `json:"allow-methods" validate:"dive,required,http_method"`
	// AllowHeaders represents the list of allowed headers for Cross-Origin Resource Sharing (CORS) settings in Gopen.
	AllowHeaders []string `json:"allow-headers" validate:"dive,required"`
}

// EndpointJson represents the configuration for an API endpoint in the Gopen application.
type EndpointJson struct {
	// Comment represents a string that can be used to provide additional information or documentation about an API endpoint.
	Comment string `json:"@comment,omitempty"`
	// Path is a string representing the path of the API endpoint. It is a field in the Endpoint struct.
	Path string `json:"path,omitempty" validate:"required,url_path"`
	// Method represents the HTTP method of an API endpoint.
	Method string `json:"method,omitempty" validate:"required,http_method"`
	// Timeout represents the timeout duration for the API endpoint.
	// It is a string value specified in the JSON configuration.
	// The default value is empty. If not provided, the timeout will be Gopen.Timeout.
	Timeout Duration `json:"timeout,omitempty"`
	// Limiter represents the configuration for rate limiting in the Gopen application.
	// The default value is nil. If not provided, the limiter will be Gopen.Limiter.
	Limiter *EndpointLimiterJson `json:"limiter,omitempty"`
	// Cache represents the cache configuration for an endpoint.
	// The default value is EndpointCache empty with enabled false.
	Cache *EndpointCacheJson `json:"cache,omitempty"`
	// ResponseEncode represents the encoding format for the API endpoint response. The ResponseEncode
	// field is an enum.ResponseEncode value, which can have one of the following values:
	// - enum.ResponseEncodeText: for encoding the response as plain text.
	// - enum.ResponseEncodeJson: for encoding the response as JSON.
	// - enum.ResponseEncodeXml: for encoding the response as XML.
	// The default value is empty. If not provided, the response will be encoded by type, if the string is json it
	// returns json, otherwise it responds to plain text
	ResponseEncode enum.ResponseEncode `json:"response-encode,omitempty" validate:"omitempty,enum"`
	// AggregateResponses represents a boolean indicating whether the API endpoint should aggregate responses
	// from multiple backends.
	AggregateResponses bool `json:"aggregate-responses,omitempty"`
	// AbortIfStatusCodes represents a slice of integers representing the HTTP status codes
	// for which the API endpoint should abort. It is a field in the Endpoint struct.
	AbortIfStatusCodes *[]int `json:"abort-if-status-codes,omitempty" validate:"dive,gte=100,lte=599"`
	// Beforewares represents a slice of strings containing the names of the beforeware middlewares that should be
	// applied before processing the API endpoint.
	Beforewares []string `json:"beforewares,omitempty"`
	// Afterwares represents the configuration for the afterware middlewares to apply after processing the API endpoint.
	// It is a slice of strings representing the names of the afterware middlewares to apply.
	// The names specify the behavior and settings of each afterware middleware.
	// If not provided, the default value is an empty slice.
	// The afterware middleware is executed after processing the API endpoint, allowing for modification or
	// transformation of the response or performing any additional actions.
	// Afterwares can be used for logging, error handling, response modification, etc.
	Afterwares []string `json:"afterwares,omitempty"`
	// Backends represents the backend configurations for an API endpoint in the Gopen application.
	// It is a slice of Backend structs.
	Backends []BackendJson `json:"backends,omitempty" validate:"required,min=1"`
}

// BackendJson represents the configuration for a backend in the Gopen application.
type BackendJson struct {
	// Comment represents a string comment field in the Gopen application configuration.
	Comment string `json:"@comment,omitempty"`
	// Name is a field in the Backend struct that represents the name of the backend configuration.
	Name string `json:"name,omitempty"`
	// Hosts represents a slice of strings that specifies the hosts for a backend configuration.
	Hosts []string `json:"hosts,omitempty" validate:"required,min=1,dive,url"`
	// Path is a field in the Backend struct that represents the path for a backend request.
	// Example: "/api/users"
	Path UrlPath `json:"path,omitempty" validate:"required,url_path"`
	// Method represents the HTTP method for a backend request in the Gopen application.
	// It is a field in the Backend struct and is specified in the Gopen configuration JSON file.
	// The value should be a string and can be one of the following: "GET", "POST", "PUT", "PATCH", "DELETE".
	// It is used to specify the HTTP method for the backend request.
	Method string `json:"method,omitempty" validate:"required,http_method"`
	// ForwardHeaders is a field in the Backend struct that represents a list of headers to be forwarded
	// in the backend request. It is specified as a slice of strings in the Gopen configuration JSON file.
	// Each string represents a header name.
	// Example: ["Content-Type", "User-Agent"]
	ForwardHeaders []string `json:"forward-headers,omitempty" validate:"dive,required"`
	// ForwardQueries represents the list of query parameters that will be forwarded to the backend server.
	// The query parameters are specified as string elements in a slice.
	// It is a field in the Backend struct and is specified in the Gopen configuration JSON file.
	// The ForwardQueries field is used to specify which query parameters of the incoming request will be included in the
	// request sent to the backend server.
	ForwardQueries []string `json:"forward-queries,omitempty" validate:"dive,required"`
	// Modifiers represent the configuration to modify the request and response of a backend and endpoint in the Gopen application.
	Modifiers *BackendModifiersJson `json:"modifiers,omitempty"`
	// ExtraConfig represents additional configuration options for a backend in the Gopen application.
	ExtraConfig *BackendExtraConfigJson `json:"extra-config,omitempty"`
}

// BackendModifiersJson represents a set of modifiers that can be applied to different parts of the request and response
// in the Gopen application.
type BackendModifiersJson struct {
	// Comment represents a string comment field in the Gopen application configuration.
	Comment string `json:"@comment,omitempty"`
	// StatusCode represents the status code that can be applied to a request or response. It is an integer value and is
	// specified in the Gopen configuration JSON file. The status code is used to indicate the status of the response,
	// such as success, failure, or error. The status code is optional and can be omitted from the configuration.
	StatusCode int `json:"status-code,omitempty" validate:"omitempty,gte=100,lte=599"`
	// Header represents a slice of modifying structures that can be applied to the header of a request or response from
	// the Endpoint or just the current Backend.
	Header []ModifierJson `json:"header,omitempty" validate:"dive,required"`
	// Param is a slice of modifiers that can be applied to the parameters of a request from the Endpoint
	// or just the current Backend.
	Param []ModifierJson `json:"param,omitempty" validate:"dive,required"`
	// Query represents a slice of Modifier structs that can be applied to the query parameters of a request
	// from the Endpoint or just the current Backend.
	Query []ModifierJson `json:"query,omitempty" validate:"dive,required"`
	// Body represents a slice of Modifier structs that can be applied to the body of a request or response
	// from the Endpoint or just the current Backend.
	Body []ModifierJson `json:"body,omitempty" validate:"dive,required"`
}

// BackendExtraConfigJson represents additional configuration options for a backend in the Gopen application.
// - OmitRequestBody: a boolean flag indicating whether the backend should omit the request body in the outgoing request.
// If set to true, the backend will not include the request body in the outgoing request.
// If set to false, the request body will be included in the outgoing request.
// The default value is false.
// - OmitResponse: a boolean flag indicating whether the backend should omit the response in the incoming request.
// If set to true, the backend will not include the response in the incoming request.
// If set to false, the response will be included in the incoming request.
// The default value is false.
type BackendExtraConfigJson struct {
	// GroupResponse is a boolean flag indicating whether the backend should group response.
	// The default value is false.
	GroupResponse bool `json:"group-response"`
	// OmitRequestBody represents a boolean flag indicating whether the backend should omit the request body in request.
	// If set to true, the backend will not include the request body in the request.
	// If set to false, the request body will be included in the request. The default value is false.
	OmitRequestBody bool `json:"omit-request-body"`
	// OmitResponse represents a boolean flag indicating whether the backend should omit the response in the incoming request.
	// If set to true, the backend will not include the response in the incoming request.
	// If set to false, the response will be included in the incoming request.
	// The default value is false.
	OmitResponse bool `json:"omit-response"`
}

// ModifierJson represents a modification that can be applied to a request or response in the Gopen application.
type ModifierJson struct {
	// Comment represents a string comment field in the Gopen application configuration.
	Comment string `json:"@comment,omitempty"`
	// Context represents the context in which a modification should be applied.
	// It is an enum.ModifierContext value.
	// Valid values for Context are "request" and "response".
	Context enum.ModifierContext `json:"context,omitempty" validate:"required,enum"`
	// Scope represents the scope of a modification in the Backend or Endpoint.
	// It is an enum.ModifierScope value that specifies where the modification should be applied.
	// Valid values for Scope are "request" and "response".
	Scope enum.ModifierScope `json:"scope,omitempty" validate:"omitempty,enum"`
	// Action represents the action to be performed in the Modifier struct.
	// It is an enum.ModifierAction value and can be one of the following values:
	// - ModifierActionSet: to set a value.
	// - ModifierActionAdd: to add a value.
	// - ModifierActionDel: to delete a value.
	// - ModifierActionReplace: to replace a value.
	// - ModifierActionRename: to rename a value.
	Action enum.ModifierAction `json:"action,omitempty" validate:"required,enum"`
	// Propagate represents a Boolean flag that indicates whether the modification should be propagated to subsequent
	// Backend requests.
	Propagate bool `json:"propagate,omitempty"`
	// Key represents a string value that serves as the key for a modification in the Modifier structure.
	// Indicates the field that you want to modify.
	Key string `json:"key,omitempty"`
	// Value represents a string value in the Modifier struct.
	// It is used as a field to store the value of a modification.
	Value string `json:"value,omitempty"`
}

// NewGopenJson parses the given JSON file bytes and returns a GopenJson object.
// It converts the fileJsonBytes into a GopenJson value object using the ConvertToDest function.
// Then, it validates the GopenJson object using the Validate function from the helper package.
// If any error occurs during validation, it returns an error.
// Finally, it returns the GopenJson object and nil as the error if everything is successful.
func NewGopenJson(fileJsonBytes []byte) (*GopenJson, error) {
	// convertemos o valor em bytes em VO
	var gopenJsonVO GopenJson
	err := helper.ConvertToDest(fileJsonBytes, &gopenJsonVO)
	if helper.IsNil(err) {
		err = helper.Validate().Struct(gopenJsonVO)
	}

	// se ocorreu algum erro retornamos
	if helper.IsNotNil(err) {
		return nil, errors.New("Error parse Gopen json file to VO:", err)
	}

	// retornamos o DTO que é a configuração do Gopen
	return &gopenJsonVO, nil
}

// Json returns a new GopenJson object with the same values as the receiver.
// It initializes and returns a new GopenJson object by assigning the values of the receiver
// to the corresponding fields of the new object.
// This method can be used to create a copy of the GopenJson object.
// The returned GopenJson object shares the same values as the receiver,
// but they are separate objects in memory.
func (g GopenJson) Json() *GopenJson {
	return &GopenJson{
		Version:      g.Version,
		Port:         g.Port,
		HotReload:    g.HotReload,
		Timeout:      g.Timeout,
		Cache:        g.Cache,
		Limiter:      g.Limiter,
		SecurityCors: g.SecurityCors,
		Middlewares:  g.Middlewares,
		Endpoints:    g.Endpoints,
	}
}

// CountMiddlewares returns the number of middlewares in the Gopen instance.
func (g GopenJson) CountMiddlewares() int {
	return len(g.Middlewares)
}

// CountEndpoints returns the number of endpoints in the Gopen struct.
func (g GopenJson) CountEndpoints() int {
	return len(g.Middlewares)
}

// CountBackends returns the total number of backends present in the `Gopen` struct and its nested `Endpoint` structs.
// It calculates the count by summing the number of middlewares in `Gopen` and recursively iterating through each `Endpoint`
// to count their backends.
// Returns an integer indicating the total count of backends.
func (g GopenJson) CountBackends() (count int) {
	count += g.CountMiddlewares()
	for _, endpointVO := range g.Endpoints {
		count += endpointVO.CountBackends()
	}
	return count
}

// CountModifiers counts the total number of modifiers in the Gopen struct.
// It iterates through all the middleware backends and endpoint VOs,
// and calls the CountModifiers method on each of them to calculate the count.
// The count is incremented for each modifier found and the final count is returned.
func (g GopenJson) CountModifiers() (count int) {
	for _, middlewareBackend := range g.Middlewares {
		count += middlewareBackend.CountModifiers()
	}
	for _, endpointDTO := range g.Endpoints {
		count += endpointDTO.CountModifiers()
	}
	return count
}

// CountBackends returns the number of backends in the Endpoint struct.
func (e EndpointJson) CountBackends() int {
	if helper.IsNil(e.Backends) {
		return 0
	}
	return len(e.Backends)
}

// CountModifiers counts the total number of modifiers in an Endpoint by summing the count of modifiers in each
// Backend associated with it.
func (e EndpointJson) CountModifiers() (count int) {
	for _, backendJsonVO := range e.Backends {
		count += backendJsonVO.CountModifiers()
	}
	return count
}

// HasMaxHeaderSize checks if the MaxHeaderSize field in the EndpointLimiterJson object is greater than 0.
// If it is, it returns true. Otherwise, it returns false.
func (e EndpointLimiterJson) HasMaxHeaderSize() bool {
	return helper.IsGreaterThan(e.MaxHeaderSize, 0)
}

// HasMaxBodySize checks if the MaxBodySize field in the EndpointLimiterJson object is greater than 0.
// If it is, it returns true. Otherwise, it returns false.
func (e EndpointLimiterJson) HasMaxBodySize() bool {
	return helper.IsGreaterThan(e.MaxBodySize, 0)
}

// HasMaxMultipartMemorySize checks if the MaxMultipartMemorySize field in the EndpointLimiterJson object is greater than 0.
// If it is, it returns true. Otherwise, it returns false.
func (e EndpointLimiterJson) HasMaxMultipartMemorySize() bool {
	return helper.IsGreaterThan(e.MaxMultipartMemorySize, 0)
}

// HasEvery checks if the Every field in the RateJson object is greater than 0.
// It returns true if the Every field is greater than 0, and false otherwise.
// This method is typically used to determine if the rate limiting configuration
// for a specific endpoint includes a valid value for the frequency of allowed requests.
func (r RateJson) HasEvery() bool {
	return helper.IsGreaterThan(r.Every, 0)
}

// HasCapacity checks if the Capacity field in the RateJson object is greater than 0.
// It returns true if the Capacity field is greater than 0, and false otherwise.
// This method is typically used to determine if the rate limiting configuration
// for a specific endpoint includes a valid value for the maximum number of allowed requests.
func (r RateJson) HasCapacity() bool {
	return helper.IsGreaterThan(r.Capacity, 0)
}

// CountModifiers returns the number of modifiers present in the Backend instance.
// If the modifiers field is not nil, it counts all the modifiers using the CountAll() method of BackendModifiers.
// Otherwise, it returns 0.
func (b BackendJson) CountModifiers() int {
	if helper.IsNotNil(b.Modifiers) {
		return b.Modifiers.CountAll()
	}
	return 0
}

// CountAll returns the total count of modifiers for a BackendModifiers instance.
// It counts the number of valid `statusCode` and the length of `header`, `params`, `query`, and `body` slices,
// and adds them up to get the total count.
func (b BackendModifiersJson) CountAll() (count int) {
	if helper.IsNotEmpty(b.StatusCode) {
		count++
	}
	count += len(b.Header) + len(b.Param) + len(b.Query) + len(b.Body)
	return count
}

// HasDuration returns a boolean value indicating whether the duration of the cache in the EndpointCacheJson object is
// greater than zero.
func (e EndpointCacheJson) HasDuration() bool {
	return helper.IsGreaterThan(e.Duration, 0)
}

// HasStrategyHeaders returns a boolean value indicating whether the `strategyHeaders` field in the EndpointCacheJson.
// struct is not nil.
func (e EndpointCacheJson) HasStrategyHeaders() bool {
	return helper.IsNotNil(e.StrategyHeaders)
}

// HasOnlyIfStatusCodes returns a boolean value indicating whether the `onlyIfStatusCodes` field in the EndpointCache struct
// is not nil. If the field is not nil, it means that the cache should only be applied to the specified status codes, and the
// function returns true. Otherwise, it returns false.
func (e EndpointCacheJson) HasOnlyIfStatusCodes() bool {
	return helper.IsNotNil(e.OnlyIfStatusCodes)
}

// HasAllowCacheControl returns a boolean value indicating whether the `allowCacheControl` field in the EndpointCacheJson struct
// is not nil. If the field is not nil, it means that the cache control header is allowed for the endpoint cache, and the
// function returns true. Otherwise, it returns false.
func (e EndpointCacheJson) HasAllowCacheControl() bool {
	return helper.IsNotNil(e.AllowCacheControl)
}
