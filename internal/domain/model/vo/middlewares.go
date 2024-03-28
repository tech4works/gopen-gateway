package vo

import "github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"

type Middlewares map[string]Backend

func newMiddlewares(middlewaresDTO map[string]dto.Backend) (m Middlewares) {
	m = Middlewares{}
	for k, v := range middlewaresDTO {
		m[k] = newBackend(v)
	}
	return m
}

func (m Middlewares) Get(key string) (Backend, bool) {
	backend, ok := m[key]
	if !ok {
		return Backend{}, false
	}
	return newMiddlewareBackend(backend, backendExtraConfig{
		omitResponse: true,
	}), true
}
