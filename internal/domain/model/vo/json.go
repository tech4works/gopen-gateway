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
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"strings"
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
	// Timeout represents the timeout duration for a httpRequest or operation.
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
	// Address is a string field representing the address of the Redis server.
	Address string `json:"address,omitempty" validate:"required,url"`
	// Password represents the password to authenticate with a Redis server.
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
	// MaxMultipartMemorySize represents the maximum memory size for multipart httpRequest bodies.
	MaxMultipartMemorySize Bytes `json:"max-multipart-memory-size,omitempty"`
	// Rate represents the configuration for rate limiting in the Limiter struct. It specifies the capacity and
	// frequency of allowed requests.
	Rate *RateJson `json:"rate,omitempty"`
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
	// MaxMultipartMemorySize represents the maximum memory size for multipart httpRequest bodies.
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

// MiddlewaresJson represents a map of string keys to BackendJson values. It is used to configure and store middleware
// settings in the Gopen application.
// Each key in the map represents the name of the middleware, and the corresponding BackendJson value defines the
// properties and behavior of that middleware.
type MiddlewaresJson map[string]BackendJson

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
	// The default value is Cache empty with enabled false.
	Cache *EndpointCacheJson `json:"cache,omitempty"`
	// AbortIfStatusCodes represents a slice of integers representing the HTTP status codes
	// for which the API endpoint should abort. It is a field in the Endpoint struct.
	AbortIfStatusCodes *[]int `json:"abort-if-status-codes,omitempty" validate:"dive,gte=100,lte=599"`
	// Response is the field in the `EndpointJson` struct that represents the configuration for an API endpoint response.
	// It is of type `EndpointResponseJson` and is used to define how the response should be encoded and if the responses
	// should be aggregated from multiple backends.
	Response *EndpointResponseJson `json:"response,omitempty"`
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

// EndpointResponseJson represents the configuration for an API endpoint response.
type EndpointResponseJson struct {
	// Aggregate represents a boolean indicating whether the API endpoint should aggregate responses
	// from multiple backends.
	Aggregate bool `json:"aggregate,omitempty"`
	// Encode represents the encoding format for the API endpoint response. The ResponseEncode
	// field is an enum.Encode value, which can have one of the following values:
	// - enum.EncodeText: for encoding the response as plain text.
	// - enum.EncodeJson: for encoding the response as JSON.
	// - enum.EncodeXml: for encoding the response as XML.
	// The default value is empty. If not provided, the response will be encoded by type, if the string is json it
	// returns json, otherwise it responds to plain text
	Encode enum.Encode `json:"encode,omitempty" validate:"omitempty,enum"`
	// Nomenclature field represents the case format for json text fields.
	Nomenclature enum.Nomenclature `json:"nomenclature,omitempty" validate:"omitempty,enum"`
	// OmitEmpty represents a boolean value indicating whether the field should be omitted
	OmitEmpty bool `json:"omit-empty,omitempty"`
}

// BackendJson represents the configuration for a backend in the Gopen application.
type BackendJson struct {
	// Comment represents a string comment field in the Gopen application configuration.
	Comment string `json:"@comment,omitempty"`
	// Hosts represents a slice of strings that specifies the hosts for a backend configuration.
	Hosts []string `json:"hosts,omitempty" validate:"required,min=1,dive,url"`
	// Path is a field in the Backend struct that represents the path for a backend httpRequest.
	// Example: "/api/users"
	Path string `json:"path,omitempty" validate:"required,url_path"`
	// Method represents the HTTP method for a backend httpRequest in the Gopen application.
	// It is a field in the Backend struct and is specified in the Gopen configuration JSON file.
	// The value should be a string and can be one of the following: "GET", "POST", "PUT", "PATCH", "DELETE".
	// It is used to specify the HTTP method for the backend httpRequest.
	Method string `json:"method,omitempty" validate:"required,http_method"`
	// Request is a field in the BackendJson struct that represents the configuration for a backend's request.
	// It contains various properties and settings related to the request, such as omitting headers, queries, or body,
	// mapper and projection configurations for headers, queries, and body, and modifier configurations for headers,
	// params, queries, and body.
	Request *BackendRequestJson `json:"request,omitempty"`
	// Response is a field in the BackendJson struct that represents the response configuration for a backend in
	// the Gopen application.
	Response *BackendResponseJson `json:"response,omitempty"`
}

