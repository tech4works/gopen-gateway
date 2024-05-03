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
type httpResponseHistory []*httpBackendResponse

// Size returns the number of elements in the httpResponseHistory.
func (r httpResponseHistory) Size() int {
	return len(r)
}

// Success returns a boolean value indicating whether all backend responses in the httpResponseHistory are successful.
// It iterates through each backend httpResponse in the httpResponseHistory.
// If any backend httpResponse is not successful (Ok() returns false), it returns false.
// Otherwise, it returns true, indicating that all backend responses are successful.
func (r httpResponseHistory) Success() bool {
	for _, backendResponseVO := range r {
		if !backendResponseVO.Ok() {
			return false
		}
	}
	return true
}

func (r httpResponseHistory) Filter() (filteredHistory httpResponseHistory) {
	// iteramos o histórico para ser filtrado
	for _, httpBackendResponseVO := range r {
		// aplicamos a config de forma LATE
		httpBackendResponseAppliedVO := httpBackendResponseVO.ApplyConfig(enum.BackendResponseApplyLate)
		// setamos a resposta caso não esteja nil
		if helper.IsNotNil(httpBackendResponseAppliedVO) {
			filteredHistory = append(filteredHistory, httpBackendResponseAppliedVO)
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

// StatusCode function provides HTTP Status code based on the httpResponse history.
// It checks if there are multiple or single responses, and returns the relevant HTTP status code.
// If there is more than one httpResponse, it returns HTTP status code 200 (OK).
// If there is a single httpResponse, it returns the HTTP status code of that particular httpResponse.
// If there are no responses, it returns HTTP status code 204 (No Content).
func (r httpResponseHistory) StatusCode() int {
	// se tiver mais de 1 resposta obtemos o código de status mais frequente
	if r.MultipleResponse() {
		return r.mostFrequentStatusCode()
	} else if r.SingleResponse() {
		return r[0].statusCode
	}
	// resposta padrão de sucesso
	return http.StatusNoContent
}

// Header iterates over the httpResponseHistory and aggregates non-nil headers from each backendResponseVO.
// It constructs a final Header value object, which is an aggregation of all the individual non-nil headers.
// The method returns this final aggregated header.
//
// This function wouldn't consider any backendResponseVO whose header is nil.
// Also, with every aggregation, a new Header value object is created.
func (r httpResponseHistory) Header() Header {
	h := Header{}
	for _, httpBackendResponseVO := range r {
		// se tiver nil ou vazio pulamos para o próximo
		if helper.IsNil(httpBackendResponseVO.Header()) || helper.IsEmpty(httpBackendResponseVO.Header()) {
			continue
		}
		// agregamos os valores ao header de resultado, gerando um novo objeto de valor a cada agregação
		h = h.Aggregate(httpBackendResponseVO.Header())
	}
	return h
}

// Body function takes 'aggregate' boolean as argument.
// If there are multiple responses present, it aggregates them based on
// the boolean parameter. If the boolean parameter is true, it returns the result
// of 'aggregateBody' function else 'aggregatedBodies' function is called.
// For a single httpResponse case, the body of the single httpResponse is returned.
// If no responses are available, returns an empty Body struct.
func (r httpResponseHistory) Body(aggregate bool) *Body {
	if r.MultipleResponse() {
		if aggregate {
			return r.aggregateBody()
		} else {
			return r.sliceOfBodies()
		}
	} else if r.SingleResponse() {
		return r.body()
	}
	return nil
}

// Map executes the Map method on each backendResponseVO in the httpResponseHistory slice,
// collecting the results in a new slice of interface{}.
// It then converts the collected results to a string using the SimpleConvertToString function from the helper package.
// Returns the string representation of the collected results.
func (r httpResponseHistory) Map() string {
	var evalHistory []any
	for _, backendResponseVO := range r {
		evalHistory = append(evalHistory, backendResponseVO.Map())
	}
	return helper.SimpleConvertToString(evalHistory)
}

// last returns the last httpBackendResponse in the httpResponseHistory list.
// Returns the last httpBackendResponse object.
// Does not modify the httpResponseHistory.
func (r httpResponseHistory) last() *httpBackendResponse {
	return r[len(r)-1]
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

// aggregateBody aggregates the body from each backend httpResponse in the httpResponse history.
// It creates an initial Body object with empty JSON content.
// Then, it iterates through each backend httpResponse,
// skipping responses with a nil body.
// If the backend httpResponse has a group httpResponse flag set to true,
// it aggregates the body by key using the AggregateByKey method of the Body object.
// Otherwise, it aggregates all the JSON fields into the body using the Aggregate method of the Body object.
// Returns the final aggregated body.
// If there are no non-nil bodies in the httpResponse history, returns an empty Body struct.
func (r httpResponseHistory) aggregateBody() *Body {
	// instanciamos primeiro o aggregate body para retornar
	bodyHistory := &Body{
		contentType: enum.ContentTypeJson,
		value:       helper.SimpleConvertToBuffer("{}"),
	}

	// iteramos o histórico de backends httpResponse
	for index, backendResponseVO := range r {
		// se tiver nil pulamos para o próximo
		if helper.IsNil(backendResponseVO.body) {
			continue
		}
		// instânciamos o body
		body := backendResponseVO.Body()
		// caso seja string ou slice agregamos na chave
		if backendResponseVO.GroupByType() {
			bodyHistory = bodyHistory.AggregateByKey(backendResponseVO.Key(index), body)
		} else {
			bodyHistory = bodyHistory.Aggregate(body)
		}
	}

	// se tudo ocorreu bem retornamos o corpo agregado
	return bodyHistory
}

// sliceOfBodies iterates over the httpResponse history and constructs a slice of bodies from the backend responses.
// If a backend httpResponse has an empty body, it is skipped.
// For each backend httpResponse with a non-empty body, a bodyBackendResponse object is created and added to the list of bodies.
// Returns a new Body object that contains the aggregated list of bodies from the httpResponse history.
func (r httpResponseHistory) sliceOfBodies() *Body {
	// instanciamos o valor a ser construído
	var bodies []*Body
	// iteramos o histórico para listar os bodies de resposta
	for index, backendResponseVO := range r {
		// se tiver vazio pulamos para o próximo
		if helper.IsNil(backendResponseVO.body) {
			continue
		}
		// inicializamos o body agregado
		bodyBackendResponse := newBodyByHttpBackendResponse(index, backendResponseVO)
		// inserimos na lista de retorno
		bodies = append(bodies, bodyBackendResponse)
	}
	// se tudo ocorrer bem, teremos o body agregado
	return newBodyBySlice(bodies)
}

// mostFrequentStatusCode returns the most frequent status code from the httpResponse history.
// It counts the occurrences of each status code and finds the one with the highest count.
// If multiple status codes have the same highest count, it returns the first one encountered.
// Returns the most frequent status code as an integer.
// If the httpResponse history is empty, it returns 0.
func (r httpResponseHistory) mostFrequentStatusCode() int {
	statusCodes := make(map[int]int)
	for _, backendResponseVO := range r {
		statusCodes[backendResponseVO.statusCode]++
	}

	maxCount := 0
	mostFrequentCode := 0
	for code, count := range statusCodes {
		if count >= maxCount {
			mostFrequentCode = code
			maxCount = count
		}
	}

	return mostFrequentCode
}
