package usecase

import (
	"github.com/GabrielHCataldo/martini-gateway/internal/domain/service"
)

func BuildModifierRequest(request Request) service.ModifierRequest {
	return service.ModifierRequest{
		Host:     request.Host,
		Endpoint: request.Endpoint,
		Url:      request.Url,
		Method:   request.Method,
		Header:   request.Header,
		Query:    request.Query,
		Params:   request.Params,
		Body:     request.Body,
	}
}

func BuildModifierRequestByBackendRequest(request service.BackendRequest) service.ModifierRequest {
	return service.ModifierRequest{
		Host:     request.Host,
		Endpoint: request.Endpoint,
		Url:      request.Url,
		Method:   request.Method,
		Header:   request.Header,
		Query:    request.Query,
		Params:   request.Params,
		Body:     request.Body,
	}
}

func BuildModifierRequests(requests []Request) []service.ModifierRequest {
	var result []service.ModifierRequest
	for _, request := range requests {
		result = append(result, BuildModifierRequest(request))
	}
	return result
}

func BuildModifierResponse(response Response) service.ModifierResponse {
	return service.ModifierResponse{
		StatusCode: response.StatusCode,
		Header:     response.Header,
		Body:       response.Body,
	}
}

func BuildModifierResponseByBackendResponse(response service.BackendResponse) service.ModifierResponse {
	return service.ModifierResponse{
		StatusCode: response.StatusCode,
		Header:     response.Header,
		Body:       response.Body,
	}
}

func BuildModifierResponses(responses []Response) []service.ModifierResponse {
	var result []service.ModifierResponse
	for _, response := range responses {
		result = append(result, BuildModifierResponse(response))
	}
	return result
}
