package valueobject

import "github.com/GabrielHCataldo/martini-gateway/internal/application/model/dto"

type Backend struct {
	Host           []string
	Endpoint       string
	Method         string
	Group          string
	HideResponse   bool
	ForwardHeaders []string
	Query          []string
	Headers        []Modifier
	Params         []Modifier
	Queries        []Modifier
	Body           []Modifier
}

func BuildBackend(backend dto.Backend) Backend {
	return Backend{
		Host:           backend.Host,
		Endpoint:       backend.Endpoint,
		Method:         backend.Method,
		Group:          backend.Group,
		HideResponse:   backend.HideResponse,
		ForwardHeaders: backend.ForwardHeaders,
		Query:          backend.Query,
		Headers:        BuildModifiers(backend.Headers),
		Params:         BuildModifiers(backend.Params),
		Queries:        BuildModifiers(backend.Queries),
		Body:           BuildModifiers(backend.Body),
	}
}
