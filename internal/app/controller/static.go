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
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"net/http"
)

type staticController struct {
	version        string
	settingViewDTO dto.SettingView
}

type Static interface {
	Ping(ctx *api.Context)
	Version(ctx *api.Context)
	Settings(ctx *api.Context)
}

func NewStatic(version string, settingView dto.SettingView) Static {
	return staticController{
		version:        version,
		settingViewDTO: settingView,
	}
}

func (s staticController) Ping(ctx *api.Context) {
	ctx.WriteString(http.StatusOK, "Pong!")
}

func (s staticController) Version(ctx *api.Context) {
	if helper.IsNotEmpty(s.version) {
		ctx.WriteString(http.StatusOK, s.version)
		return
	}
	ctx.WriteStatusCode(http.StatusNotFound)
}

func (s staticController) Settings(ctx *api.Context) {
	ctx.WriteJson(http.StatusOK, s.settingViewDTO)
}
