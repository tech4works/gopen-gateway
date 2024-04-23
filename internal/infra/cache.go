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

package infra

import (
	"context"
	"time"
)

// CacheStore is an interface that defines methods for interacting with a cache store.
type CacheStore interface {
	// Set stores the given value with the provided key in the cache store. The value will expire after the specified duration.
	// If an error occurs during the set operation, it is returned.
	Set(ctx context.Context, key string, value any, expire time.Duration) error
	// Del deletes the cache entry with the specified key from the cache store.
	// The key is a string that serves as the identifier of the cache entry.
	// The error returned indicates any issues encountered while deleting the cache entry.
	// This method is used to remove a key-value pair from the cache store.
	// Implementing the CacheStore interface, this method is responsible for removing the cache entry using the underlying implementation.
	Del(ctx context.Context, key string) error
	// Get retrieves the value from the cache store with the given key and stores it in the provided destination object.
	// The destination object must be a pointer to the type that matches the retrieved value.
	// If the cache entry is not found, or an error occurs while retrieving the value, an error is returned.
	// The method takes in a context.Context object to support cancellation and timeouts.
	Get(ctx context.Context, key string, dest any) error
	// Close is a method defined in the CacheStore interface. It is used to close the cache store and release any resources
	// associated with it. The method returns an error if there was a problem closing the store.
	Close() error
}
