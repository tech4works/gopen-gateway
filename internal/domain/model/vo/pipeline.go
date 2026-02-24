package vo

import (
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
)

type HostPipelineSpec struct {
	hosts []string
}

type HeaderPipelineSpec struct {
	omit bool

	mapper    *Mapper
	projector *Projector
	modifiers []Modifier
}

type URLPathPipelineSpec struct {
	modifiers []Modifier
}

type QueryPipelineSpec struct {
	omit bool

	mapper    *Mapper
	projector *Projector
	modifiers []Modifier
}

type BodyPipelineSpec struct {
	omit      bool
	omitEmpty bool

	group string

	nomenclature    enum.Nomenclature
	contentType     enum.ContentType
	contentEncoding enum.ContentEncoding

	mapper    *Mapper
	projector *Projector
	modifiers []Modifier
	joins     []Join
}

type GroupIDPipelineSpec struct {
	value string
}

type DeduplicationIDPipelineSpec struct {
	value string
}

type PublisherAttributesPipelineSpec struct {
	attributes map[string]PublisherMessageAttribute
}

func NewHostPipelineSpec(hosts []string) *HostPipelineSpec {
	return &HostPipelineSpec{hosts: hosts}
}

func NewHeaderPipelineSpecFromBackendRequest(b *BackendRequest) *HeaderPipelineSpec {
	if checker.IsNil(b) || !b.HasHeader() {
		return nil
	}
	return &HeaderPipelineSpec{
		omit:      b.Header().Omit(),
		mapper:    b.Header().Mapper(),
		projector: b.Header().Projector(),
		modifiers: b.Header().Modifiers(),
	}
}

func NewHeaderPipelineSpecFromBackendResponse(b *BackendResponse) *HeaderPipelineSpec {
	if checker.IsNil(b) || !b.HasHeader() {
		return nil
	}
	return &HeaderPipelineSpec{
		omit:      b.Header().Omit(),
		mapper:    b.Header().Mapper(),
		projector: b.Header().Projector(),
		modifiers: b.Header().Modifiers(),
	}
}

func NewURLPathPipelineSpecFromBackendRequest(b *BackendRequest) *URLPathPipelineSpec {
	if checker.IsNil(b) || !b.HasURLPath() {
		return nil
	}
	return &URLPathPipelineSpec{
		modifiers: b.URLPath().Modifiers(),
	}
}

func NewQueryPipelineSpecFromBackendRequest(b *BackendRequest) *QueryPipelineSpec {
	if checker.IsNil(b) || !b.HasQuery() {
		return nil
	}
	return &QueryPipelineSpec{
		omit:      b.Query().Omit(),
		mapper:    b.Query().Mapper(),
		projector: b.Query().Projector(),
		modifiers: b.Query().Modifiers(),
	}
}

func NewBodyPipelineSpecFromBackendRequest(b *BackendRequest) *BodyPipelineSpec {
	if checker.IsNil(b) || !b.HasBody() {
		return nil
	}
	return &BodyPipelineSpec{
		omit:            b.Body().Omit(),
		omitEmpty:       b.Body().OmitEmpty(),
		nomenclature:    b.Body().Nomenclature(),
		contentType:     b.Body().ContentType(),
		contentEncoding: b.Body().ContentEncoding(),
		mapper:          b.Body().Mapper(),
		projector:       b.Body().Projector(),
		modifiers:       b.Body().Modifiers(),
		joins:           b.Body().Joins(),
	}
}

func NewBodyPipelineSpecFromPublisherMessage(b *PublisherMessage) *BodyPipelineSpec {
	if checker.IsNil(b) || !b.HasBody() {
		return nil
	}
	return &BodyPipelineSpec{
		omitEmpty:       b.Body().OmitEmpty(),
		nomenclature:    b.Body().Nomenclature(),
		contentType:     b.Body().ContentType(),
		contentEncoding: b.Body().ContentEncoding(),
		mapper:          b.Body().Mapper(),
		projector:       b.Body().Projector(),
		modifiers:       b.Body().Modifiers(),
		joins:           b.Body().Joins(),
	}
}

