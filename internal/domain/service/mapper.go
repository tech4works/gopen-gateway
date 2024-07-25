package service

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain"
	mapper2 "github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"strings"
)

type mapperService struct {
	jsonPath domain.JSONPath
}

type Mapper interface {
	MapHeader(header vo.Header, mapper *vo.Mapper) vo.Header
	MapQuery(query vo.Query, mapper *vo.Mapper) vo.Query
	MapBody(body *vo.Body, mapper *vo.Mapper) (*vo.Body, []error)
}

func NewMapper(jsonPath domain.JSONPath) Mapper {
	return mapperService{
		jsonPath: jsonPath,
	}
}

func (m mapperService) MapHeader(header vo.Header, mapper *vo.Mapper) vo.Header {
	if helper.IsNil(mapper) || mapper.IsEmpty() {
		return header
	}
	mappedHeader := map[string][]string{}
	for _, key := range header.Keys() {
		if mapper2.IsNotHeaderMandatoryKey(key) && mapper.Exists(key) {
			mappedHeader[mapper.Get(key)] = header.GetAll(key)
		} else {
			mappedHeader[key] = header.GetAll(key)
		}
	}
	return vo.NewHeader(mappedHeader)
}

func (m mapperService) MapQuery(query vo.Query, mapper *vo.Mapper) vo.Query {
	if helper.IsNil(mapper) || mapper.IsEmpty() {
		return query
	}

	mappedQuery := map[string][]string{}
	for _, key := range query.Keys() {
		if mapper.Exists(key) {
			mappedQuery[mapper.Get(key)] = query.GetAll(key)
		} else {
			mappedQuery[key] = query.GetAll(key)
		}
	}
	return vo.NewQuery(mappedQuery)
}

func (m mapperService) MapBody(body *vo.Body, mapper *vo.Mapper) (*vo.Body, []error) {
	if helper.IsNil(mapper) || mapper.IsNotEmpty() || helper.IsNil(body) {
		return body, nil
	}

	if body.ContentType().IsText() {
		return m.mapBodyText(body, mapper)
	} else if body.ContentType().IsJSON() {
		return m.mapBodyJson(body, mapper)
	} else {
		return body, nil
	}
}

func (m mapperService) mapBodyText(body *vo.Body, mapper *vo.Mapper) (*vo.Body, []error) {
	mappedBody, err := body.String()
	if helper.IsNotNil(err) {
		return body, []error{err}
	}

	for _, key := range mapper.Keys() {
		newKey := mapper.Get(key)
		if helper.IsNotEqualTo(key, newKey) {
			mappedBody = strings.ReplaceAll(mappedBody, key, newKey)
		}
	}

	buffer, err := helper.ConvertToBuffer(mappedBody)
	if helper.IsNotNil(err) {
		return body, []error{err}
	}

	return vo.NewBodyWithContentType(body.ContentType(), buffer), nil
}

func (m mapperService) mapBodyJson(body *vo.Body, mapper *vo.Mapper) (*vo.Body, []error) {
	bodyStr, err := body.String()
	if helper.IsNotNil(err) {
		return body, []error{err}
	}

	var mappedBodyStr string
	var errs []error

	parsedJson := m.jsonPath.Parse(bodyStr)
	if parsedJson.IsArray() {
		mappedBodyStr, errs = m.mapBodyJsonArray(parsedJson, mapper)
	} else {
		mappedBodyStr, errs = m.mapBodyJsonObject(parsedJson, mapper)
	}
	if helper.IsNotEmpty(errs) {
		return body, errs
	}

	buffer, err := helper.ConvertToBuffer(mappedBodyStr)
	if helper.IsNotNil(err) {
		return body, []error{err}
	}

	return vo.NewBodyWithContentType(body.ContentType(), buffer), nil
}

func (m mapperService) mapBodyJsonArray(jsonArray domain.JSONValue, mapper *vo.Mapper) (string, []error) {
	var mappedArray = "[]"
	var errs []error

	jsonArray.ForEach(func(key string, value domain.JSONValue) bool {
		var newMappedArray string
		var err error
		if value.IsObject() {
			childObject, childErrs := m.mapBodyJsonObject(value, mapper)
			if helper.IsNotEmpty(childErrs) {
				errs = append(errs, childErrs...)
				return true
			}
			newMappedArray, err = m.jsonPath.AppendOnArray(mappedArray, childObject)
		} else if value.IsArray() {
			childArray, childErrs := m.mapBodyJsonArray(value, mapper)
			if helper.IsNotEmpty(childErrs) {
				errs = append(errs, childErrs...)
				return true
			}
			newMappedArray, err = m.jsonPath.AppendOnArray(mappedArray, childArray)
		} else {
			newMappedArray, err = m.jsonPath.AppendOnArray(mappedArray, value.Raw())
		}

		if helper.IsNotNil(err) {
			errs = append(errs, err)
			return true
		}

		mappedArray = newMappedArray
		return true
	})

	return mappedArray, errs
}

func (m mapperService) mapBodyJsonObject(jsonObject domain.JSONValue, mapper *vo.Mapper) (string, []error) {
	var mappedJson = jsonObject.Raw()
	var errs []error

	for _, key := range mapper.Keys() {
		newKey := mapper.Get(key)
		if helper.Equals(key, newKey) {
			continue
		}
		jsonValue := jsonObject.Get(key)
		if jsonValue.NotExists() {
			continue
		}

		newMappedJson, err := m.jsonPath.Set(mappedJson, newKey, jsonValue.Raw())
		if helper.IsNotNil(err) {
			errs = append(errs, err)
			continue
		}

		newMappedJson, err = m.jsonPath.Delete(newMappedJson, key)
		if helper.IsNotNil(err) {
			errs = append(errs, err)
			continue
		}

		mappedJson = newMappedJson
	}

	return mappedJson, errs
}
