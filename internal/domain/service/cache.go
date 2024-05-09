package service

import (
	"context"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/interfaces"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
)

type cacheService struct {
	cacheStore interfaces.CacheStore
}

type Cache interface {
	Read(ctx context.Context, cache *vo.Cache, httpRequest *vo.HttpRequest) (*vo.CacheResponse, error)
	Write(ctx context.Context, cache *vo.Cache, httpRequest *vo.HttpRequest, httpResponse *vo.HttpResponse) error
}

func NewCache(cacheStore interfaces.CacheStore) Cache {
	return cacheService{
		cacheStore: cacheStore,
	}
}

func (c cacheService) Read(ctx context.Context, cache *vo.Cache, httpRequest *vo.HttpRequest) (*vo.CacheResponse, error) {
	// verificamos se ele não pode ler, retornamos nil
	if cache.CantRead(httpRequest) {
		return nil, nil
	}

	// inicializamos o valor a ser obtido
	var cacheResponse vo.CacheResponse

	// obtemos através do cache store se a chave exists respondemos, se não seguimos normalmente
	err := c.cacheStore.Get(ctx, cache.StrategyKey(httpRequest), &cacheResponse)
	if errors.Is(err, mapper.ErrCacheNotFound) {
		return nil, nil
	} else if helper.IsNotNil(err) {
		return nil, err
	}

	// se tudo ocorreu bem, retornamos a resposta do cache
	return &cacheResponse, nil
}

func (c cacheService) Write(ctx context.Context, cache *vo.Cache, httpRequest *vo.HttpRequest,
	httpResponse *vo.HttpResponse) error {
	// verificamos se ele não pode escrever, retornamos nil
	if cache.CantWrite(httpRequest, httpResponse) {
		return nil
	}

	// instanciamos a duração
	duration := cache.Duration()
	// construímos o valor a ser setado no cache
	cacheResponse := vo.NewCacheResponse(httpResponse, duration)

	// transformamos em cacheResponse e setamos
	return c.cacheStore.Set(ctx, cache.StrategyKey(httpRequest), cacheResponse, duration.Time())
}
