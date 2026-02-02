package factory

import (
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"github.com/tech4works/gopen-gateway/internal/domain/service"
)

type message struct {
	dynamicValueService service.DynamicValue
	mapperService       service.Mapper
	projectorService    service.Projector
	modifierService     service.Modifier
}

type Message interface {
	Build(request *vo.HTTPRequest, history *vo.History, publisher *vo.Publisher) (*vo.Message, []error)
}

func NewMessage(
	dynamicValueService service.DynamicValue,
	mapperService service.Mapper,
	projectorService service.Projector,
	modifierService service.Modifier,
) Message {
	return message{
		dynamicValueService: dynamicValueService,
		mapperService:       mapperService,
		projectorService:    projectorService,
		modifierService:     modifierService,
	}
}

func (m message) Build(request *vo.HTTPRequest, history *vo.History, publisher *vo.Publisher) (*vo.Message, []error) {
	groupID, groupErrs := m.dynamicValueService.Get(publisher.GroupID(), request, history)
	if checker.IsNotEmpty(groupErrs) {
		return nil, groupErrs
	}

	deduplicateID, deduplicateErrs := m.dynamicValueService.Get(publisher.DeduplicationID(), request, history)
	if checker.IsNotEmpty(deduplicateErrs) {
		return nil, deduplicateErrs
	}

	body := request.Body()

	body, modifyErrs := m.modifyBody(publisher.BodyModifiers(), body, request, history)
	if checker.IsNotEmpty(modifyErrs) {
		return nil, modifyErrs
	}

	body, mapErrs := m.mapperService.MapBody(body, publisher.BodyMapper())
	if checker.IsNotEmpty(mapErrs) {
		return nil, mapErrs
	}

	body, projectErrs := m.projectorService.ProjectBody(body, publisher.BodyProjection())
	if checker.IsNotEmpty(projectErrs) {
		return nil, projectErrs
	}

	bodyString, err := body.String()
	if checker.NonNil(err) {
		return nil, []error{err}
	}

	return vo.NewMessage(bodyString, groupID, deduplicateID, publisher.Delay()), nil
}

func (m message) modifyBody(modifiers []vo.Modifier, body *vo.Body, request *vo.HTTPRequest, history *vo.History) (
	*vo.Body, []error) {
	var errs []error

	for _, bodyModifier := range modifiers {
		modifierValue, dynamicValueErrs := m.dynamicValueService.Get(bodyModifier.Value(), request, history)
		if checker.IsNotEmpty(dynamicValueErrs) {
			errs = append(errs, dynamicValueErrs...)
		}
		modifiedBody, err := m.modifierService.ModifyBody(body, bodyModifier.Action(), bodyModifier.Key(), modifierValue)
		if checker.NonNil(err) {
			errs = append(errs, err)
		}
		body = modifiedBody
	}

	return body, errs
}
