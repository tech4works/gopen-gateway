package service

import (
	"math/rand"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"

	"github.com/tech4works/gopen-gateway/internal/domain/model/aggregate"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type BuildPipeline struct {
	modifierService     Modifier
	joinService         Join
	mapperService       Mapper
	projectorService    Projector
	omitterService      Omitter
	nomenclatureService Nomenclature
	contentService      Content
	aggregatorService   Aggregator
	dynamicValueService DynamicValue
}

func NewBuildPipeline(
	modifierService Modifier,
	joinService Join,
	mapperService Mapper,
	projectorService Projector,
	omitterService Omitter,
	nomenclatureService Nomenclature,
	contentService Content,
	aggregatorService Aggregator,
	dynamicValueService DynamicValue,
) BuildPipeline {
	return BuildPipeline{
		modifierService:     modifierService,
		joinService:         joinService,
		mapperService:       mapperService,
		projectorService:    projectorService,
		omitterService:      omitterService,
		nomenclatureService: nomenclatureService,
		contentService:      contentService,
		aggregatorService:   aggregatorService,
		dynamicValueService: dynamicValueService,
	}
}

func (p BuildPipeline) ApplyHost(
	spec vo.HostPipelineSpec,
) string {
	hosts := spec.Hosts()
	if checker.IsLengthEquals(hosts, 1) {
		return hosts[0]
	}
	return hosts[rand.Intn(len(hosts))]
}

func (p BuildPipeline) ApplyMetadata(
	spec vo.MetadataPipelineSpec,
	metadata vo.Metadata,
	ignoreKeys []string,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (vo.Metadata, []error) {
	if checker.IsNil(spec) || spec.Omit() {
		return metadata, nil
	}
	return apply(
		metadata,
		step[vo.Metadata]{
			label: "modifiers",
			run: func(in vo.Metadata) (vo.Metadata, []error) {
				return p.modifierService.ExecuteMetadataModifiers(spec.Modifiers(), in, ignoreKeys, request, history)
			},
		},
		step[vo.Metadata]{
			label: "mapper",
			run: func(in vo.Metadata) (vo.Metadata, []error) {
				out, err := p.mapperService.MapMetadata(spec.Mapper(), in, ignoreKeys, request, history)
				return out, converter.ToSliceIfNonNil(err)
			},
		},
		step[vo.Metadata]{
			label: "projector",
			run: func(in vo.Metadata) (vo.Metadata, []error) {
				out, err := p.projectorService.ProjectMetadata(spec.Projector(), in, ignoreKeys, request, history)
				return out, converter.ToSliceIfNonNil(err)
			},
		},
	)
}

func (p BuildPipeline) ApplyURLPath(
	spec vo.URLPathPipelineSpec,
	urlPath vo.URLPath,
	request *vo.EndpointRequest,
	hist *aggregate.History,
) (vo.URLPath, []error) {
	if checker.IsNil(spec) {
		return urlPath, nil
	}
	return apply(
		urlPath,
		step[vo.URLPath]{
			label: "modifiers",
			run: func(in vo.URLPath) (vo.URLPath, []error) {
				return p.modifierService.ExecuteURLPathModifiers(spec.Modifiers(), in, request, hist)
			},
		},
	)
}

func (p BuildPipeline) ApplyQuery(
	spec vo.QueryPipelineSpec,
	query vo.Query,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (vo.Query, []error) {
	if checker.IsNil(spec) {
		return query, nil
	} else if spec.Omit() {
		return vo.NewEmptyQuery(), nil
	}
	return apply(
		query,
		step[vo.Query]{
			label: "modifiers",
			run: func(in vo.Query) (vo.Query, []error) {
				return p.modifierService.ExecuteQueryModifiers(spec.Modifiers(), in, request, history)
			},
		},
		step[vo.Query]{
			label: "mapper",
			run: func(in vo.Query) (vo.Query, []error) {
				out, err := p.mapperService.MapQuery(spec.Mapper(), in, request, history)
				return out, converter.ToSliceIfNonNil(err)
			},
		},
		step[vo.Query]{
			label: "projector",
			run: func(in vo.Query) (vo.Query, []error) {
				out, err := p.projectorService.ProjectQuery(spec.Projector(), in, request, history)
				return out, converter.ToSliceIfNonNil(err)
			},
		},
	)
}

func (p BuildPipeline) ApplyPayload(
	spec vo.PayloadPipelineSpec,
	payload *vo.Payload,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (*vo.Payload, []error) {
	if checker.IsNil(spec) {
		return payload, nil
	} else if spec.Omit() || checker.IsNil(payload) {
		return nil, nil
	} else if payload.IsNotValid() || payload.ContentType().IsUnsupported() {
		return payload, errors.NewAsSlicef("payload is not valid or unsupported to transforms content-type=%s isValid=%v",
			payload.ContentType().String(), payload.IsValid())
	}
	return apply(
		payload,
		step[*vo.Payload]{
			label: "modifiers",
			run: func(in *vo.Payload) (*vo.Payload, []error) {
				return p.modifierService.ExecutePayloadModifiers(spec.Modifiers(), in, request, history)
			},
		},
		step[*vo.Payload]{
			label: "joins",
			run: func(in *vo.Payload) (*vo.Payload, []error) {
				return p.joinService.ExecutePayloadJoins(spec.Joins(), in, request, history)
			},
		},
		step[*vo.Payload]{
			label: "mapper",
			run: func(in *vo.Payload) (*vo.Payload, []error) {
				return p.mapperService.MapPayload(spec.Mapper(), in, request, history)
			},
		},
		step[*vo.Payload]{
			label: "projector",
			run: func(in *vo.Payload) (*vo.Payload, []error) {
				return p.projectorService.ProjectPayload(spec.Projector(), in, request, history)
			},
		},
		step[*vo.Payload]{
			label: "omit-empty",
			run: func(in *vo.Payload) (*vo.Payload, []error) {
				if spec.OmitEmpty() {
					return p.omitterService.OmitEmptyValuesFromPayload(in)
				}
				return in, nil
			},
		},
		step[*vo.Payload]{
			label: "aggregate",
			run: func(in *vo.Payload) (*vo.Payload, []error) {
				if spec.HasGroup() {
					out, err := p.aggregatorService.AggregatePayloadToKey(spec.Group(), in)
					return out, converter.ToSliceIfNonNil(err)
				}
				return in, nil
			},
		},
	)
}

func (p BuildPipeline) ApplyGroupID(
	spec vo.GroupIDPipelineSpec,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (string, []error) {
	if checker.IsNil(spec) || !spec.HasGroupID() {
		return "", nil
	}
	return apply(
		spec.GroupID(),
		step[string]{
			label: "get-dynamic-value",
			run: func(in string) (string, []error) {
				value, allErrs := p.dynamicValueService.Get(in, request, history)
				if checker.IsNotEmpty(allErrs) {
					allErrs = errors.JoinInheritAsSlicef(allErrs, ", ",
						"pipeline failed: op=build-group-id value=%s", in)
				}
				return value, allErrs
			},
		})
}

func (p BuildPipeline) ApplyDeduplicationID(
	spec vo.DeduplicationIDPipelineSpec,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (string, []error) {
	if checker.IsNil(spec) || !spec.HasDeduplicationID() {
		return "", nil
	}
	return apply(
		spec.DeduplicationID(),
		step[string]{
			label: "get-dynamic-value",
			run: func(in string) (string, []error) {
				value, allErrs := p.dynamicValueService.Get(in, request, history)
				if checker.IsNotEmpty(allErrs) {
					allErrs = errors.JoinInheritAsSlicef(allErrs, ", ",
						"pipeline failed: op=build-deduplicate-id value=%s", in)
				}
				return value, allErrs
			},
		})
}

func (p BuildPipeline) ApplyAttributes(
	spec vo.AttributesPipelineSpec,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (map[string]vo.AttributeValueConfig, []error) {
	if checker.IsNil(spec) {
		return nil, nil
	}
	return apply(
		spec.Attributes(),
		step[map[string]vo.AttributeValueConfig]{
			label: "build-message-attribute",
			run: func(m map[string]vo.AttributeValueConfig) (map[string]vo.AttributeValueConfig, []error) {
				attributes := make(map[string]vo.AttributeValueConfig, len(spec.Attributes()))
				allErrs := make([]error, 0)

				for key, attribute := range spec.Attributes() {
					value, errs := p.dynamicValueService.Get(attribute.Value(), request, history)
					if checker.IsNotEmpty(errs) {
						allErrs = append(allErrs, errors.JoinInheritf(errs, ", ",
							"pipeline failed: op=build-message-attribute key=%s value=%s", key, attribute.Value()))
						continue
					}
					attributes[key] = vo.NewAttributeValueConfig(attribute.Type(), value)
				}

				return attributes, allErrs
			},
		},
	)
}

type step[T any] struct {
	label string
	run   func(T) (T, []error)
}

func apply[T any](
	in T,
	steps ...step[T],
) (T, []error) {
	cur := in

	var allErrs []error
	for _, s := range steps {
		out, errs := s.run(cur)
		if checker.IsNotEmpty(errs) {
			allErrs = append(allErrs, errs...)
		}
		cur = out
	}

	return cur, allErrs
}
