package service

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type omitterService struct {
	jsonPath domain.JSONPath
}

type Omitter interface {
	OmitEmptyValuesFromBody(body *vo.Body) (*vo.Body, []error)
}

func NewOmitter(jsonPath domain.JSONPath) Omitter {
	return omitterService{
		jsonPath: jsonPath,
	}
}

func (o omitterService) OmitEmptyValuesFromBody(body *vo.Body) (*vo.Body, []error) {
	if body.ContentType().IsText() {
		return o.omitEmptyValuesFromBodyText(body)
	} else if body.ContentType().IsJSON() {
		return o.omitEmptyValuesFromBodyJson(body)
	}
	return body, nil
}

func (o omitterService) omitEmptyValuesFromBodyText(body *vo.Body) (*vo.Body, []error) {
	bodyStr, err := body.String()
	if checker.NonNil(err) {
		return nil, []error{err}
	}

	buffer, err := helper.ConvertToBuffer(helper.CleanAllRepeatSpaces(bodyStr))
	if checker.NonNil(err) {
		return nil, []error{err}
	}

	return vo.NewBodyWithContentType(body.ContentType(), buffer), nil
}

func (o omitterService) omitEmptyValuesFromBodyJson(body *vo.Body) (*vo.Body, []error) {
	raw, err := body.Raw()
	if checker.NonNil(err) {
		return nil, []error{err}
	}

	newBodyStr, errs := o.removeAllEmptyFields(raw)
	if checker.IsNotEmpty(errs) {
		return nil, errs
	}

	buffer, err := helper.ConvertToBuffer(newBodyStr)
	if checker.NonNil(err) {
		return nil, []error{err}
	}

	return vo.NewBodyWithContentType(body.ContentType(), buffer), nil
}

func (o omitterService) removeAllEmptyFields(jsonStr string) (string, []error) {
	var errs []error

	o.jsonPath.Parse(jsonStr).ForEach(func(key string, value domain.JSONValue) bool {
		if value.IsObject() || value.IsArray() {
			childJson, childErrs := o.removeAllEmptyFields(value.Raw())
			if checker.IsNotEmpty(childErrs) {
				errs = append(errs, childErrs...)
				return true
			}
			value = o.jsonPath.Parse(childJson)
		}

		var newJsonStr string
		var err error
		if checker.IsEmpty(value.Interface()) {
			newJsonStr, err = o.jsonPath.Delete(jsonStr, key)
		} else {
			newJsonStr, err = o.jsonPath.Set(jsonStr, key, value.Raw())
		}

		if checker.NonNil(err) {
			errs = append(errs, err)
			return true
		}
		jsonStr = newJsonStr

		return true
	})

	return jsonStr, errs
}