func NewBodyPipelineSpecFromBackendResponse(b *BackendResponse) *BodyPipelineSpec {
	if checker.IsNil(b) || !b.HasBody() {
		return nil
	}
	return &BodyPipelineSpec{
		omit:      b.Body().Omit(),
		group:     b.Body().Group(),
		mapper:    b.Body().Mapper(),
		projector: b.Body().Projector(),
		modifiers: b.Body().Modifiers(),
		joins:     b.Body().Joins(),
	}
}

func NewBodyPipelineSpecFromEndpointResponse(e *EndpointResponse) *BodyPipelineSpec {
	if checker.IsNil(e) || !e.HasBody() {
		return nil
	}
	return &BodyPipelineSpec{
		omitEmpty:       e.Body().OmitEmpty(),
		nomenclature:    e.Body().Nomenclature(),
		contentType:     e.Body().ContentType(),
		contentEncoding: e.Body().ContentEncoding(),
		mapper:          e.Body().Mapper(),
		projector:       e.Body().Projector(),
	}
}

func NewGroupIDPipelineSpecFromPublisher(p *Publisher) *GroupIDPipelineSpec {
	if checker.IsNil(p) {
		return nil
	}
	return &GroupIDPipelineSpec{value: p.GroupID()}
}

func NewDeduplicationIDPipelineSpecFromPublisher(p *Publisher) *DeduplicationIDPipelineSpec {
	if checker.IsNil(p) {
		return nil
	}
	return &DeduplicationIDPipelineSpec{value: p.DeduplicationID()}
}

func NewPublisherAttributesPipelineSpecFromPublisher(p *Publisher) *PublisherAttributesPipelineSpec {
	if checker.IsNil(p) || !p.HasMessage() {
		return nil
	}
	return &PublisherAttributesPipelineSpec{attributes: p.Message().Attributes()}
}

func (s HostPipelineSpec) Hosts() []string { return s.hosts }

func (s HeaderPipelineSpec) Omit() bool            { return s.omit }
func (s HeaderPipelineSpec) Mapper() *Mapper       { return s.mapper }
func (s HeaderPipelineSpec) Projector() *Projector { return s.projector }
func (s HeaderPipelineSpec) Modifiers() []Modifier { return s.modifiers }

func (s URLPathPipelineSpec) Modifiers() []Modifier { return s.modifiers }

func (s QueryPipelineSpec) Omit() bool            { return s.omit }
func (s QueryPipelineSpec) Mapper() *Mapper       { return s.mapper }
func (s QueryPipelineSpec) Projector() *Projector { return s.projector }
func (s QueryPipelineSpec) Modifiers() []Modifier { return s.modifiers }

func (s BodyPipelineSpec) Omit() bool                            { return s.omit }
func (s BodyPipelineSpec) OmitEmpty() bool                       { return s.omitEmpty }
func (s BodyPipelineSpec) HasGroup() bool                        { return checker.IsNotEmpty(s.group) }
func (s BodyPipelineSpec) Group() string                         { return s.group }
func (s BodyPipelineSpec) Nomenclature() enum.Nomenclature       { return s.nomenclature }
func (s BodyPipelineSpec) ContentType() enum.ContentType         { return s.contentType }
func (s BodyPipelineSpec) ContentEncoding() enum.ContentEncoding { return s.contentEncoding }
func (s BodyPipelineSpec) Mapper() *Mapper                       { return s.mapper }
func (s BodyPipelineSpec) Projector() *Projector                 { return s.projector }
func (s BodyPipelineSpec) Modifiers() []Modifier                 { return s.modifiers }
func (s BodyPipelineSpec) Joins() []Join                         { return s.joins }

func (s *GroupIDPipelineSpec) Value() string { return s.value }

func (s *DeduplicationIDPipelineSpec) Value() string { return s.value }

func (s *PublisherAttributesPipelineSpec) Attributes() map[string]PublisherMessageAttribute {
	return s.attributes
}
