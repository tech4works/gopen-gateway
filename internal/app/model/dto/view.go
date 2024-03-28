package dto

import "time"

type ErrorView struct {
	File      string    `json:"file,omitempty"`
	Line      int       `json:"line,omitempty"`
	Endpoint  string    `json:"endpoint,omitempty"`
	Message   string    `json:"message,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

type ConfigView struct {
	Version     string `json:"version,omitempty"`
	VersionDate string `json:"version-date,omitempty"`
	Founder     string `json:"founder,omitempty"`
	CodeHelpers string `json:"code-helpers,omitempty"`
	Endpoints   int    `json:"endpoints"`
	Middlewares int    `json:"middlewares"`
	Backends    int    `json:"backends"`
	Modifiers   int    `json:"modifiers"`
	Config      GOpen  `json:"config,omitempty"`
}
