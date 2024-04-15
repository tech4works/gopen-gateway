package dto

import (
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
)

// Gopen represents the configuration json for the Gopen application.
type Gopen struct {
	// Version represents the version of the Gopen configuration.
	Version string `json:"version,omitempty"`
	// Port represents the port number on which the Gopen application will listen for incoming requests.
	// It is an integer value and can be specified in the Gopen configuration JSON file.
	Port int `json:"port,omitempty"`
	// HotReload represents a boolean flag indicating whether hot-reloading is enabled or not.
	// It is a field in the Gopen struct and is specified in the Gopen configuration JSON file.
	// It is used to control whether the Gopen application will automatically reload the configuration file
	// and apply the changes and restart the server.
	// If the value is true, hot-reloading is enabled. If the value is false, hot-reloading is disabled.
	// By default, hot-reloading is disabled, so if the field is not specified in the JSON file, it will be set to false.
	HotReload bool `json:"hot-reload,omitempty"`
	// Store represents the store configuration for the Gopen application.
	// It contains the Redis configuration.
	Store Store `json:"store,omitempty"`
	// Timeout represents the timeout duration for a request or operation.
	// It is specified in string format and can be parsed into a time.Duration value.
	// The default value is empty. If not provided, the timeout will be 30s.
	Timeout string `json:"timeout,omitempty"`
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
	Cache Cache `json:"cache,omitempty"`
	// Limiter represents the configuration for rate limiting.
	// It specifies the maximum header size, maximum body size, maximum multipart memory size, and the rate of allowed requests.
	Limiter Limiter `json:"limiter,omitempty"`
	// SecurityCors represents the configuration options for Cross-Origin Resource Sharing (CORS) settings in Gopen.
	SecurityCors SecurityCors `json:"security-cors,omitempty"`
	// Middlewares is a map that represents the middleware configuration in Gopen.
	// The keys of the map are the names of the middlewares, and the values are
	// Backend objects that define the properties of each middleware.
	// The Backend struct contains fields like Name, Hosts, Path, Method, ForwardHeaders,
	// ForwardQueries, Modifiers, and ExtraConfig, which specify the behavior
	// and settings of the middleware.
	Middlewares map[string]Backend `json:"middlewares,omitempty"`
	// Endpoints is a field in the Gopen struct that represents a slice of Endpoint objects.
	// Each Endpoint object defines a specific API endpoint with its corresponding settings such as path, method,
	// timeout, limiter, cache, etc.
	Endpoints []Endpoint `json:"endpoints,omitempty"`
}

// Store represents the store configuration for the Gopen application.
// It contains the Redis configuration.
type Store struct {
	// Redis represents the Redis configuration for the Gopen application.
	Redis Redis `json:"redis,omitempty"`
}

// Redis represents the configuration for connecting to a Redis server.
// It contains the following fields:
// - Address: a string representing the address of the Redis server. It defaults to an empty string.
// - Password: a string representing the password to authenticate with the Redis server. It defaults to an empty string.
type Redis struct {
	Address  string `json:"address,omitempty"`
	Password string `json:"password,omitempty"`
}

// Cache represents the cache configuration in the Gopen struct.
type Cache struct {
	// Duration is a string field in the Cache struct that represents the duration of the cache.
	// It is specified in a format compatible with Go's time.ParseDuration function.
	// The default value is an empty string. If not provided, the duration will be 30s.
	Duration string `json:"duration,omitempty"`
	// StrategyHeaders represents a slice of strings that contains the headers used to determine the cache strategy key.
	StrategyHeaders []string `json:"strategy-headers,omitempty"`
	// OnlyIfStatusCodes represents a slice of integers representing the HTTP status codes
	// for which the cache should be used. Default is an empty slice. If not provided,
	// the default value is 2xx success HTTP status codes.
	OnlyIfStatusCodes []int `json:"only-if-status-codes,omitempty"`
	// OnlyIfMethods represents a slice of strings that contains the HTTP methods for which the cache should be used.
	// The default value is an empty slice. If not provided, the default value is an empty slice.
	// Example: []string{"GET", "POST"}
	OnlyIfMethods []string `json:"only-if-methods,omitempty"`
	// AllowCacheControl represents a pointer to a boolean indicating whether the cache should
	// honor the Cache-Control header. It defaults to nil. If not provided, the default value is false.
	AllowCacheControl *bool `json:"allow-cache-control,omitempty"`
}

// EndpointCache represents the cache configuration for an endpoint.
type EndpointCache struct {
	// Enabled represents a boolean indicating whether caching is enabled for an endpoint.
	Enabled bool `json:"enabled,omitempty"`
	// IgnoreQuery represents a boolean indicating whether to ignore query parameters when caching.
	IgnoreQuery bool `json:"ignore-query,omitempty"`
	// Duration represents the duration configuration for caching an endpoint response.
	Duration string `json:"duration,omitempty"`
	// StrategyHeaders represents a slice of strings for strategy headers
	StrategyHeaders []string `json:"strategy-headers,omitempty"`
	// OnlyIfStatusCodes represents the status codes that the cache should be applied to.
	OnlyIfStatusCodes []int `json:"only-if-status-codes,omitempty"`
	// AllowCacheControl represents a boolean value indicating whether the cache control header is allowed for the endpoint cache.
	AllowCacheControl *bool `json:"allow-cache-control,omitempty"`
}

