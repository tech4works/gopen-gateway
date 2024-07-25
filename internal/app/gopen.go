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

package app

import (
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/controller"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/middleware"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/usecase"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/factory"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/service"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/boot"
	"github.com/gin-gonic/gin"
	"net/http"
)

type gopenAPP struct {
	gopen                   *vo.Gopen
	panicRecoveryMiddleware middleware.PanicRecovery
	securityCorsMiddleware  middleware.SecurityCors
	timeoutMiddleware       middleware.Timeout
	limiterMiddleware       middleware.Limiter
	cacheMiddleware         middleware.Cache
	staticController        controller.Static
	endpointController      controller.Endpoint
	httpServer              *http.Server
}

type Gopen interface {
	ListerAndServer()
	Shutdown(ctx context.Context) error
}

func NewGopen(gopenDTO *dto.Gopen, httpClient HTTPClient, jsonPath domain.JSONPath, converter domain.Converter,
	store domain.Store) Gopen {
	boot.PrintInfo("Building value objects...")
	gopen := vo.NewGopen(gopenDTO)

	boot.PrintInfo("Building domain...")
	mapperService := service.NewMapper(jsonPath)
	projectorService := service.NewProjector(jsonPath)
	dynamicValueService := service.NewDynamicValue(jsonPath)
	modifierService := service.NewModifier(jsonPath)
	omitterService := service.NewOmitter(jsonPath)
	nomenclatureService := service.NewNomenclature(jsonPath)
	contentService := service.NewContent(converter)
	aggregatorService := service.NewAggregator(jsonPath)
	cacheService := service.NewCache(store)

	boot.PrintInfo("Building factories...")
	httpBackendFactory := factory.NewHTTPBackend(mapperService, projectorService, dynamicValueService, modifierService,
		omitterService, nomenclatureService, contentService, aggregatorService)
	httpResponseFactory := factory.NewHTTPResponse(aggregatorService, omitterService, nomenclatureService, contentService,
		httpBackendFactory)

	boot.PrintInfo("Building use cases...")
	endpointUseCase := usecase.NewEndpoint(httpBackendFactory, httpResponseFactory, httpClient)

	boot.PrintInfo("Building middlewares...")
	panicRecoveryMiddleware := middleware.NewPanicRecovery(logProvider)
	securityCorsMiddleware := middleware.NewSecurityCors(gopen.SecurityCors())
	timeoutMiddleware := middleware.NewTimeout()
	limiterMiddleware := middleware.NewLimiter()
	cacheMiddleware := middleware.NewCache(cacheService, logProvider)

	boot.PrintInfo("Building controllers...")
	staticController := controller.NewStatic(gopenDTO.Version, factory.BuildSettingViewDTO(gopenDTO, gopen))
	endpointController := controller.NewEndpoint(endpointUseCase)

	return gopenAPP{
		gopen:                   gopen,
		panicRecoveryMiddleware: panicRecoveryMiddleware,
		timeoutMiddleware:       timeoutMiddleware,
		limiterMiddleware:       limiterMiddleware,
		cacheMiddleware:         cacheMiddleware,
		securityCorsMiddleware:  securityCorsMiddleware,
		staticController:        staticController,
		endpointController:      endpointController,
	}
}

func (g gopenAPP) ListerAndServer() {
	boot.PrintInfo("Starting lister and server...")

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	g.buildStaticRoutes(engine)

	boot.PrintInfo("Starting to read endpoints to register routes...")
	for _, endpoint := range g.gopen.Endpoints() {
		handles := g.buildEndpointHandles()
		api.Handle(engine, g.gopen, &endpoint, handles...)

		lenString := helper.SimpleConvertToString(len(handles))
		boot.PrintInfof("Registered route with %s handles: %s", lenString, endpoint.Resume())
	}

	address := fmt.Sprint(":", g.gopen.Port())
	boot.PrintInfof("Listening and serving HTTP on %s!", address)

	g.httpServer = &http.Server{
		Addr:    address,
		Handler: engine,
	}

	fmt.Println()
	fmt.Println()
	boot.PrintTitle("LISTEN AND SERVER")

	g.httpServer.ListenAndServe()
}

func (g gopenAPP) Shutdown(ctx context.Context) error {
	if helper.IsNil(g.httpServer) {
		return nil
	}
	return g.httpServer.Shutdown(ctx)
}

func (g gopenAPP) buildStaticRoutes(engine *gin.Engine) {
	boot.PrintInfo("Configuring static routes...")
	formatLog := "Registered route with 5 handles: %s --> \"%s\""

	pingEndpoint := g.registerStaticPingRoute(engine)
	boot.PrintInfof(formatLog, pingEndpoint.Method(), pingEndpoint.Path())

	versionEndpoint := g.registerStaticVersionRoute(engine)
	boot.PrintInfof(formatLog, versionEndpoint.Method(), versionEndpoint.Path())

	settingsEndpoint := g.registerStaticSettingsRoute(engine)
	boot.PrintInfof(formatLog, settingsEndpoint.Method(), settingsEndpoint.Path())
}

func (g gopenAPP) registerStaticPingRoute(engine *gin.Engine) *vo.Endpoint {
	endpoint := vo.NewEndpointStatic("/ping", http.MethodGet)
	g.registerStaticRoute(engine, &endpoint, g.staticController.Ping)
	return &endpoint
}

func (g gopenAPP) registerStaticVersionRoute(engine *gin.Engine) *vo.Endpoint {
	endpoint := vo.NewEndpointStatic("/version", http.MethodGet)
	g.registerStaticRoute(engine, &endpoint, g.staticController.Version)
	return &endpoint
}

func (g gopenAPP) registerStaticSettingsRoute(engine *gin.Engine) *vo.Endpoint {
	endpoint := vo.NewEndpointStatic("/settings", http.MethodGet)
	g.registerStaticRoute(engine, &endpoint, g.staticController.Settings)
	return &endpoint
}

func (g gopenAPP) registerStaticRoute(engine *gin.Engine, endpointStatic *vo.Endpoint, handler api.HandlerFunc) {
	timeoutHandler := g.timeoutMiddleware.Do
	panicHandler := g.panicRecoveryMiddleware.Do
	limiterHandler := g.limiterMiddleware.Do
	api.Handle(engine, g.gopen, endpointStatic, timeoutHandler, panicHandler, limiterHandler, handler)
}

func (g gopenAPP) buildEndpointHandles() []api.HandlerFunc {
	timeoutHandler := g.timeoutMiddleware.Do
	panicHandler := g.panicRecoveryMiddleware.Do
	securityCorsHandler := g.securityCorsMiddleware.Do
	limiterHandler := g.limiterMiddleware.Do
	cacheHandler := g.cacheMiddleware.Do
	endpointHandler := g.endpointController.Execute
	return []api.HandlerFunc{
		timeoutHandler,
		panicHandler,
		securityCorsHandler,
		limiterHandler,
		cacheHandler,
		endpointHandler,
	}
}
