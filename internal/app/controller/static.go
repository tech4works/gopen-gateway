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
	"github.com/GabrielHCataldo/gopen-gateway/internal/app"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"net/http"
)

type staticController struct {
	gopenDTO *dto.Gopen
}

type Static interface {
	Ping(ctx app.Context)
	Version(ctx app.Context)
	Settings(ctx app.Context)
}

func NewStatic(gopenDTO *dto.Gopen) Static {
	return staticController{
		gopenDTO: gopenDTO,
	}
}

func (s staticController) Ping(ctx app.Context) {
	ctx.WriteString(http.StatusOK, "Pong!")
}

func (s staticController) Version(ctx app.Context) {
	if helper.IsNotEmpty(s.gopenDTO.Version) {
		ctx.WriteString(http.StatusOK, s.gopenDTO.Version)
		return
	}
	ctx.WriteStatusCode(http.StatusNotFound)
}

func (s staticController) Settings(ctx app.Context) {
	// todo: aq fazer o build SettingView a partir do DTO
	//		ctx.WriteJson(http.StatusOK, s.settingViewDTO)
	ctx.WriteStatusCode(http.StatusNotImplemented)
}
