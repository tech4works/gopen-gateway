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
	"github.com/tech4works/checker"
)

type Gopen struct {
	proxy        *Proxy
	securityCors *SecurityCors
	endpoints    []Endpoint
}

func NewGopen(proxy *Proxy, securityCors *SecurityCors, endpoints []Endpoint) *Gopen {
	return &Gopen{
		proxy:        proxy,
		securityCors: securityCors,
		endpoints:    endpoints,
	}
}

func (g Gopen) Proxy() *Proxy {
	return g.proxy
}

func (g Gopen) HasProxy() bool {
	return checker.NonNil(g.proxy)
}

func (g Gopen) SecurityCors() *SecurityCors {
	return g.securityCors
}

func (g Gopen) HasSecurityCors() bool {
	return checker.NonNil(g.securityCors)
}

func (g Gopen) Endpoints() []Endpoint {
	return g.endpoints
}
