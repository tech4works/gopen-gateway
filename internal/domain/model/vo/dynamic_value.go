package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/tidwall/gjson"
	"regexp"
	"strings"
)

type DynamicValue struct {
	value string
}

func NewDynamicValue(value string) DynamicValue {
	return DynamicValue{
		value: value,
	}
}

func (m DynamicValue) AsInt(httpRequest *HttpRequest, httpResponse *HttpResponse) int {
	// inicializamos o valor dinâmico
	dynamicValueStr := m.AsString(httpRequest, httpResponse)
	// retornamos o valor modificado ou não
	return helper.SimpleConvertToInt(dynamicValueStr)
}

func (m DynamicValue) AsSliceOfString(httpRequest *HttpRequest, httpResponse *HttpResponse) []string {
	// inicializamos o valor dinâmico
	dynamicValueStr := m.AsString(httpRequest, httpResponse)
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

func (m DynamicValue) AsString(httpRequest *HttpRequest, httpResponse *HttpResponse) string {
	// inicializamos o valor dinâmico como ele mesmo
	value := m.value
	// iteramos os valores com base na sintaxe do mesmo
	for _, word := range m.findAllByDynamicValueSyntax() {
		// processamos o palavra para converter em um valor dinâmico
		value = m.processDynamicValueWord(word, httpRequest, httpResponse)
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

func (m DynamicValue) processDynamicValueWord(word string, httpRequest *HttpRequest, httpResponse *HttpResponse) string {
	// obtemos o valor pela palavra
	dynamicValue := m.getDynamicValuePerWord(word, httpRequest, httpResponse)
	// caso o valor não encontrado, retornamos o próprio campo value
	if helper.IsEmpty(dynamicValue) {
		return m.value
	}
	// substituímos a palavra pelo eval
	return strings.Replace(m.value, word, dynamicValue, 1)
}

func (m DynamicValue) getDynamicValuePerWord(word string, httpRequest *HttpRequest, httpResponse *HttpResponse) string {
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
		return m.getHttpRequestValueByJsonPath(cleanSintaxe, httpRequest)
	} else if helper.Contains(dotSplit[0], "response") {
		return m.getHttpResponseValueByJsonPath(cleanSintaxe, httpResponse)
	}
	// se não esta no padrão, retornamos vazio
	return ""
}

func (m DynamicValue) getHttpRequestValueByJsonPath(jsonPath string, httpRequest *HttpRequest) string {
	jsonPath = strings.Replace(jsonPath, "request.", "", 1)
	result := gjson.Get(httpRequest.Map(), jsonPath)
	return result.String()
}

func (m DynamicValue) getHttpResponseValueByJsonPath(jsonPath string, httpResponse *HttpResponse) string {
	jsonPath = strings.Replace(jsonPath, "response.", "", 1)
	result := gjson.Get(httpResponse.Map(), jsonPath)
	return result.String()
}
