package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
)

type HTTPResponse struct {
	statusCode StatusCode
	header     Header
	body       *Body
}

func NewHTTPResponse(statusCode StatusCode, header Header, body *Body) *HTTPResponse {
	return &HTTPResponse{
		statusCode: statusCode,
		header:     header,
		body:       body,
	}
}

func NewHTTPResponseStatusCode(statusCode StatusCode, header Header) *HTTPResponse {
	return &HTTPResponse{
		statusCode: statusCode,
		header:     header,
	}
}

func (h *HTTPResponse) StatusCode() StatusCode {
	return h.statusCode
}

func (h *HTTPResponse) Header() Header {
	return h.header
}

func (h *HTTPResponse) Body() *Body {
	return h.body
}

func (h *HTTPResponse) HasBody() bool {
	return helper.IsNotNil(h.body)
}
