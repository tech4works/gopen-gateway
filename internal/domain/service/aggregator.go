package service

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
)

type aggregatorService struct {
	jsonPath domain.JSONPath
}

type Aggregator interface {
	AggregateHeaders(base, value vo.Header) vo.Header
	AggregateBodyToKey(key string, value *vo.Body) (*vo.Body, error)
	AggregateBodiesIntoSlice(history *vo.History) (*vo.Body, []error)
	AggregateBodies(history *vo.History) (*vo.Body, []error)
}

func NewAggregator(jsonPath domain.JSONPath) Aggregator {
	return aggregatorService{
		jsonPath: jsonPath,
	}
}

func (a aggregatorService) AggregateHeaders(base, value vo.Header) vo.Header {
	aggregated := base.Copy()
	for _, key := range value.Keys() {
		if mapper.IsNotHeaderMandatoryKey(key) {
			aggregated[key] = append(aggregated[key], value.GetAll(key)...)
		}
	}
	return vo.NewHeader(aggregated)
}

func (a aggregatorService) AggregateBodyToKey(key string, value *vo.Body) (*vo.Body, error) {
	if helper.IsEmpty(key) || value.ContentType().IsNotJSON() || helper.IsNil(value) {
		return value, nil
	}

	raw, err := value.Raw()
	if helper.IsNotNil(err) {
		return value, err
	}

	jsonValue, err := a.jsonPath.Set("{}", key, raw)
	if helper.IsNotNil(err) {
		return value, err
	}

	buffer, err := helper.ConvertToBuffer(jsonValue)
	if helper.IsNotNil(err) {
		return value, err
	}

	return vo.NewBodyWithContentType(vo.NewContentTypeJson(), buffer), nil
}

func (a aggregatorService) AggregateBodiesIntoSlice(history *vo.History) (*vo.Body, []error) {
	result := "[]"

	var errs []error
	for i := 0; i < history.Size(); i++ {
		_, httpBackendResponse := history.Get(i)

		newJsonStr := a.buildBodyDefaultForSlice(httpBackendResponse)

		if !httpBackendResponse.HasBody() {
			result, _ = a.jsonPath.AppendOnArray(result, newJsonStr)
			continue
		}

		raw, err := httpBackendResponse.Body().Raw()
		if helper.IsNotNil(err) {
			errs = append(errs, err)
			continue
		}

		newJsonStr, mergeErrs := a.merge(i, newJsonStr, raw)
		if helper.IsNotEmpty(mergeErrs) {
			errs = append(errs, mergeErrs...)
			continue
		}

		result, _ = a.jsonPath.AppendOnArray(result, newJsonStr)
	}

	return a.buildBodyJson(result, errs)
}

func (a aggregatorService) AggregateBodies(history *vo.History) (*vo.Body, []error) {
	result := "{}"

	var errs []error
	for i := 0; i < history.Size(); i++ {
		_, httpBackendResponse := history.Get(i)
		if !httpBackendResponse.HasBody() {
			continue
		}

		raw, err := httpBackendResponse.Body().Raw()
		if helper.IsNotNil(err) {
			errs = append(errs, err)
			continue
		}

		newJsonStr, mergeErrs := a.merge(i, result, raw)
		if helper.IsNotEmpty(mergeErrs) {
			errs = append(errs, mergeErrs...)
			continue
		}

		result = newJsonStr
	}

	return a.buildBodyJson(result, errs)
}

func (a aggregatorService) buildBodyDefaultForSlice(httpBackendResponse *vo.HTTPBackendResponse) string {
	code := httpBackendResponse.StatusCode()

	jsonStr := "{}"
	jsonStr, _ = a.jsonPath.Set(jsonStr, "ok", helper.SimpleConvertToString(httpBackendResponse.OK()))
	jsonStr, _ = a.jsonPath.Set(jsonStr, "code", code.String())

	return jsonStr
}

func (a aggregatorService) merge(i int, jsonStr, raw string) (string, []error) {
	if helper.IsNotJson(raw) || helper.IsSlice(raw) {
		return a.mergeJSONByKey(i, jsonStr, raw)
	}
	return a.mergeJSON(jsonStr, raw)
}

func (a aggregatorService) mergeJSONByKey(i int, jsonStr, raw string) (string, []error) {
	newJsonStr, err := a.jsonPath.Set(jsonStr, fmt.Sprintf("backend%v", i), raw)
	if helper.IsNotNil(err) {
		return jsonStr, []error{err}
	}
	return newJsonStr, nil
}

func (a aggregatorService) mergeJSON(jsonStr, raw string) (string, []error) {
	var result string
	var errs []error

	result = jsonStr
	a.jsonPath.Parse(raw).ForEach(func(key string, value domain.JSONValue) bool {
		newResult, err := a.jsonPath.Add(result, key, value.Raw())
		if helper.IsNotNil(err) {
			errs = append(errs, err)
			return true
		}
		result = newResult
		return true
	})

	return result, errs
}

func (a aggregatorService) buildBodyJson(result string, errs []error) (*vo.Body, []error) {
	buffer, err := helper.ConvertToBuffer(result)
	if helper.IsNotNil(err) {
		errs = append(errs, err)
		return nil, errs
	}

	return vo.NewBodyWithContentType(vo.NewContentTypeJson(), buffer), errs
}
