package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/tidwall/gjson"
	"regexp"
	"strings"
)

type DynamicValue struct {
	value        string
	httpRequest  *HttpRequest
	httpResponse *HttpResponse
}

func NewDynamicValue(value string, httpRequest *HttpRequest, httpResponse *HttpResponse) DynamicValue {
	return DynamicValue{
		httpRequest:  httpRequest,
		httpResponse: httpResponse,
		value:        value,
	}
}

func (m DynamicValue) AsInt() int {
	// inicializamos o valor dinâmico
	dynamicValueStr := m.AsString()
	// retornamos o valor modificado ou não
	return helper.SimpleConvertToInt(dynamicValueStr)
}

func (m DynamicValue) AsSliceOfString() []string {
	// inicializamos o valor dinâmico
	dynamicValueStr := m.AsString()
	// verificamos se o mesmo é um slice
	if helper.IsSliceType(dynamicValueStr) {
		var ss []string
		helper.SimpleConvertToDest(dynamicValueStr, &ss)
		if helper.IsNotEmpty(ss) {
			return ss
		}
	}
	// caso não seja uma string de slice, forcamos o retorno apenas com o valor do mesmo
	return []string{dynamicValueStr}
}

func (m DynamicValue) AsString() string {
	// inicializamos o valor dinâmico como ele mesmo
	value := m.value
	// iteramos os valores com base na sintaxe do mesmo
	for _, word := range m.findAllByDynamicValueSyntax() {
		// processamos o palavra para converter em um valor dinâmico
		value = m.processDynamicValueWord(word)
	}
	// damos o parse do valor em string caso tenha um valor do tipo não string
	return value
}

func (m DynamicValue) findAllByDynamicValueSyntax() []string {
	// criamos o regex esperado para obter o valor dinâmico
	regex := regexp.MustCompile(`\B#[a-zA-Z0-9_.\-\[\]]+`)
	// buscamos todos os valores no campo value com essa sintaxe
	return regex.FindAllString(m.value, -1)
}

func (m DynamicValue) processDynamicValueWord(word string) string {
	// obtemos o valor pela palavra
	dynamicValue := m.getDynamicValuePerWord(word)
	// caso o valor não encontrado, retornamos o próprio campo value
	if helper.IsEmpty(dynamicValue) {
		return m.value
	}
	// substituímos a palavra pelo eval
	return strings.Replace(m.value, word, dynamicValue, 1)
}

func (m DynamicValue) getDynamicValuePerWord(word string) string {
	// limpamos a #
	cleanSintaxe := strings.ReplaceAll(word, "#", "")
	// damos o split pela pontuação
	dotSplit := strings.Split(cleanSintaxe, ".")
	// caso esteja vazio já retornamos
	if helper.IsEmpty(dotSplit) {
		return ""
	}
	// obtemos o valor dinâmico pelo httpRequest or httpResponse
	if helper.Contains(dotSplit[0], "request") {
		return m.getHttpRequestValueByJsonPath(cleanSintaxe)
	} else if helper.Contains(dotSplit[0], "response") {
		return m.getHttpResponseValueByJsonPath(cleanSintaxe)
	}
	// se não esta no padrão, retornamos vazio
	return ""
}

func (m DynamicValue) getHttpRequestValueByJsonPath(jsonPath string) string {
	jsonPath = strings.Replace(jsonPath, "request.", "", 1)
	result := gjson.Get(m.httpRequest.Map(), jsonPath)
	return result.String()
}

func (m DynamicValue) getHttpResponseValueByJsonPath(jsonPath string) string {
	jsonPath = strings.Replace(jsonPath, "response.", "", 1)
	result := gjson.Get(m.httpResponse.Map(), jsonPath)
	return result.String()
}
