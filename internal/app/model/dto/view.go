package dto

type SettingView struct {
	Version     string    `json:"version,omitempty"`
	VersionDate string    `json:"version-date,omitempty"`
	Founder     string    `json:"founder,omitempty"`
	CodeHelpers string    `json:"code-helpers,omitempty"`
	Endpoints   int       `json:"endpoints"`
	Middlewares int       `json:"middlewares"`
	Backends    int       `json:"backends"`
	Modifiers   int       `json:"modifiers"`
	Setting     GOpenView `json:"setting"`
}
