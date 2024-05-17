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

package domain

import (
	"context"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
)

// CacheStore is an interface that defines methods for interacting with a cache store.
type CacheStore interface {
	Set(ctx context.Context, key string, value *vo.CacheResponse) error
	// Del deletes the cache entry with the specified key from the cache store.
	// The key is a string that serves as the identifier of the cache entry.
	// The error returned indicates any issues encountered while deleting the cache entry.
	// This method is used to remove a key-value pair from the cache store.
	// Implementing the CacheStore interface, this method is responsible for removing the cache entry using the underlying implementation.
	Del(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (*vo.CacheResponse, error)
	// Close is a method defined in the CacheStore interface. It is used to close the cache store and release any resources
	// associated with it. The method returns an error if there was a problem closing the store.
	Close() error
}

type RestTemplate interface {
	MakeRequest(ctx context.Context, backend *vo.Backend, httpBackendRequest *vo.HttpBackendRequest) (
		*vo.HttpBackendResponse, error)
}
