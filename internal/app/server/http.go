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

package server

import (
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/controller"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/factory"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/middleware"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/usecase"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain"
	domainFactory "github.com/GabrielHCataldo/gopen-gateway/internal/domain/factory"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/service"
	net "net/http"
)

type http struct {
	*net.Server
	gopen                   *vo.Gopen
	logger                  app.Logger
	router                  app.Router
	panicRecoveryMiddleware middleware.PanicRecovery
	securityCorsMiddleware  middleware.SecurityCors
	timeoutMiddleware       middleware.Timeout
	limiterMiddleware       middleware.Limiter
	cacheMiddleware         middleware.Cache
	staticController        controller.Static
	endpointController      controller.Endpoint
}

type HTTP interface {
	ListerAndServe()
	Shutdown(ctx context.Context) error
}

func New(
	gopen *dto.Gopen,
	logger app.Logger,
	router app.Router,
	httpClient app.HTTPClient,
	jsonPath domain.JSONPath,
	converter domain.Converter,
	store domain.Store,
) HTTP {
	logger.PrintInfo("Building domain...")
	mapperService := service.NewMapper(jsonPath)
	projectorService := service.NewProjector(jsonPath)
	dynamicValueService := service.NewDynamicValue(jsonPath)
	modifierService := service.NewModifier(jsonPath)
	omitterService := service.NewOmitter(jsonPath)
	nomenclatureService := service.NewNomenclature(jsonPath)
	contentService := service.NewContent(converter)
	aggregatorService := service.NewAggregator(jsonPath)
	limiterService := service.NewLimiter()
	securityCorsService := service.NewSecurityCors()
	cacheService := service.NewCache(store)

	logger.PrintInfo("Building factories...")
	httpBackendFactory := domainFactory.NewHTTPBackend(mapperService, projectorService, dynamicValueService,
		modifierService, omitterService, nomenclatureService, contentService, aggregatorService)
	httpResponseFactory := domainFactory.NewHTTPResponse(aggregatorService, omitterService, nomenclatureService,
		contentService, httpBackendFactory)

	logger.PrintInfo("Building use cases...")
	endpointUseCase := usecase.NewEndpoint(httpBackendFactory, httpResponseFactory, httpClient)

	logger.PrintInfo("Building middlewares...")
	panicRecoveryMiddleware := middleware.NewPanicRecovery()
	securityCorsMiddleware := middleware.NewSecurityCors(securityCorsService)
	timeoutMiddleware := middleware.NewTimeout()
	limiterMiddleware := middleware.NewLimiter(limiterService)
	cacheMiddleware := middleware.NewCache(cacheService)

	logger.PrintInfo("Building controllers...")
	staticController := controller.NewStatic(gopen)
	endpointController := controller.NewEndpoint(endpointUseCase)

	logger.PrintInfo("Building value objects...")
	return http{
		gopen:                   factory.BuildGopen(gopen),
		logger:                  logger,
		router:                  router,
		panicRecoveryMiddleware: panicRecoveryMiddleware,
		timeoutMiddleware:       timeoutMiddleware,
		limiterMiddleware:       limiterMiddleware,
		cacheMiddleware:         cacheMiddleware,
		securityCorsMiddleware:  securityCorsMiddleware,
		staticController:        staticController,
		endpointController:      endpointController,
	}
}

func (h http) ListerAndServe() {
	h.logger.PrintInfo("Starting lister and server...")

	h.buildStaticRoutes()

	h.logger.PrintInfo("Starting to read endpoints to register routes...")
	for _, endpoint := range h.gopen.Endpoints() {
		handles := h.buildEndpointHandles()
		h.router.Handle(h.gopen, &endpoint, handles...)

		lenString := helper.SimpleConvertToString(len(handles))
		h.logger.PrintInfof("Registered route with %s handles: %s", lenString, endpoint.Resume())
	}

	address := fmt.Sprint(":", h.gopen.Port())
	h.logger.PrintInfof("Listening and serving HTTP on %s!", address)

	h.Server = &net.Server{
		Addr:    address,
		Handler: h.router.Engine(),
	}

	fmt.Println()
	fmt.Println()
	h.logger.PrintTitle("LISTEN AND SERVER")

	h.Server.ListenAndServe()
}

func (h http) Shutdown(ctx context.Context) error {
	if helper.IsNil(h.Server) {
		return nil
	}
	return h.Server.Shutdown(ctx)
}

func (h http) buildStaticRoutes() {
	h.logger.PrintInfo("Configuring static routes...")
	formatLog := "Registered route with 5 handles: %s --> \"%s\""

	pingEndpoint := h.buildStaticPingRoute()
	h.logger.PrintInfof(formatLog, pingEndpoint.Method(), pingEndpoint.Path())

	versionEndpoint := h.buildStaticVersionRoute()
	h.logger.PrintInfof(formatLog, versionEndpoint.Method(), versionEndpoint.Path())

	settingsEndpoint := h.buildStaticSettingsRoute()
	h.logger.PrintInfof(formatLog, settingsEndpoint.Method(), settingsEndpoint.Path())
}

func (h http) buildStaticPingRoute() *vo.Endpoint {
	endpoint := vo.NewEndpointStatic("/ping", net.MethodGet)
	h.buildStaticRoute(&endpoint, h.staticController.Ping)
	return &endpoint
}

func (h http) buildStaticVersionRoute() *vo.Endpoint {
	endpoint := vo.NewEndpointStatic("/version", net.MethodGet)
	h.buildStaticRoute(&endpoint, h.staticController.Version)
	return &endpoint
}

func (h http) buildStaticSettingsRoute() *vo.Endpoint {
	endpoint := vo.NewEndpointStatic("/settings", net.MethodGet)
	h.buildStaticRoute(&endpoint, h.staticController.Settings)
	return &endpoint
}

func (h http) buildStaticRoute(endpointStatic *vo.Endpoint, handler app.HandlerFunc) {
	timeoutHandler := h.timeoutMiddleware.Do
	panicHandler := h.panicRecoveryMiddleware.Do
	limiterHandler := h.limiterMiddleware.Do
	h.router.Handle(h.gopen, endpointStatic, timeoutHandler, panicHandler, limiterHandler, handler)
}

func (h http) buildEndpointHandles() []app.HandlerFunc {
	timeoutHandler := h.timeoutMiddleware.Do
	panicHandler := h.panicRecoveryMiddleware.Do
	securityCorsHandler := h.securityCorsMiddleware.Do
	limiterHandler := h.limiterMiddleware.Do
	cacheHandler := h.cacheMiddleware.Do
	endpointHandler := h.endpointController.Execute
	return []app.HandlerFunc{
		timeoutHandler,
		panicHandler,
		securityCorsHandler,
		limiterHandler,
		cacheHandler,
		endpointHandler,
	}
}
