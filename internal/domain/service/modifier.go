package service

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"strings"
)

type modifierService struct {
	jsonPath domain.JSONPath
}

type Modifier interface {
	ModifyUrlPath(urlPath vo.URLPath, action enum.ModifierAction, key, value string) (vo.URLPath, error)
	ModifyHeader(header vo.Header, action enum.ModifierAction, key string, value []string) (vo.Header, error)
	ModifyQuery(query vo.Query, action enum.ModifierAction, key string, values []string) (vo.Query, error)
	ModifyBody(body *vo.Body, action enum.ModifierAction, key string, value string) (*vo.Body, error)
}

func NewModifier(jsonPath domain.JSONPath) Modifier {
	return modifierService{
		jsonPath: jsonPath,
	}
}

func (s modifierService) ModifyUrlPath(urlPath vo.URLPath, action enum.ModifierAction, key, value string) (
	vo.URLPath, error) {
	switch action {
	case enum.ModifierActionSet:
		return s.setUrlPath(urlPath, key, value)
	case enum.ModifierActionRpl:
		return s.replaceUrlPath(urlPath, key, value)
	case enum.ModifierActionDel:
		return s.deleteUrlPath(urlPath, key)
	default:
		return urlPath, mapper.NewErrInvalidAction("params", action)
	}
}

func (s modifierService) ModifyHeader(header vo.Header, action enum.ModifierAction, key string, value []string) (
	vo.Header, error) {
	switch action {
	case enum.ModifierActionAdd:
		return s.addHeader(header, key, value)
	case enum.ModifierActionApd:
		return s.appendHeader(header, key, value)
	case enum.ModifierActionSet:
		return s.setHeader(header, key, value)
	case enum.ModifierActionRpl:
		return s.replaceHeader(header, key, value)
	case enum.ModifierActionDel:
		return s.deleteHeader(header, key)
	default:
		return header, mapper.NewErrInvalidAction("header", action)
	}
}

func (s modifierService) ModifyQuery(query vo.Query, action enum.ModifierAction, key string, values []string) (
	vo.Query, error) {
	switch action {
	case enum.ModifierActionAdd:
		return s.addQuery(query, key, values)
	case enum.ModifierActionApd:
		return s.appendQuery(query, key, values)
	case enum.ModifierActionSet:
		return s.setQuery(query, key, values)
	case enum.ModifierActionRpl:
		return s.replaceQuery(query, key, values)
	case enum.ModifierActionDel:
		return s.deleteQuery(query, key)
	default:
		return query, mapper.NewErrInvalidAction("query", action)
	}
}

func (s modifierService) ModifyBody(body *vo.Body, action enum.ModifierAction, key string, value string) (*vo.Body, error) {
	if helper.IsNil(body) {
		return nil, nil
	}

	switch action {
	case enum.ModifierActionAdd:
		return s.addBody(body, key, value)
	case enum.ModifierActionApd:
		return s.appendBody(body, key, value)
	case enum.ModifierActionSet:
		return s.setBody(body, key, value)
	case enum.ModifierActionRpl:
		return s.replaceBody(body, key, value)
	case enum.ModifierActionDel:
		return s.deleteBody(body, key)
	default:
		return body, mapper.NewErrInvalidAction("body", action)
	}
}

func (s modifierService) validateKey(key string) error {
	if helper.IsEmpty(key) {
		return mapper.NewErrEmptyKey()
	}
	return nil
}

func (s modifierService) validateValue(value any) error {
	if helper.IsEmpty(value) {
		return mapper.NewErrEmptyValue()
	}
	return nil
}

func (s modifierService) setUrlPath(urlPath vo.URLPath, key, value string) (vo.URLPath, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return urlPath, err
	} else if err = s.validateValue(value); helper.IsNotNil(err) {
		return urlPath, err
	}

	path := urlPath.Raw()
	paramValues := urlPath.Params().Copy()

	paramValues[key] = value
	if helper.NotContains(path, fmt.Sprintf("/:%s", key)) {
		path = fmt.Sprintf("%s/:%s", path, key)
	}

	return vo.NewURLPath(path, paramValues), nil
}

func (s modifierService) replaceUrlPath(urlPath vo.URLPath, key, value string) (vo.URLPath, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return urlPath, err
	} else if err = s.validateValue(value); helper.IsNotNil(err) {
		return urlPath, err
	} else if urlPath.NotExists(key) {
		return urlPath, nil
	}

	return s.setUrlPath(urlPath, key, value)
}

func (s modifierService) deleteUrlPath(urlPath vo.URLPath, key string) (vo.URLPath, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return urlPath, err
	}

	path := strings.ReplaceAll(urlPath.Raw(), fmt.Sprintf("/:%s", key), "")

	paramValues := urlPath.Params().Copy()
	delete(paramValues, key)

	return vo.NewURLPath(path, paramValues), nil
}

