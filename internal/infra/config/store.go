package config

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

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra"
)

// NewCacheStore configures and builds a cache store based on the provided storeJsonVO.
// If storeJsonVO is not nil, it creates a Redis cache store with the specified Redis address and password.
// Otherwise, it creates a Memory cache store.
// It returns a CacheStore interface that can be used to interact with the cache store.
func NewCacheStore(storeJsonVO *vo.StoreJson) infra.CacheStore {
	PrintInfoLogCmd("Configuring cache store...")
	if helper.IsNotNil(storeJsonVO) {
		return infra.NewRedisStore(storeJsonVO.Redis.Address, storeJsonVO.Redis.Password)
	}
	return infra.NewMemoryStore()
}

// CloseCacheStore closes the cache store by calling the Close method on the provided infra.CacheStore object.
// If an error occurs during the close operation, a warning log is printed.
func CloseCacheStore(store infra.CacheStore) {
	err := store.Close()
	if helper.IsNotNil(err) {
		PrintWarningLogCmd("Error close cache store:", err)
	}
}
