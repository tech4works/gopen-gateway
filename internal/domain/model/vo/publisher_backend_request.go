package vo

import (
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type PublisherBackendRequest struct {
	degradation     Degradation
	broker          enum.BackendBroker
	path            string
	groupID         *string
	deduplicationID *string
	delay           Duration
	attributes      map[string]AttributeValueConfig
	body            *Payload
}

func NewPublisherBackendRequest(
	degradation Degradation,
	broker enum.BackendBroker,
	path,
	groupID,
	deduplicationID string,
	delay Duration,
	attributes map[string]AttributeValueConfig,
	body *Payload,
) *PublisherBackendRequest {
	return &PublisherBackendRequest{
		degradation:     degradation,
		broker:          broker,
		path:            path,
		groupID:         checker.IfEmptyReturns(&groupID, nil),
		deduplicationID: checker.IfEmptyReturns(&deduplicationID, nil),
		delay:           delay,
		attributes:      attributes,
		body:            body,
	}
}

func (m PublisherBackendRequest) Degradation() Degradation {
	return m.degradation
}

func (m PublisherBackendRequest) Degraded() bool {
	return m.Degradation().Any()
}

func (m PublisherBackendRequest) GroupIDDegraded() bool {
	return m.Degradation().Has(enum.DegradationKindGroupID)
}

func (m PublisherBackendRequest) DeduplicationIDDegraded() bool {
	return m.Degradation().Has(enum.DegradationKindDeduplicationID)
}

func (m PublisherBackendRequest) AttributesDegraded() bool {
	return m.Degradation().Has(enum.DegradationKindAttributes)
}

func (m PublisherBackendRequest) BodyDegraded() bool {
	return m.Degradation().Has(enum.DegradationKindAttributes)
}

func (m PublisherBackendRequest) Broker() enum.BackendBroker {
	return m.broker
}

func (m PublisherBackendRequest) Path() string {
	return m.path
}

func (m PublisherBackendRequest) HasGroupID() bool {
	return checker.IsNotEmpty(m.groupID)
}

func (m PublisherBackendRequest) GroupID() *string {
	return m.groupID
}

func (m PublisherBackendRequest) HasDeduplicationID() bool {
	return checker.IsNotEmpty(m.deduplicationID)
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

func (m PublisherBackendRequest) Attributes() map[string]AttributeValueConfig {
	return m.attributes
}

func (m PublisherBackendRequest) HasBody() bool {
	return checker.NonNil(m.body)
}

func (m PublisherBackendRequest) Body() *Payload {
	return m.body
}
