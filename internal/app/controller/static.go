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

package controller

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"net/http"
)

// static represents a static handler that handles requests for ping and version.
// It contains a gopenJsonVO field of type vo.Gopen, which holds configuration information.
type static struct {
	gopenJsonVO *vo.GopenJson
}

// Static represents an interface for handling requests related to ping, version, and settings.
type Static interface {
	Ping(ctx *api.Context)
	Version(ctx *api.Context)
	Settings(ctx *api.Context)
}

// NewStatic is a function that creates a new instance of the Static interface.
// It takes a vo.Gopen parameter and returns a Static object.
// The returned Static object has a gopenVO field which is initialized with the provided vo.Gopen object.
func NewStatic(gopenJsonVO *vo.GopenJson) Static {
	return static{
		gopenJsonVO: gopenJsonVO,
	}
}

func (s static) Ping(ctx *api.Context) {
	ctx.WriteString(http.StatusOK, "Pong!")
}

func (s static) Version(ctx *api.Context) {
	if helper.IsNotEmpty(s.gopenJsonVO.Version) {
		ctx.WriteString(http.StatusOK, s.gopenJsonVO.Version)
		return
	}
	ctx.WriteStatusCode(http.StatusNotFound)
}

func (s static) Settings(ctx *api.Context) {
	settingViewDTO := mapper.BuildSettingViewDTO(s.gopenJsonVO)
	ctx.WriteJson(http.StatusOK, settingViewDTO)
}
