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

package vo

import (
	"time"
)

type CacheResponse struct {
	StatusCode StatusCode `json:"statusCode"`
	Header     Header     `json:"header"`
	Body       *Body      `json:"body,omitempty"`
	Duration   Duration   `json:"duration"`
	CreatedAt  time.Time  `json:"createdAt"`
}

func NewCacheResponse(cacheConfig *Cache, response *HTTPResponse) *CacheResponse {
	return &CacheResponse{
		StatusCode: response.StatusCode(),
		Header:     response.Header(),
		Body:       response.Body(),
		Duration:   cacheConfig.Duration(),
		CreatedAt:  time.Now(),
	}
}

func (r CacheResponse) TTL() string {
	timeDuration := r.Duration.Time()
	sub := r.CreatedAt.Add(timeDuration).Sub(time.Now())
	return sub.String()
}