// BackendRequestJson represents the JSON structure for a backend request in the Gopen application configuration.
type BackendRequestJson struct {
	// OmitHeader is a boolean flag that indicates whether the header should be omitted in the backend HTTP request.
	OmitHeader bool `json:"omit-header,omitempty"`
	// OmitQuery is a boolean flag that indicates whether the query should be omitted in the backend HTTP request.
	OmitQuery bool `json:"omit-query,omitempty"`
	// OmitBody is a boolean flag that indicates whether the body should be omitted in the backend HTTP request.
	OmitBody bool `json:"omit-body,omitempty"`
	// HeaderMapper represents the mapping of header fields in a backend request.
	// It is used to rename keys in the header of a request.
	HeaderMapper *Mapper `json:"header-mapper,omitempty"`
	// QueryMapper represents a mapper for mapping keys to values in the query parameters of a backend request.
	// It is used to rename keys in the query of a request.
	QueryMapper *Mapper `json:"query-mapper,omitempty"`
	// BodyMapper is a field in the BackendRequestJson structure that represents a mapper for the body of a backend request.
	// It is used to rename keys in the body of a request.
	BodyMapper *Mapper `json:"body-mapper,omitempty"`
	// HeaderProjection represents a projection of headers in the BackendRequestJson structure.
	HeaderProjection *Projection `json:"header-projection,omitempty"`
	// QueryProjection represents a projection of keys and values that can be applied to the query parameters
	// of a backend request in the Gopen application configuration.
	QueryProjection *Projection `json:"query-projection,omitempty"`
	// BodyProjection is a struct that represents a projection for the body of a BackendRequestJson.
	BodyProjection *Projection `json:"body-projection,omitempty"`
	// HeaderModifiers represents a slice of ModifierJson objects that can be applied to the header of a backend request
	// in the Gopen application configuration.
	HeaderModifiers []ModifierJson `json:"header-modifiers,omitempty" validate:"dive,required"`
	// ParamModifiers represents a list of modifications that can be applied to request parameters
	ParamModifiers []ModifierJson `json:"param-modifiers,omitempty" validate:"dive,required"`
	// QueryModifiers is a slice of ModifierJson structs representing modifications that can be applied to an HTTP
	// request or response.
	QueryModifiers []ModifierJson `json:"query-modifiers,omitempty" validate:"dive,required"`
	// BodyModifiers represents an array of ModifierJson objects. It is a field
	// in the BackendRequestJson struct and is used to specify modifications that
	// can be applied to the body of an HTTP request.
	BodyModifiers []ModifierJson `json:"body-modifiers,omitempty" validate:"dive,required"`
}

// BackendResponseJson represents the JSON response from a backend.
type BackendResponseJson struct {
	// Apply represents the application scope of a BackendResponse. It is an optional field that indicates whether
	// these settings must be applied sooner or later. Possible values for Apply are "EARLY" and "LATE".
	Apply enum.BackendResponseApply `json:"apply,omitempty" validate:"omitempty,enum"`
	// Omit is a field that indicates whether the HTTP response from the backend should be omitted or not.
	// If Omit is true, the HTTP response from the backend will be omitted to the final HTTP client. Otherwise,
	// it will be returned to the end customer.
	Omit bool `json:"omit,omitempty"`
	// OmitHeader is a boolean field that indicates whether the backend HTTP response header should be omitted or not.
	// If OmitHeader is true, the header property will be omitted for the final HTTP client. Otherwise, the header will
	// be displayed to the end customer.
	OmitHeader bool `json:"omit-header,omitempty"`
	// OmitBody is a boolean field that indicates whether the body of the backend response should be omitted or not.
	// If OmitBody is true, the body property will be omitted for the final HTTP client. Otherwise, it will be displayed
	// if it has a body.
	OmitBody bool `json:"omit-body,omitempty"`
	// Group is a string field, which represents the field to be used to group the backend HTTP response body, to the
	// final HTTP response body. It is optional and must have a minimum length of 1 character.
	Group string `json:"group,omitempty" validate:"omitempty,min=1"`
	// HeaderMapper represents a mapper for headers in a BackendResponseJson.
	// Allows you to rename header fields using the key to identify the current field, and the value representing the
	// new value of the header key.
	HeaderMapper *Mapper `json:"header-mapper,omitempty"`
	// BodyMapper represents a mapper for the body in a BackendResponseJson.
	// Allows you to rename JSON fields using the key to identify the current field and the value that represents the
	// new JSON key value.
	BodyMapper *Mapper `json:"body-mapper,omitempty"`
	// HeaderProjection represents the header projection of a BackendResponseJson.
	// It is used to define which header fields should be included or excluded when processing the response.
	HeaderProjection *Projection `json:"header-projection,omitempty"`
	// BodyProjection represents the body projection of a BackendResponseJson.
	// It is used to define which body fields should be included or excluded when processing the response.
	BodyProjection *Projection `json:"body-projection,omitempty"`
	// HeaderModifiers is a slice of ModifierJson objects representing modifications that can be applied to backend
	// response HTTP headers.
	HeaderModifiers []ModifierJson `json:"header-modifiers,omitempty" validate:"dive,required"`
	// BodyModifiers is a slice of ModifierJson representing modifications that can be applied to backend
	// response HTTP body.
	BodyModifiers []ModifierJson `json:"body-modifiers,omitempty" validate:"dive,required"`
}