func (s modifierService) addHeader(header vo.Header, key string, value []string) (vo.Header, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return header, err
	} else if err = s.validateValue(value); helper.IsNotNil(err) {
		return header, err
	} else if mapper.IsHeaderMandatoryKey(key) {
		return header, nil
	}

	values := header.Copy()
	values[key] = append(header.GetAll(key), value...)

	return vo.NewHeader(values), nil
}

func (s modifierService) appendHeader(header vo.Header, key string, value []string) (vo.Header, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return header, err
	} else if err = s.validateValue(value); helper.IsNotNil(err) {
		return header, err
	} else if mapper.IsHeaderMandatoryKey(key) || header.NotExists(key) {
		return header, nil
	}

	values := header.Copy()
	values[key] = append(header.GetAll(key), value...)

	return vo.NewHeader(values), nil
}

func (s modifierService) setHeader(header vo.Header, key string, value []string) (vo.Header, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return header, err
	} else if err = s.validateValue(value); helper.IsNotNil(err) {
		return header, err
	} else if mapper.IsHeaderMandatoryKey(key) {
		return header, nil
	}

	values := header.Copy()
	values[key] = value

	return vo.NewHeader(values), nil
}

func (s modifierService) replaceHeader(header vo.Header, key string, value []string) (vo.Header, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return header, err
	} else if err = s.validateValue(value); helper.IsNotNil(err) {
		return header, err
	} else if mapper.IsHeaderMandatoryKey(key) || header.NotExists(key) {
		return header, nil
	}

	return s.setHeader(header, key, value)
}

func (s modifierService) deleteHeader(header vo.Header, key string) (vo.Header, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return header, err
	} else if mapper.IsHeaderMandatoryKey(key) {
		return header, nil
	}

	values := header.Copy()
	delete(values, key)

	return vo.NewHeader(values), nil
}

func (s modifierService) addQuery(query vo.Query, key string, value []string) (vo.Query, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return query, err
	} else if err = s.validateValue(value); helper.IsNotNil(err) {
		return query, err
	}

	values := query.Copy()
	values[key] = append(query.GetAll(key), value...)

	return vo.NewQuery(values), nil
}

func (s modifierService) appendQuery(query vo.Query, key string, value []string) (vo.Query, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return query, err
	} else if err = s.validateValue(value); helper.IsNotNil(err) {
		return query, err
	} else if query.NotExists(key) {
		return query, nil
	}

	values := query.Copy()
	values[key] = append(query.GetAll(key), value...)

	return vo.NewQuery(values), nil
}

func (s modifierService) setQuery(query vo.Query, key string, value []string) (vo.Query, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return query, err
	} else if err = s.validateValue(value); helper.IsNotNil(err) {
		return query, err
	}

	values := query.Copy()
	values[key] = value

	return vo.NewQuery(values), nil
}

func (s modifierService) replaceQuery(query vo.Query, key string, value []string) (vo.Query, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return query, err
	} else if err = s.validateValue(value); helper.IsNotNil(err) {
		return query, err
	} else if query.NotExists(key) {
		return query, nil
	}

	return s.setQuery(query, key, value)
}

func (s modifierService) deleteQuery(query vo.Query, key string) (vo.Query, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return query, err
	}

	values := query.Copy()
	delete(values, key)

	return vo.NewQuery(values), nil
}

func (s modifierService) addBody(body *vo.Body, key, value string) (*vo.Body, error) {
	if body.ContentType().IsText() {
		return s.addBodyText(body, value)
	} else if body.ContentType().IsJSON() {
		return s.addBodyJson(body, key, value)
	}
	return body, mapper.NewErrIncompatibleBodyType(body.ContentType())
}

func (s modifierService) addBodyText(body *vo.Body, value string) (*vo.Body, error) {
	if err := s.validateValue(value); helper.IsNotNil(err) {
		return body, err
	}

	bodyStr, err := body.String()
	if helper.IsNotNil(err) {
		return body, err
	}

	modifiedBodyText := fmt.Sprintf("%s%s", bodyStr, value)
	return s.newBodyByString(body, modifiedBodyText)
}

func (s modifierService) addBodyJson(body *vo.Body, key, value string) (*vo.Body, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return body, err
	} else if err = s.validateValue(value); helper.IsNotNil(err) {
		return body, err
	}

	bodyRaw, err := body.Raw()
	if helper.IsNotNil(err) {
		return body, err
	}

	modifiedBodyJson, err := s.jsonPath.Add(bodyRaw, key, value)
	if helper.IsNotNil(err) {
		return body, err
	}

	return s.newBodyByString(body, modifiedBodyJson)
}

func (s modifierService) appendBody(body *vo.Body, key, value string) (*vo.Body, error) {
	if body.ContentType().IsText() {
		return s.appendBodyText(body, value)
	} else if body.ContentType().IsJSON() {
		return s.appendBodyJson(body, key, value)
	}
	return body, mapper.NewErrIncompatibleBodyType(body.ContentType())
}

