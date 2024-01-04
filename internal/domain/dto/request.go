package dto

import (
	"io"
	"net/http"
	"net/url"
)

type BackendRequest struct {
	Endpoint   string            `json:"endpoint,omitempty"`
	Header     http.Header       `json:"header,omitempty"`
	Query      url.Values        `json:"query,omitempty"`
	Params     map[string]string `json:"params,omitempty"`
	Body       any               `json:"body,omitempty"`
	BodyToSend io.ReadCloser     `json:"bodyToSend,omitempty"`
}

type BackendResponse struct {
	StatusCode int         `json:"statusCode"`
	Header     http.Header `json:"header,omitempty"`
	Body       any         `json:"body,omitempty"`
	Group      string      `json:"group,omitempty"`
	Remove     bool        `json:"remove,omitempty"`
}
