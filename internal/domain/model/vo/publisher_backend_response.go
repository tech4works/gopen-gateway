package vo

import (
	"net/http"

	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type PublisherBackendResponse struct {
	outcome enum.BackendOutcome

	ok   bool
	body *Body
}

func NewPublisherBackendResponse(outcome enum.BackendOutcome, ok bool, body *Body) *PublisherBackendResponse {
	return &PublisherBackendResponse{
		outcome: outcome,
		ok:      ok,
		body:    body,
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

func (p *PublisherBackendResponse) HasBody() bool {
	return checker.NonNil(p.body) && checker.IsGreaterThan(p.body.Size(), 0)
}

func (p *PublisherBackendResponse) Body() *Body {
	return p.body
}

func (p *PublisherBackendResponse) Executed() bool {
	return checker.Equals(p.outcome, enum.BackendOutcomeExecuted)
}

func (p *PublisherBackendResponse) Cancelled() bool {
	return checker.Equals(p.outcome, enum.BackendOutcomeCancelled)
}

func (p *PublisherBackendResponse) Ignored() bool {
	return checker.Equals(p.outcome, enum.BackendOutcomeIgnored)
}

func (p *PublisherBackendResponse) Map() (map[string]any, error) {
	return map[string]any{
		"ok":   p.ok,
		"body": p.body,
	}, nil
}
