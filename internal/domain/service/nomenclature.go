package service

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
)

type nomenclatureService struct {
	jsonPath     domain.JSONPath
	nomenclature domain.Nomenclature
}

type Nomenclature interface {
	ToCase(body *vo.Body, nomenclature enum.Nomenclature) (*vo.Body, []error)
}

func NewNomenclature(jsonPath domain.JSONPath) Nomenclature {
	return nomenclatureService{
		jsonPath: jsonPath,
	}
}

func (n nomenclatureService) ToCase(body *vo.Body, nomenclature enum.Nomenclature) (*vo.Body, []error) {
	if body.ContentType().IsNotJSON() {
		return body, nil
	}

	raw, err := body.Raw()
	if helper.IsNotNil(err) {
		return body, []error{err}
	}

	jsonStr, errs := n.convertKeysToCase(raw, nomenclature)
	if helper.IsNotEmpty(errs) {
		return body, errs
	}

	buffer, err := helper.ConvertToBuffer(jsonStr)
	if helper.IsNotNil(err) {
		return body, []error{err}
	}

	return vo.NewBodyWithContentType(body.ContentType(), buffer), nil
}

func (n nomenclatureService) convertKeysToCase(jsonStr string, nomenclature enum.Nomenclature) (string, []error) {
	jsonValue := n.jsonPath.Parse(jsonStr)

	var result = "{}"
	if jsonValue.IsArray() {
		result = "[]"
	}
	var errs []error

	jsonValue.ForEach(func(key string, value domain.JSONValue) bool {
		var newResult string
		var err error

		newKey := n.nomenclature.Parse(nomenclature, key)
		if value.IsObject() || value.IsArray() {
			childJson, childErrs := n.convertKeysToCase(value.Raw(), nomenclature)
			if helper.IsNotEmpty(childErrs) {
				errs = append(errs, childErrs...)
				return true
			}
			newResult, err = n.jsonPath.Set(result, newKey, childJson)
		} else {
			newResult, err = n.jsonPath.Set(result, newKey, value.Raw())
		}

		if helper.IsNotNil(err) {
			errs = append(errs, err)
			return true
		}

		result = newResult
		return true
	})

	return result, errs
}
