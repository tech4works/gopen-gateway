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
	"github.com/gin-gonic/gin"
	"net/http"
)

// static represents a static handler that handles requests for ping and version.
// It contains a gopenVO field of type vo.Gopen, which holds configuration information.
type static struct {
	gopenVO *vo.Gopen
}

// Static represents an interface for handling requests related to ping, version, and settings.
type Static interface {
	// Ping handles the GET request to the "/ping" endpoint. It is a method of the Static interface
	// that returns a response containing the string "pong". This method is used in the buildStaticRoutes
	// method of the gopen type to configure the "/ping" route for the Gin engine.
	Ping(ctx *gin.Context)
	// Version is a method of the Static interface that handles the GET request to the "/version" endpoint.
	// It takes a *gin.Context parameter and performs the necessary actions to return a response.
	Version(ctx *gin.Context)
	// Settings behaves as a method of the Static interface that handles the GET request to the "/settings" endpoint.
	// It takes a *gin.Context parameter and performs the necessary actions to return a response.
	Settings(ctx *gin.Context)
}

// NewStatic is a function that creates a new instance of the Static interface.
// It takes a vo.Gopen parameter and returns a Static object.
// The returned Static object has a gopenVO field which is initialized with the provided vo.Gopen object.
func NewStatic(gopenVO *vo.Gopen) Static {
	return static{
		gopenVO: gopenVO,
	}
}

// Ping is a method that handles the "Ping" request.
// It responds with the string "Pong!" and a status code of 200 (OK).
func (s static) Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "%s", "Pong!")
}

// Version is a method that handles the "Version" request.
// It checks if the version is not empty and responds with the version string and a status code of 200 (OK).
// If the version is empty, it responds with a status code of 404 (Not Found).
func (s static) Version(ctx *gin.Context) {
	if helper.IsNotEmpty(s.gopenVO.Version()) {
		ctx.String(http.StatusOK, "%s", s.gopenVO.Version())
		return
	}
	ctx.Status(http.StatusNotFound)
}

// Settings is a method that handles the "Settings" request.
// It retrieves the necessary data from the gopenVO object and constructs a dto.SettingView object.
// The SettingView object contains information such as version, version date, founder, code helpers,
// number of endpoints, number of middlewares, number of backends, number of modifiers, and the gopenVO object itself.
// The method then responds with the constructed SettingView object in JSON format and a status code of 200 (OK).
func (s static) Settings(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, mapper.BuildSettingViewDTO(s.gopenVO))
}
