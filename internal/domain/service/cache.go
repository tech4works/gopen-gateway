/*
 * Copyright 2024 Tech4Works
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
	"fmt"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/model/aggregate"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type cache struct {
	dynamicValueService DynamicValue
	store               domain.Store
}

type Cache interface {
	Read(
		ctx context.Context,
		config *vo.CacheConfig,
		request *vo.EndpointRequest,
		history *aggregate.History,
		dest any,
	) error
	Write(
		ctx context.Context,
		config *vo.CacheConfig,
		cacheable vo.Cacheable,
		request *vo.EndpointRequest,
		history *aggregate.History,
	) error
}

func NewCache(dynamicValueService DynamicValue, store domain.Store) Cache {
	return cache{
		dynamicValueService: dynamicValueService,
		store:               store,
	}
}

func (c cache) Read(
	ctx context.Context,
	config *vo.CacheConfig,
	request *vo.EndpointRequest,
	history *aggregate.History,
	dest any,
) error {
	shouldRun, err := c.evalCacheGuards("read", config.Read(), request, history)
	if checker.NonNil(err) {
		return errors.Inheritf(err, "cache failed: op=eval-guards from=write")
	} else if !shouldRun {
		return nil
	}

	key, err := c.buildKey(config, request, history)
	if checker.NonNil(err) {
		return err
	}

	entry, err := c.store.Get(ctx, key)
	if errors.Is(err, domain.ErrCacheNotFound) {
		return nil
	} else if checker.NonNil(err) {
		return errors.Inheritf(err, "cache failed: unexpected error reading cache key=%s", key)
	}

	return converter.ToDestWithErr(entry, dest)
}

func (c cache) Write(
	ctx context.Context,
	config *vo.CacheConfig,
	cacheable vo.Cacheable,
	request *vo.EndpointRequest,
	history *aggregate.History,
) error {
	shouldRun, err := c.evalCacheGuards("write", config.Write(), request, history)
	if checker.NonNil(err) {
		return errors.Inheritf(err, "cache failed: op=eval-guards from=write")
	} else if !shouldRun {
		return nil
	}

	key, err := c.buildKey(config, request, history)
	if checker.NonNil(err) {
		return err
	}

	entry, err := cacheable.Entry()
	if checker.NonNil(err) {
		return errors.Inheritf(err, "cache failed: op=build-data-entry")
	}

	err = c.store.Set(ctx, key, entry, config.TTL().Time())
	if checker.NonNil(err) {
		return errors.Inheritf(err, "cache failed: unexpected error writing cache key=%s kind=%s", key, config.Kind())
	}

	return nil
}

func (c cache) evalCacheGuards(
	operation string,
	decision vo.CacheDecisionConfig,
	request *vo.EndpointRequest,
	history *aggregate.History,
) (bool, error) {
	shouldRun, _, errs := c.dynamicValueService.EvalGuards(decision.OnlyIf(), decision.IgnoreIf(), request, history)
	if checker.IsNotEmpty(errs) {
		return false, errors.JoinInheritf(errs, ", ", "failed to evaluate guard for cache %s", operation)
	}
	return shouldRun, nil
}

func (c cache) buildKey(config *vo.CacheConfig, request *vo.EndpointRequest, history *aggregate.History) (string, error) {
	key, errs := c.dynamicValueService.Get(config.Key(), request, history)
	if checker.IsNotEmpty(errs) {
		return "", errors.JoinInheritf(errs, ", ", "cache failed: op=build-key")
	}
	return fmt.Sprintf("%s:%s", config.Kind(), key), nil
}
