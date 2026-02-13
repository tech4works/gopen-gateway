package vo

import (
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type PublisherBackendRequest struct {
	provider        enum.PublisherProvider
	path            string
	groupID         *string
	deduplicationID *string
	delay           Duration
	attributes      map[string]PublisherMessageAttribute
	body            string
}

func NewPublisherBackendRequest(
	provider enum.PublisherProvider,
	path,
	groupID,
	deduplicationID string,
	delay Duration,
	attributes map[string]PublisherMessageAttribute,
	body string,
) *PublisherBackendRequest {
	return &PublisherBackendRequest{
		provider:        provider,
		path:            path,
		groupID:         checker.IfEmptyReturns(&groupID, nil),
		deduplicationID: checker.IfEmptyReturns(&deduplicationID, nil),
		delay:           delay,
		attributes:      attributes,
		body:            body,
	}
}

func (m PublisherBackendRequest) Provider() enum.PublisherProvider {
	return m.provider
}

func (m PublisherBackendRequest) Path() string {
	return m.path
}

func (m PublisherBackendRequest) GroupID() *string {
	return m.groupID
}

func (m PublisherBackendRequest) DeduplicationID() *string {
	return m.deduplicationID
}

func (m PublisherBackendRequest) Delay() Duration {
	return m.delay
}

func (m PublisherBackendRequest) HasAttributes() bool {
	return checker.IsNotNilOrEmpty(m.attributes)
}

func (m PublisherBackendRequest) Attributes() map[string]PublisherMessageAttribute {
	return m.attributes
}

func (m PublisherBackendRequest) Body() string {
	return m.body
}
