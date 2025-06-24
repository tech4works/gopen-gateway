package factory

import (
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"github.com/tech4works/gopen-gateway/internal/domain/service"
)

type message struct {
	dynamicValueService service.DynamicValue
}

type Message interface {
	Build(request *vo.HTTPRequest, history *vo.History, publisher *vo.Publisher) (*vo.Message, []error)
}

func NewMessage(dynamicValueService service.DynamicValue) Message {
	return message{
		dynamicValueService: dynamicValueService,
	}
}

func (m message) Build(request *vo.HTTPRequest, history *vo.History, publisher *vo.Publisher) (*vo.Message, []error) {
	body, err := request.Body().String()
	if checker.NonNil(err) {
		return nil, []error{err}
	}

	groupID, groupErrs := m.dynamicValueService.Get(publisher.GroupID(), request, history)
	deduplicateID, deduplicateErrs := m.dynamicValueService.Get(publisher.DeduplicationID(), request, history)

	var allErrs []error
	allErrs = append(allErrs, groupErrs...)
	allErrs = append(allErrs, deduplicateErrs...)

	return vo.NewMessage(body, groupID, deduplicateID, publisher.Delay()), allErrs
}
