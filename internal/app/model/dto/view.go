package dto

// SettingView represents the configuration view for the application.
// It includes properties such as version, founder, code helpers, and various counts.
// It also contains a setting object of type GopenView, which represents the detailed configuration.
type SettingView struct {
	// Version represents the version of the application.
	Version string `json:"version,omitempty"`
	// VersionDate represents the date of the application version.
	VersionDate string `json:"version-date,omitempty"`
	// Founder represents the founder of a software application.
	Founder string `json:"founder,omitempty"`
	// CodeHelpers represents the code helpers configuration in the SettingView struct.
	// It is a string field that represents the code helpers for the application.
	CodeHelpers string `json:"code-helpers,omitempty"`
	// Endpoints represents the number of APIs in the Gopen application.
	Endpoints int `json:"endpoints"`
	// Middlewares represents the number of middlewares in the SettingView struct.
	// It is an integer field that specifies the count of middlewares used in the Gopen application.
	Middlewares int `json:"middlewares"`
	// Backends represents the number of backends configured in the SettingView struct.
	Backends int `json:"backends"`
	// Modifiers represents the count of modifiers in the SettingView struct.
	Modifiers int `json:"modifiers"`
	// Setting represents the detailed configuration view for the Gopen application.
	Setting GopenView `json:"setting"`
}

// GopenView represents the configuration view for the Gopen application.
type GopenView struct {
	// Version represents the version of the gopen json configuration.
	Version string `json:"version,omitempty"`
	// Port represents the port of the gopen json configuration.
	Port int `json:"port,omitempty"`
	// HotReload allows for dynamic reloading of the json configuration.
	HotReload bool `json:"hot-reload,omitempty"`
	// Timeout represents the timeout duration for an endpoint of the json configuration.
	Timeout string `json:"timeout,omitempty"`
	// Limiter represents rate limiting configuration settings.
	Limiter *Limiter `json:"limiter,omitempty"`
	// Cache represents the caching configuration for an endpoint or the application.
	Cache *Cache `json:"cache,omitempty"`
	// SecurityCors represents the Cross-Origin Resource Sharing (CORS) configuration for security.
	SecurityCors *SecurityCors `json:"security-cors,omitempty"`
	// Middlewares is a map that represents the configuration for middleware backends.
	// The key is the name of the middleware, and the value is a Backend object that specifies the properties of the middleware.
	Middlewares map[string]Backend `json:"middlewares,omitempty"`
	// Endpoints represents a collection of Endpoint objects.
	// Each Endpoint object represents a specific path and method configuration.
	//
	// The Endpoints field is used in the GopenView struct to define the various endpoints of the application.
	// It is a slice of Endpoint objects.
	Endpoints []Endpoint `json:"endpoints,omitempty"`
}
