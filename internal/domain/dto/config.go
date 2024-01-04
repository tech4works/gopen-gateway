package dto

type Config struct {
	Port        int         `json:"port,omitempty" validate:"required"`
	Cache       string      `json:"cache,omitempty" validate:"omitempty"`
	Timeout     Timeout     `json:"timeout,omitempty" validate:"required"`
	Limiter     Limiter     `json:"limiter,omitempty" validate:"omitempty"`
	ExtraConfig ExtraConfig `json:"extra-config,omitempty" validate:"omitempty"`
	Endpoints   []Endpoint  `json:"endpoints,omitempty" validate:"required"`
}

type Timeout struct {
	ReadHeader string `json:"read-header,omitempty"`
	Read       string `json:"read,omitempty"`
	Handler    string `json:"handler,omitempty" validate:"required"`
	Write      string `json:"write,omitempty"`
}

type Limiter struct {
	MaxSizeRequestHeader   string `json:"max-size-request-header,omitempty"`
	MaxSizeRequestBody     string `json:"max-size-request-body,omitempty"`
	MaxSizeMultipartMemory string `json:"max-size-multipart-memory,omitempty"`
	MaxIpRequestPerSeconds int    `json:"max-ip-request-per-seconds,omitempty"`
}

type ExtraConfig struct {
	SecurityCors   SecurityCors       `json:"security-cors,omitempty"`
	Authorizations map[string]Backend `json:"authorizations,omitempty"`
}

type SecurityCors struct {
	AllowCountries []string `json:"allow-countries,omitempty"`
	AllowOrigins   []string `json:"allow-origins,omitempty"`
	AllowMethods   []string `json:"allow-methods,omitempty"`
	AllowHeaders   []string `json:"allow-headers,omitempty"`
}

type Endpoint struct {
	Endpoint           string    `json:"endpoint,omitempty" validate:"required"`
	Method             string    `json:"method,omitempty" validate:"required"`
	Cacheable          bool      `json:"cacheable,omitempty"`
	AggregateResponses bool      `json:"aggregate-responses,omitempty"`
	AbortSequential    bool      `json:"abort-sequential,omitempty"`
	Backends           []Backend `json:"backends,omitempty" validate:"required"`
}

type Backend struct {
	Host           []string   `json:"host,omitempty" validate:"required"`
	Endpoint       string     `json:"endpoint,omitempty" validate:"required,url"`
	Method         string     `json:"method,omitempty" validate:"required"`
	Group          string     `json:"group,omitempty"`
	RemoveResponse bool       `json:"remove-response,omitempty"`
	ForwardHeaders []string   `json:"forward-headers,omitempty" validate:"required,min=1"`
	Authorizations []string   `json:"authorizations,omitempty"`
	Query          []string   `json:"query,omitempty"`
	Headers        []Modifier `json:"headers,omitempty"`
	Params         []Modifier `json:"params,omitempty"`
	Queries        []Modifier `json:"queries,omitempty"`
	Body           []Modifier `json:"body,omitempty"`
}

type Modifier struct {
	Scope  []string `json:"scope,omitempty" validate:"required"`
	Action string   `json:"action,omitempty" validate:"required"`
	Key    string   `json:"key,omitempty" validate:"required"`
	Value  string   `json:"value,omitempty" validate:"required"`
}
