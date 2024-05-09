package vo

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"strings"
)

type Params = map[string]string

// UrlPath represents a string type that represents a URL path.
type UrlPath struct {
	value  string
	params Params
}

func NewUrlPath(value string, params Params) UrlPath {
	// filtramos os parâmetros que contem no valor passado
	filteredParams := Params{}
	for k, v := range params {
		paramPathKey := patternParamPathKey(k)
		if helper.Contains(value, paramPathKey) {
			filteredParams[k] = v
		}
	}
	return UrlPath{
		value:  value,
		params: filteredParams,
	}
}

func (u UrlPath) Modify(modifier *Modifier, httpRequest *HttpRequest, httpResponse *HttpResponse) UrlPath {
	// instanciamos o novo valor a ser utilizado ao modificar
	newValue := modifier.ValueAsString(httpRequest, httpResponse)
	// alteramos conforme a ação mapeada
	switch modifier.Action() {
	case enum.ModifierActionSet:
		return u.SetParam(modifier.Key(), newValue)
	case enum.ModifierActionRpl:
		return u.ReplaceParam(modifier.Key(), newValue)
	case enum.ModifierActionDel:
		return u.DeleteParam(modifier.Key())
	default:
		return u
	}
}

func (u UrlPath) Params() Params {
	return u.params
}

func (u UrlPath) String() string {
	// instanciamos o valor a ser retornado
	urlPathString := u.value
	// preenchemos todos os parâmetros do value e retornamos
	for k, v := range u.params {
		paramPathKey := patternParamPathKey(k)
		urlPathString = strings.ReplaceAll(urlPathString, paramPathKey, v)
	}
	return urlPathString
}

func (u UrlPath) SetParam(key, value string) UrlPath {
	// copiamos os parâmetros atuais
	urlPathParams := u.copyParams()
	// inserimos a nova chave e valor
	urlPathParams[key] = value

	// instanciamos o valor atual do url path
	urlPathValue := u.value
	// verificamos se o parâmetro contain no valor do urlPath
	paramPathKey := patternParamPathKey(key)
	if helper.NotContains(urlPathValue, paramPathKey) {
		urlPathValue = fmt.Sprintf("%s/%s", urlPathValue, paramPathKey)
	}
	// retornamos os valores modificados em um novo VO
	return UrlPath{
		value:  urlPathValue,
		params: urlPathParams,
	}
}

func (u UrlPath) ReplaceParam(key, value string) UrlPath {
	// caso não existe o parâmetro retornamos a UrlPath atual
	if u.NotExistsParam(key) {
		return u
	}
	// se existe alteramos
	return u.SetParam(key, value)
}

func (u UrlPath) DeleteParam(key string) UrlPath {
	// copiamos os parâmetros atuais
	urlPathParams := u.copyParams()
	// removemos a chave indicada
	delete(urlPathParams, key)
	// removemos a chave do valor da urlPath caso informada
	urlPathValue := strings.ReplaceAll(u.value, patternParamUrlKey(key), "")
	// retornamos os valores modificados em um novo VO
	return UrlPath{
		value:  urlPathValue,
		params: urlPathParams,
	}
}

func (u UrlPath) ExistsParam(key string) bool {
	_, ok := u.params[key]
	return ok && helper.Contains(u.value, patternParamPathKey(key))
}

func (u UrlPath) NotExistsParam(key string) bool {
	return !u.ExistsParam(key)
}

func (u UrlPath) copyParams() map[string]string {
	copiedParams := make(map[string]string)
	for k, v := range u.params {
		copiedParams[k] = v
	}
	return copiedParams
}

func patternParamPathKey(k string) string {
	return fmt.Sprintf(":%s", k)
}

func patternParamUrlKey(k string) string {
	return fmt.Sprintf("/:%s", k)
}
