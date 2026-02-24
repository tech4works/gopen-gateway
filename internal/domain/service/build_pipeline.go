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
	spec *vo.HostPipelineSpec,
) (string, []error) {
	return spec.Hosts()[rand.Intn(len(spec.Hosts())-1)], nil
}

func (p BuildPipeline) ApplyHeader(
	spec *vo.HeaderPipelineSpec,
	header vo.Header,
	request *vo.HTTPRequest,
	history *aggregate.History,
) (vo.Header, []error) {
	if checker.IsNil(spec) || spec.Omit() {
		return header, nil
	}
	return apply(
		header,
		step[vo.Header]{
			label: "aggregate",
			run: func(in vo.Header) (vo.Header, []error) {
				return p.aggregatorService.AggregateHeaders(in, request.Header()), nil
			},
		},
		step[vo.Header]{
			label: "modifiers",
			run: func(in vo.Header) (vo.Header, []error) {
				return p.modifierService.ExecuteHeaderModifiers(spec.Modifiers(), in, request, history)
			},
		},
		step[vo.Header]{
			label: "mapper",
			run: func(in vo.Header) (vo.Header, []error) {
				out, err := p.mapperService.MapHeader(spec.Mapper(), in, request, history)
				return out, converter.ToSliceIfNonNil(err)
			},
		},
		step[vo.Header]{
			label: "projector",
			run: func(in vo.Header) (vo.Header, []error) {
				out, err := p.projectorService.ProjectHeader(spec.Projector(), in, request, history)
				return out, converter.ToSliceIfNonNil(err)
			},
		},
	)
}

func (p BuildPipeline) ApplyURLPath(
	spec *vo.URLPathPipelineSpec,
	urlPath vo.URLPath,
	req *vo.HTTPRequest,
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
				return p.modifierService.ExecuteURLPathModifiers(spec.Modifiers(), in, req, hist)
			},
		},
	)
}

func (p BuildPipeline) ApplyQuery(
	spec *vo.QueryPipelineSpec,
	query vo.Query,
	request *vo.HTTPRequest,
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

func (p BuildPipeline) ApplyBody(
	spec *vo.BodyPipelineSpec,
	body *vo.Body,
	request *vo.HTTPRequest,
	history *aggregate.History,
) (*vo.Body, []error) {
	if checker.IsNil(spec) {
		return body, nil
	} else if spec.Omit() || checker.IsNil(body) {
		return nil, nil
	}
	return apply(
		body,
		step[*vo.Body]{
			label: "modifiers",
			run: func(in *vo.Body) (*vo.Body, []error) {
				return p.modifierService.ExecuteBodyModifiers(spec.Modifiers(), in, request, history)
			},
		},
		step[*vo.Body]{
			label: "joins",
			run: func(in *vo.Body) (*vo.Body, []error) {
				return p.joinService.ExecuteBodyJoins(spec.Joins(), in, request, history)
			},
		},
		step[*vo.Body]{
			label: "mapper",
			run: func(in *vo.Body) (*vo.Body, []error) {
				return p.mapperService.MapBody(spec.Mapper(), in, request, history)
			},
		},
		step[*vo.Body]{
			label: "projector",
			run: func(in *vo.Body) (*vo.Body, []error) {
				return p.projectorService.ProjectBody(spec.Projector(), in, request, history)
			},
		},
		step[*vo.Body]{
			label: "omit-empty",
			run: func(in *vo.Body) (*vo.Body, []error) {
				if spec.OmitEmpty() {
					return p.omitterService.OmitEmptyValuesFromBody(in)
				}
				return in, nil
			},
		},
		step[*vo.Body]{
			label: "aggregate",
			run: func(in *vo.Body) (*vo.Body, []error) {
				if spec.HasGroup() {
					out, err := p.aggregatorService.AggregateBodyToKey(spec.Group(), in)
					return out, converter.ToSliceIfNonNil(err)
				}
				return in, nil
			},
		},
	)
}

func (p BuildPipeline) ApplyGroupID(
	spec *vo.GroupIDPipelineSpec,
	request *vo.HTTPRequest,
	history *aggregate.History,
) (string, []error) {
	if checker.IsNil(spec) {
		return "", nil
	}
	return apply(
		spec.Value(),
		step[string]{
			label: "get-dynamic-value",
			run: func(in string) (string, []error) {
				value, allErrs := p.dynamicValueService.Get(spec.Value(), request, history)
				if checker.IsNotEmpty(allErrs) {
					allErrs = errors.JoinInheritAsSlicef(allErrs, ", ",
						"pipeline failed: op=build-group-id value=%s", spec.Value())
				}
				return value, allErrs
			},
		})
}

func (p BuildPipeline) ApplyDeduplicationID(
	spec *vo.DeduplicationIDPipelineSpec,
	request *vo.HTTPRequest,
	history *aggregate.History,
) (string, []error) {
	if checker.IsNil(spec) {
		return "", nil
	}
	return apply(
		spec.Value(),
		step[string]{
			label: "get-dynamic-value",
			run: func(in string) (string, []error) {
				value, allErrs := p.dynamicValueService.Get(spec.Value(), request, history)
				if checker.IsNotEmpty(allErrs) {
					allErrs = errors.JoinInheritAsSlicef(allErrs, ", ",
						"pipeline failed: op=build-deduplicate-id value=%s", spec.Value())
				}
				return value, allErrs
			},
		})
}

func (p BuildPipeline) ApplyPublisherAttributes(
	spec *vo.PublisherAttributesPipelineSpec,
	request *vo.HTTPRequest,
	history *aggregate.History,
) (map[string]vo.PublisherMessageAttribute, []error) {
	if checker.IsNil(spec) {
		return nil, nil
	}
	return apply(
		spec.Attributes(),
		step[map[string]vo.PublisherMessageAttribute]{
			label: "build-message-attribute",
			run: func(m map[string]vo.PublisherMessageAttribute) (map[string]vo.PublisherMessageAttribute, []error) {
				attributes := make(map[string]vo.PublisherMessageAttribute, len(spec.Attributes()))
				allErrs := make([]error, 0)

				for key, attribute := range spec.Attributes() {
					value, errs := p.dynamicValueService.Get(attribute.Value(), request, history)
					if checker.IsNotEmpty(errs) {
						allErrs = append(allErrs, errors.JoinInheritf(errs, ", ",
							"pipeline failed: op=build-message-attribute key=%s value=%s", key, attribute.Value()))
						continue
					}
					attributes[key] = vo.NewPublisherMessageAttribute(attribute.DataType(), value)
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