// ModifierJson represents a modification that can be applied to a httpRequest or response in the Gopen application.
type ModifierJson struct {
	// Comment represents a string comment field in the Gopen application configuration.
	Comment string `json:"@comment,omitempty"`
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

// NewGopenJson parses the given bytes as a GopenJson object and returns a pointer to it.
// It first converts the bytes into a GopenJson object using the ConvertToDest function.
// Then it validates the GopenJson object using the Validate function.
// If any validation errors occur, it returns an error with a descriptive message.
// Next, it calls the ValidateEndpoints function to validate the endpoints in the GopenJson object.
// If any errors occur during endpoint validation, it returns an error with a descriptive message.
// Finally, it returns the pointer to the parsed and validated GopenJson object and nil error.
// If any errors occur during the parsing or validation process, it returns nil and an error message.
func NewGopenJson(bytes []byte) (*GopenJson, error) {
	var gopenJson GopenJson
	err := helper.ConvertToDest(bytes, &gopenJson)
	if helper.IsNil(err) {
		err = helper.Validate().Struct(gopenJson)
	}

	if helper.IsNotNil(err) {
		return nil, errors.New("Error parse Gopen json file to value object:", err)
	} else if err = ValidateEndpoints(gopenJson.Endpoints); helper.IsNotNil(err) {
		return nil, err
	}

	return &gopenJson, nil
}

// ValidateEndpoints validates the given endpoints for any duplicate registrations.
// It compares each endpoint with every other endpoint to check if the path and method are already registered.
// If any duplicate registrations are found, it adds a warning message to the warnings slice.
// If no warnings are found, it returns nil.
// Otherwise, it returns an error with all the warning messages joined by a newline character.
func ValidateEndpoints(endpoints []EndpointJson) error {
	var warnings []string
	for index, endpoint := range endpoints {
		for anotherIndex, anotherEndpoint := range endpoints {
			if helper.IsNotEqualTo(index, anotherIndex) && endpoint.Equals(anotherEndpoint) {
				format := "endpoint path: %s method: %s on index: %v already registered on index %v!"
				warnings = append(warnings, fmt.Sprintf(format, endpoint.Path, endpoint.Method, anotherIndex, index))
			}
		}
	}
	if helper.IsEmpty(warnings) {
		return nil
	}
	return errors.New(fmt.Sprintf("Error %s", strings.Join(warnings, "\n ")))
}

// Filter returns a new GopenJson object with the same values as the receiver.
// It initializes and returns a new GopenJson object by assigning the values of the receiver
// to the corresponding fields of the new object.
// This method can be used to create a copy of the GopenJson object.
// The returned GopenJson object shares the same values as the receiver,
// but they are separate objects in memory.
func (g GopenJson) Filter() *GopenJson {
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

// Equals checks if the current EndpointJson object is equal to another EndpointJson object.
// Two EndpointJson objects are considered equal if their `Path` and `Method` fields are equal.
// If both fields are equal, the method returns true. Otherwise, it returns false.
func (e EndpointJson) Equals(another EndpointJson) bool {
	return helper.Equals(e.Path, another.Path) && helper.Equals(e.Method, another.Method)
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

// HasOnlyIfStatusCodes returns a boolean value indicating whether the `onlyIfStatusCodes` field in the Cache struct
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
