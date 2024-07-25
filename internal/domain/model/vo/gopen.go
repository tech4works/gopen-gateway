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

package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
)

type Gopen struct {
	port         int
	securityCors *SecurityCors
	middlewares  Middlewares
	endpoints    []Endpoint
}

func NewGopen(gopenJson *GopenJson) *Gopen {
	var endpoints []Endpoint
	for _, endpointJson := range gopenJson.Endpoints {
		endpoints = append(endpoints, newEndpoint(gopenJson, &endpointJson))
	}
	return &Gopen{
		port:         gopenJson.Port,
		securityCors: newSecurityCors(gopenJson.SecurityCors),
		middlewares:  newMiddlewares(gopenJson.Middlewares),
		endpoints:    endpoints,
	}
}

func (g Gopen) Port() int {
	return g.port
}

func (g Gopen) SecurityCors() *SecurityCors {
	return g.securityCors
}

func (g Gopen) NoSecurityCors() bool {
	return helper.IsNil(g.securityCors)
}

func (g Gopen) Middleware(key string) (*Backend, bool) {
	return g.middlewares.Get(key)
}

func (g Gopen) Endpoints() []Endpoint {
	return g.endpoints
}

func (g Gopen) CountEndpoints() int {
	return len(g.endpoints)
}

func (g Gopen) CountMiddlewares() int {
	if helper.IsNotNil(g.middlewares) {
		return len(g.middlewares)
	}
	return 0
}

func (g Gopen) CountBackends() (count int) {
	for _, endpoint := range g.Endpoints() {
		count += endpoint.CountBeforewares()
		count += endpoint.CountBackends()
		count += endpoint.CountAfterwares()
	}
	return count
}

func (g Gopen) CountAllDataTransforms() (count int) {
	for _, endpoint := range g.Endpoints() {
		count += endpoint.CountAllDataTransforms()
	}
	return count
}
