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
	Setting Gopen `json:"setting"`
}
