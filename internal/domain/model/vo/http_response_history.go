package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"net/http"
)

// httpResponseHistory represents the history of backend responses.
// It is a slice of httpBackendResponse, which represents the responses from a backend service.
// The httpResponse history can be filtered and modified based on certain conditions.
// It also provides methods to retrieve information about the httpResponse history, such as size, success, and status code.
type httpResponseHistory []*HttpBackendResponse

// Size returns the number of elements in the httpResponseHistory.
func (r httpResponseHistory) Size() int {
	return len(r)
}

// Success returns a boolean value indicating whether all backend responses in the httpResponseHistory are successful.
// It iterates through each backend httpResponse in the httpResponseHistory.
// If any backend httpResponse is not successful (Ok() returns false), it returns false.
// Otherwise, it returns true, indicating that all backend responses are successful.
func (r httpResponseHistory) Success() bool {
	for _, httpBackendResponse := range r {
		if !httpBackendResponse.Ok() {
			return false
		}
	}
	return true
}

func (r httpResponseHistory) Filter(httpRequest *HttpRequest, httpResponse *HttpResponse) (
	filteredHistory httpResponseHistory) {
	// iteramos o histórico para ser filtrado
	for _, httpBackendResponse := range r {
		// aplicamos a config de forma LATE
		httpBackendResponseApplied := httpBackendResponse.ApplyConfig(enum.BackendResponseApplyLate, httpRequest,
			httpResponse)
		// setamos a resposta caso não esteja nil
		if helper.IsNotNil(httpBackendResponseApplied) {
			filteredHistory = append(filteredHistory, httpBackendResponseApplied)
		}
	}
	return filteredHistory
}

// SingleResponse checks if the httpResponse history contains only one httpResponse.
// Returns true if the httpResponse history size is less than or equal to 1, false otherwise.
func (r httpResponseHistory) SingleResponse() bool {
	return helper.Equals(r.Size(), 1)
}

// MultipleResponse returns true if the size of the httpResponse history is greater than 1, indicating multiple responses.
// Otherwise, it returns false.
func (r httpResponseHistory) MultipleResponse() bool {
	return helper.IsGreaterThan(r.Size(), 1)
}

func (r httpResponseHistory) StatusCode() int {
	// se tiver mais de 1 resposta obtemos o código de status mais frequente
	if r.MultipleResponse() {
		return r.mostFrequentStatusCode()
	}
	return r.statusCode()
}

func (r httpResponseHistory) Header() Header {
	// instanciamos um novo header
	historyHeader := Header{}
	// iteramos as respostas do histórico
	for _, httpBackendResponse := range r {
		// se tiver nil ou vazio, ignoramos
		if helper.IsNil(httpBackendResponse.Header()) || helper.IsEmpty(httpBackendResponse.Header()) {
			continue
		}
		// agregamos os valores ao header de resultado, gerando um novo objeto de valor a cada agregação
		historyHeader = historyHeader.Aggregate(httpBackendResponse.Header())
	}
	// retornamos o header do histórico
	return historyHeader
}

func (r httpResponseHistory) Body(aggregate bool) *Body {
	// verificamos se o histórico é de múltiplas respostas
	if r.MultipleResponse() {
		// caso seja de múltiplas respostas, verificamos se precisa agregar as respostas
		return r.multipleResponseBody(aggregate)
	}
	// caso seja uma única ou não chamamos o body()
	return r.body()
}

func (r httpResponseHistory) Map() string {
	var mappedHistory []any
	for _, httpBackendResponse := range r {
		mappedHistory = append(mappedHistory, httpBackendResponse.Map())
	}
	return helper.SimpleConvertToString(mappedHistory)
}

func (r httpResponseHistory) last() *HttpBackendResponse {
	return r[len(r)-1]
}

func (r httpResponseHistory) statusCode() int {
	if helper.IsNil(r.last()) {
		return http.StatusNoContent
	}
	return r.last().StatusCode()
}

// body returns the Body object of the last httpBackendResponse in the httpResponseHistory list.
// Creates a new Body object using the last httpBackendResponse in the httpResponseHistory.
// Returns the newly created Body object.
func (r httpResponseHistory) body() *Body {
	if helper.IsNil(r.last()) {
		return nil
	}
	return r.last().Body()
}

func (r httpResponseHistory) aggregateBody() *Body {
	// instanciamos primeiro o aggregate body para retornar
	historyBody := NewBodyJson()
	// iteramos o histórico de backends httpResponse
	for index, httpBackendResponse := range r {
		// se tiver nil pulamos para o próximo
		if helper.IsNil(httpBackendResponse.Body()) {
			continue
		}
		// instânciamos o body
		body := httpBackendResponse.Body()
		// caso seja string ou slice agregamos na chave
		if httpBackendResponse.GroupByType() {
			historyBody = historyBody.AggregateByKey(httpBackendResponse.Key(index), body)
		} else {
			historyBody = historyBody.Aggregate(body)
		}
	}
	// se tudo ocorreu bem retornamos o corpo agregado
	return historyBody
}

// sliceOfBodies iterates over the httpResponse history and constructs a slice of bodies from the backend responses.
// If a backend httpResponse has an empty body, it is skipped.
// For each backend httpResponse with a non-empty body, a bodyBackendResponse object is created and added to the list of bodies.
// Returns a new Body object that contains the aggregated list of bodies from the httpResponse history.
func (r httpResponseHistory) sliceOfBodies() *Body {
	// instanciamos o valor a ser construído
	var bodies []*Body
	// iteramos o histórico para listar os bodies de resposta
	for index, httpBackendResponse := range r {
		// se tiver vazio pulamos para o próximo
		if helper.IsNil(httpBackendResponse.Body()) {
			continue
		}
		// inicializamos o body adicionando o campo "ok" e "code"
		bodyHttpBackendResponse := NewBodyByHttpBackendResponse(index, httpBackendResponse)
		// inserimos na lista de retorno
		bodies = append(bodies, bodyHttpBackendResponse)
	}
	// se tudo ocorrer bem, teremos o body agregado
	return NewBodyBySlice(bodies)
}

func (r httpResponseHistory) mostFrequentStatusCode() int {
	// instanciamos um mapa de status code com um count como valor
	statusCodes := make(map[int]int)
	// iteramos o histórico para alimentar esse map
	for _, httpBackendResponse := range r {
		statusCodes[httpBackendResponse.StatusCode()]++
	}
	// instanciamos o maxCount para manter o máximo utilizado
	maxCount := 0
	// instanciamos o mostFrequentCode para manter o statusCode mais utilizado
	mostFrequentCode := 0
	// iteramos para saber qual foi o mais frequente
	for code, count := range statusCodes {
		if count >= maxCount {
			mostFrequentCode = code
			maxCount = count
		}
	}
	// retornamos o statusCode mais frequente do histórico
	return mostFrequentCode
}

func (r httpResponseHistory) multipleResponseBody(aggregate bool) *Body {
	if aggregate {
		return r.aggregateBody()
	}
	return r.sliceOfBodies()
}