// Limiter represents the configuration for rate limiting in the Gopen application.
type Limiter struct {
	// MaxHeaderSize represents the maximum size of the header in bytes for rate limiting.
	MaxHeaderSize string `json:"max-header-size,omitempty"`
	// MaxBodySize represents the maximum size of the body in bytes for rate limiting.
	MaxBodySize string `json:"max-body-size,omitempty"`
	// MaxMultipartMemorySize represents the maximum memory size for multipart request bodies.
	MaxMultipartMemorySize string `json:"max-multipart-memory-size,omitempty"`
	// Rate represents the configuration for rate limiting in the Limiter struct. It specifies the capacity and
	// frequency of allowed requests.
	Rate Rate `json:"rate,omitempty"`
}

// Rate represents the configuration for rate limiting. It specifies the capacity
// and frequency of allowed requests.
type Rate struct {
	// Capacity represents the maximum number of allowed requests within a given time period.
	Capacity int `json:"capacity,omitempty"`
	// Every represents the frequency of allowed requests in the Rate configuration for rate limiting.
	Every string `json:"every,omitempty"`
}

// SecurityCors represents the configuration options for Cross-Origin Resource Sharing (CORS) settings in Gopen.
type SecurityCors struct {
	// AllowOrigins represents the allowed origins for Cross-Origin Resource Sharing (CORS) settings in Gopen.
	AllowOrigins []string `json:"allow-origins"`
	// AllowMethods represents the allowed HTTP methods for Cross-Origin Resource Sharing (CORS) settings in Gopen.
	AllowMethods []string `json:"allow-methods"`
	// AllowHeaders represents the list of allowed headers for Cross-Origin Resource Sharing (CORS) settings in Gopen.
	AllowHeaders []string `json:"allow-headers"`
}

// Endpoint represents the configuration for an API endpoint in the Gopen application.
type Endpoint struct {
	// Path is a string representing the path of the API endpoint. It is a field in the Endpoint struct.
	Path string `json:"path,omitempty"`
	// Method represents the HTTP method of an API endpoint.
	Method string `json:"method,omitempty"`
	// Timeout represents the timeout duration for the API endpoint.
	// It is a string value specified in the JSON configuration.
	// The default value is empty. If not provided, the timeout will be Gopen.Timeout.
	Timeout string `json:"timeout,omitempty"`
	// Limiter represents the configuration for rate limiting in the Gopen application.
	// The default value is nil. If not provided, the limiter will be Gopen.Limiter.
	Limiter Limiter `json:"limiter,omitempty"`
	// Cache represents the cache configuration for an endpoint.
	// The default value is EndpointCache empty with enabled false.
	Cache EndpointCache `json:"cache,omitempty"`
	// ResponseEncode represents the encoding format for the API endpoint response. The ResponseEncode
	// field is an enum.ResponseEncode value, which can have one of the following values:
	// - enum.ResponseEncodeText: for encoding the response as plain text.
	// - enum.ResponseEncodeJson: for encoding the response as JSON.
	// - enum.ResponseEncodeXml: for encoding the response as XML.
	// The default value is empty. If not provided, the response will be encoded by type, if the string is json it
	// returns json, otherwise it responds to plain text
	ResponseEncode enum.ResponseEncode `json:"response-encode,omitempty"`
	// AggregateResponses represents a boolean indicating whether the API endpoint should aggregate responses
	// from multiple backends.
	AggregateResponses bool `json:"aggregate-responses,omitempty"`
	// AbortIfStatusCodes represents a slice of integers representing the HTTP status codes
	// for which the API endpoint should abort. It is a field in the Endpoint struct.
	AbortIfStatusCodes []int `json:"abort-if-status-codes,omitempty"`
	// Beforeware represents a slice of strings containing the names of the beforeware middlewares that should be
	// applied before processing the API endpoint.
	Beforeware []string `json:"beforeware,omitempty"`
	// Afterware represents the configuration for the afterware middlewares to apply after processing the API endpoint.
	// It is a slice of strings representing the names of the afterware middlewares to apply.
	// The names specify the behavior and settings of each afterware middleware.
	// If not provided, the default value is an empty slice.
	// The afterware middleware is executed after processing the API endpoint, allowing for modification or
	// transformation of the response or performing any additional actions.
	// Afterware can be used for logging, error handling, response modification, etc.
	Afterware []string `json:"afterware,omitempty"`
	// Backends represents the backend configurations for an API endpoint in the Gopen application.
	// It is a slice of Backend structs.
	Backends []Backend `json:"backends,omitempty"`
}