func (s modifierService) appendBodyText(body *vo.Body, value string) (*vo.Body, error) {
	if err := s.validateValue(value); helper.IsNotNil(err) {
		return body, err
	}

	bodyStr, err := body.String()
	if helper.IsNotNil(err) {
		return body, err
	}

	modifiedBodyText := fmt.Sprintf("%s\n%s", bodyStr, value)
	return s.newBodyByString(body, modifiedBodyText)
}

func (s modifierService) appendBodyJson(body *vo.Body, key, value string) (*vo.Body, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return body, err
	} else if err = s.validateValue(value); helper.IsNotNil(err) {
		return body, err
	}

	bodyRaw, err := body.Raw()
	if helper.IsNotNil(err) {
		return body, err
	}

	if s.jsonPath.Parse(bodyRaw).Get(key).NotExists() {
		return body, nil
	}

	modifiedBodyJson, err := s.jsonPath.Add(bodyRaw, key, value)
	if helper.IsNotNil(err) {
		return body, err
	}

	return s.newBodyByString(body, modifiedBodyJson)
}

func (s modifierService) setBody(body *vo.Body, key, value string) (*vo.Body, error) {
	if body.ContentType().IsText() {
		return s.setBodyText(body, value)
	} else if body.ContentType().IsJSON() {
		return s.setBodyJson(body, key, value)
	}
	return body, mapper.NewErrIncompatibleBodyType(body.ContentType())
}

func (s modifierService) setBodyText(body *vo.Body, value string) (*vo.Body, error) {
	if err := s.validateValue(value); helper.IsNotNil(err) {
		return body, err
	}

	return s.newBodyByString(body, value)
}

func (s modifierService) setBodyJson(body *vo.Body, key, value string) (*vo.Body, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return body, err
	} else if err = s.validateValue(value); helper.IsNotNil(err) {
		return body, err
	}

	bodyRaw, err := body.Raw()
	if helper.IsNotNil(err) {
		return body, err
	}

	modifiedBodyJson, err := s.jsonPath.Set(bodyRaw, key, value)
	if helper.IsNotNil(err) {
		return body, err
	}

	return s.newBodyByString(body, modifiedBodyJson)
}

func (s modifierService) replaceBody(body *vo.Body, key, value string) (*vo.Body, error) {
	if body.ContentType().IsText() {
		return s.replaceBodyText(body, key, value)
	} else if body.ContentType().IsJSON() {
		return s.replaceBodyJson(body, key, value)
	}
	return body, mapper.NewErrIncompatibleBodyType(body.ContentType())
}

func (s modifierService) replaceBodyText(body *vo.Body, key, value string) (*vo.Body, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return body, err
	} else if err = s.validateValue(value); helper.IsNotNil(err) {
		return body, err
	}

	bodyStr, err := body.String()
	if helper.IsNotNil(err) {
		return body, err
	}

	modifiedBodyText := strings.ReplaceAll(bodyStr, key, value)
	return s.newBodyByString(body, modifiedBodyText)
}

func (s modifierService) replaceBodyJson(body *vo.Body, key, value string) (*vo.Body, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return body, err
	} else if err = s.validateValue(value); helper.IsNotNil(err) {
		return body, err
	}

	bodyRaw, err := body.Raw()
	if helper.IsNotNil(err) {
		return body, err
	}

	if s.jsonPath.Parse(bodyRaw).Get(key).NotExists() {
		return body, nil
	}

	modifiedBodyJson, err := s.jsonPath.Set(bodyRaw, key, value)
	if helper.IsNotNil(err) {
		return body, err
	}

	return s.newBodyByString(body, modifiedBodyJson)
}

func (s modifierService) deleteBody(body *vo.Body, key string) (*vo.Body, error) {
	if body.ContentType().IsText() {
		return s.deleteBodyText(body, key)
	} else if body.ContentType().IsJSON() {
		return s.deleteBodyJson(body, key)
	}
	return body, mapper.NewErrIncompatibleBodyType(body.ContentType())
}

func (s modifierService) deleteBodyText(body *vo.Body, key string) (*vo.Body, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return body, err
	}

	bodyStr, err := body.String()
	if helper.IsNotNil(err) {
		return body, err
	}

	modifiedBodyText := strings.ReplaceAll(bodyStr, key, "")
	return s.newBodyByString(body, modifiedBodyText)
}

func (s modifierService) deleteBodyJson(body *vo.Body, key string) (*vo.Body, error) {
	if err := s.validateKey(key); helper.IsNotNil(err) {
		return body, err
	}

	bodyRaw, err := body.Raw()
	if helper.IsNotNil(err) {
		return body, err
	}

	modifiedBodyJson, err := s.jsonPath.Delete(bodyRaw, key)
	if helper.IsNotNil(err) {
		return body, err
	}

	return s.newBodyByString(body, modifiedBodyJson)
}

func (s modifierService) newBodyByString(body *vo.Body, modifiedBodyJson string) (*vo.Body, error) {
	buffer, err := helper.ConvertToBuffer(modifiedBodyJson)
	if helper.IsNotNil(err) {
		return body, err
	}

	return vo.NewBodyWithContentType(body.ContentType(), buffer), nil
}
