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

// static represents a static handler that handles requests for ping and version.
// It contains a gopenJson field of type vo.Gopen, which holds configuration information.
type static struct {
	version        string
	settingViewDTO dto.SettingView
}

// Static represents an interface for handling requests related to ping, version, and settings.
type Static interface {
	Ping(ctx *api.Context)
	Version(ctx *api.Context)
	Settings(ctx *api.Context)
}

// NewStatic is a function that creates a new instance of the Static interface.
// It takes a version string and a settingView DTO as input parameters and returns a static struct
// which implements the Static interface. The version string represents the version of the application,
// and the settingView DTO contains the configuration view for the application.
func NewStatic(version string, settingView dto.SettingView) Static {
	return static{
		version:        version,
		settingViewDTO: settingView,
	}
}

// Ping is a method of the static struct that handles a ping request.
// It takes a Context parameter and writes "Pong!" as the response body.
// The Context parameter represents the context of the request.
func (s static) Ping(ctx *api.Context) {
	ctx.WriteString(http.StatusOK, "Pong!")
}

// Version is a method of the static struct that handles a version request.
// It takes a Context parameter and writes the version string as the response body.
// If the version is empty, it writes a status code of http.StatusNotFound.
// The Context parameter represents the context of the request.
func (s static) Version(ctx *api.Context) {
	if helper.IsNotEmpty(s.version) {
		ctx.WriteString(http.StatusOK, s.version)
		return
	}
	ctx.WriteStatusCode(http.StatusNotFound)
}

// Settings is a method of the static struct that handles a settings request.
// It takes a Context parameter and writes the settingViewDTO as a JSON response body.
// The Context parameter represents the context of the request.
func (s static) Settings(ctx *api.Context) {
	ctx.WriteJson(http.StatusOK, s.settingViewDTO)
}