// Backend represents the configuration for a backend in the Gopen application.
type Backend struct {
	// Name is a field in the Backend struct that represents the name of the backend configuration.
	Name string `json:"name,omitempty"`
	// Hosts represents a slice of strings that specifies the hosts for a backend configuration.
	Hosts []string `json:"hosts,omitempty"`
	// Path is a field in the Backend struct that represents the path for a backend request.
	// Example: "/api/users"
	Path string `json:"path,omitempty"`
	// Method represents the HTTP method for a backend request in the Gopen application.
	// It is a field in the Backend struct and is specified in the Gopen configuration JSON file.
	// The value should be a string and can be one of the following: "GET", "POST", "PUT", "PATCH", "DELETE".
	// It is used to specify the HTTP method for the backend request.
	Method string `json:"method,omitempty"`
	// ForwardHeaders is a field in the Backend struct that represents a list of headers to be forwarded
	// in the backend request. It is specified as a slice of strings in the Gopen configuration JSON file.
	// Each string represents a header name.
	// Example: ["Content-Type", "User-Agent"]
	ForwardHeaders []string `json:"forward-headers,omitempty"`
	// ForwardQueries represents the list of query parameters that will be forwarded to the backend server.
	// The query parameters are specified as string elements in a slice.
	// It is a field in the Backend struct and is specified in the Gopen configuration JSON file.
	// The ForwardQueries field is used to specify which query parameters of the incoming request will be included in the
	// request sent to the backend server.
	ForwardQueries []string `json:"forward-queries,omitempty"`
	// Modifiers represent the configuration to modify the request and response of a backend and endpoint in the Gopen application.
	Modifiers BackendModifiers `json:"modifiers,omitempty"`
	// ExtraConfig represents additional configuration options for a backend in the Gopen application.
	ExtraConfig BackendExtraConfig `json:"extra-config,omitempty"`
}

// BackendModifiers represents a set of modifiers that can be applied to different parts of the request and response
// in the Gopen application.
type BackendModifiers struct {
	// StatusCode represents a modifier that can be applied to the status code of Backend response.
	StatusCode Modifier `json:"status-code,omitempty"`
	// Header represents a slice of modifying structures that can be applied to the header of a request or response from
	// the Endpoint or just the current Backend.
	Header []Modifier `json:"header,omitempty"`
	// Params is a slice of modifiers that can be applied to the parameters of a request from the Endpoint
	// or just the current Backend.
	Params []Modifier `json:"params,omitempty"`
	// Query represents a slice of Modifier structs that can be applied to the query parameters of a request
	// from the Endpoint or just the current Backend.
	Query []Modifier `json:"query,omitempty"`
	// Body represents a slice of Modifier structs that can be applied to the body of a request or response
	// from the Endpoint or just the current Backend.
	Body []Modifier `json:"body,omitempty"`
}

// BackendExtraConfig represents additional configuration options for a backend in the Gopen application.
// - OmitRequestBody: a boolean flag indicating whether the backend should omit the request body in the outgoing request.
// If set to true, the backend will not include the request body in the outgoing request.
// If set to false, the request body will be included in the outgoing request.
// The default value is false.
// - OmitResponse: a boolean flag indicating whether the backend should omit the response in the incoming request.
// If set to true, the backend will not include the response in the incoming request.
// If set to false, the response will be included in the incoming request.
// The default value is false.
type BackendExtraConfig struct {
	// GroupResponse is a boolean flag indicating whether the backend should group response.
	// The default value is false.
	GroupResponse bool `json:"group-response,omitempty"`
	// OmitRequestBody represents a boolean flag indicating whether the backend should omit the request body in request.
	// If set to true, the backend will not include the request body in the request.
	// If set to false, the request body will be included in the request. The default value is false.
	OmitRequestBody bool `json:"omit-request-body,omitempty"`
	// OmitResponse represents a boolean flag indicating whether the backend should omit the response in the incoming request.
	// If set to true, the backend will not include the response in the incoming request.
	// If set to false, the response will be included in the incoming request.
	// The default value is false.
	OmitResponse bool `json:"omit-response,omitempty"`
}

// Modifier represents a modification that can be applied to a request or response in the Gopen application.
type Modifier struct {
	// Context represents the context in which a modification should be applied.
	// It is an enum.ModifierContext value.
	// Valid values for Context are "request" and "response".
	Context enum.ModifierContext `json:"context,omitempty"`
	// Scope represents the scope of a modification in the Backend or Endpoint.
	// It is an enum.ModifierScope value that specifies where the modification should be applied.
	// Valid values for Scope are "request" and "response".
	Scope enum.ModifierScope `json:"scope,omitempty"`
	// Action represents the action to be performed in the Modifier struct.
	// It is an enum.ModifierAction value and can be one of the following values:
	// - ModifierActionSet: to set a value.
	// - ModifierActionAdd: to add a value.
	// - ModifierActionDel: to delete a value.
	// - ModifierActionReplace: to replace a value.
	// - ModifierActionRename: to rename a value.
	Action enum.ModifierAction `json:"action,omitempty"`
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
