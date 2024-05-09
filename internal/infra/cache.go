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
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/interfaces"
)

type cacheStoreProvider struct {
}

func NewCacheStoreProvider() interfaces.CacheStoreProvider {
	return cacheStoreProvider{}
}

func (c cacheStoreProvider) Memory() interfaces.CacheStore {
	return NewMemoryStore()
}

func (c cacheStoreProvider) Redis(address, password string) interfaces.CacheStore {
	return NewRedisStore(address, password)
}
