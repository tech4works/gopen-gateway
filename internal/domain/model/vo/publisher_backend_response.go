package vo

import (
	"net/http"

	"github.com/tech4works/checker"
)

type PublisherBackendResponse struct {
	ok   bool
	body *Body
}

func NewPublisherBackendResponse(ok bool, body *Body) *PublisherBackendResponse {
	return &PublisherBackendResponse{
		ok:   ok,
		body: body,
	}
}

func (p *PublisherBackendResponse) OK() bool {
	return p.ok
}

func (p *PublisherBackendResponse) StatusCode() StatusCode {
	if p.ok {
		return NewStatusCode(http.StatusCreated)
	}
	return NewStatusCode(http.StatusInternalServerError)
}

func (p *PublisherBackendResponse) Header() Header {
	return NewEmptyHeader()
}

func (p *PublisherBackendResponse) Body() *Body {
	return p.body
}

func (p *PublisherBackendResponse) HasBody() bool {
	return checker.NonNil(p.body) && checker.IsGreaterThan(p.body.Size(), 0)
}

func (p *PublisherBackendResponse) Map() (map[string]any, error) {
	return map[string]any{
		"ok":   p.ok,
		"body": p.body,
	}, nil
}
