/*
 * Copyright 2024 Gabriel Cataldo
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"context"
	goerros "github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain"
	mapper2 "github.com/GabrielHCataldo/gopen-gateway/internal/domain/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
)

type cacheService struct {
	store domain.Store
}

type Cache interface {
	Read(ctx context.Context, cache *vo.Cache, request *vo.HTTPRequest) (*vo.CacheResponse, error)
	Write(ctx context.Context, cache *vo.Cache, request *vo.HTTPRequest, response *vo.HTTPResponse) error
}

func NewCache(store domain.Store) Cache {
	return cacheService{
		store: store,
	}
}

func (c cacheService) Read(ctx context.Context, cacheConfig *vo.Cache, request *vo.HTTPRequest) (
	*vo.CacheResponse, error) {
	if cacheConfig.CantRead(request) {
		return nil, nil
	}

	key := cacheConfig.StrategyKey(request)
	cacheResponse, err := c.store.Get(ctx, key)
	if goerros.Is(err, mapper2.ErrCacheNotFound) {
		return nil, nil
	} else if helper.IsNotNil(err) {
		return nil, err
	}

	return cacheResponse, nil
}

func (c cacheService) Write(ctx context.Context, cacheConfig *vo.Cache, request *vo.HTTPRequest,
	response *vo.HTTPResponse) error {
	if cacheConfig.CantWrite(request, response) {
		return nil
	}

	key := cacheConfig.StrategyKey(request)
	cacheResponse := vo.NewCacheResponse(cacheConfig, response)
	return c.store.Set(ctx, key, cacheResponse)
}
