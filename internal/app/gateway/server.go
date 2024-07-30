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

package gateway

import (
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/controller"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/middleware"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/usecase"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/factory"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/service"
	"net/http"
)

type server struct {
	*http.Server
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

type Server interface {
	Start()
	Shutdown(ctx context.Context) error
}

func NewServer(
	gopenDTO *dto.Gopen,
	logger app.Logger,
	router app.Router,
	httpClient app.HTTPClient,
	jsonPath domain.JSONPath,
	converter domain.Converter,
	store domain.Store,
) Server {
	logger.PrintInfo("Building value objects...")
	gopen := vo.NewGopen(gopenDTO)

	logger.PrintInfo("Building domain...")
	mapperService := service.NewMapper(jsonPath)
	projectorService := service.NewProjector(jsonPath)
	dynamicValueService := service.NewDynamicValue(jsonPath)
	modifierService := service.NewModifier(jsonPath)
	omitterService := service.NewOmitter(jsonPath)
	nomenclatureService := service.NewNomenclature(jsonPath)
	contentService := service.NewContent(converter)
	aggregatorService := service.NewAggregator(jsonPath)
	securityCorsService := service.NewSecurityCors()
	cacheService := service.NewCache(store)

	logger.PrintInfo("Building factories...")
	httpBackendFactory := factory.NewHTTPBackend(mapperService, projectorService, dynamicValueService, modifierService,
		omitterService, nomenclatureService, contentService, aggregatorService)
	httpResponseFactory := factory.NewHTTPResponse(aggregatorService, omitterService, nomenclatureService, contentService,
		httpBackendFactory)

	logger.PrintInfo("Building use cases...")
	endpointUseCase := usecase.NewEndpoint(httpBackendFactory, httpResponseFactory, httpClient)

	logger.PrintInfo("Building middlewares...")
	panicRecoveryMiddleware := middleware.NewPanicRecovery()
	securityCorsMiddleware := middleware.NewSecurityCors(securityCorsService)
	timeoutMiddleware := middleware.NewTimeout()
	limiterMiddleware := middleware.NewLimiter()
	cacheMiddleware := middleware.NewCache(cacheService)

	logger.PrintInfo("Building controllers...")
	staticController := controller.NewStatic(gopenDTO)
	endpointController := controller.NewEndpoint(endpointUseCase)

	return server{
		gopen:                   gopen,
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

func (s server) Start() {
	s.logger.PrintInfo("Starting lister and server...")

	s.buildStaticRoutes()

	s.logger.PrintInfo("Starting to read endpoints to register routes...")
	for _, endpoint := range s.gopen.Endpoints() {
		handles := s.buildEndpointHandles()
		s.router.Handle(s.gopen, &endpoint, handles...)

		lenString := helper.SimpleConvertToString(len(handles))
		s.logger.PrintInfof("Registered route with %s handles: %s", lenString, endpoint.Resume())
	}

	address := fmt.Sprint(":", s.gopen.Port())
	s.logger.PrintInfof("Listening and serving HTTP on %s!", address)

	s.Server = &http.Server{
		Addr:    address,
		Handler: s.router.Engine(),
	}

	fmt.Println()
	fmt.Println()
	s.logger.PrintTitle("LISTEN AND SERVER")

	s.ListenAndServe()
}

func (s server) Shutdown(ctx context.Context) error {
	if helper.IsNil(s.Server) {
		return nil
	}
	return s.Server.Shutdown(ctx)
}

func (s server) buildStaticRoutes() {
	s.logger.PrintInfo("Configuring static routes...")
	formatLog := "Registered route with 5 handles: %s --> \"%s\""

	pingEndpoint := s.buildStaticPingRoute()
	s.logger.PrintInfof(formatLog, pingEndpoint.Method(), pingEndpoint.Path())

	versionEndpoint := s.buildStaticVersionRoute()
	s.logger.PrintInfof(formatLog, versionEndpoint.Method(), versionEndpoint.Path())

	settingsEndpoint := s.buildStaticSettingsRoute()
	s.logger.PrintInfof(formatLog, settingsEndpoint.Method(), settingsEndpoint.Path())
}

func (s server) buildStaticPingRoute() *vo.Endpoint {
	endpoint := vo.NewEndpointStatic("/ping", http.MethodGet)
	s.buildStaticRoute(&endpoint, s.staticController.Ping)
	return &endpoint
}

func (s server) buildStaticVersionRoute() *vo.Endpoint {
	endpoint := vo.NewEndpointStatic("/version", http.MethodGet)
	s.buildStaticRoute(&endpoint, s.staticController.Version)
	return &endpoint
}

func (s server) buildStaticSettingsRoute() *vo.Endpoint {
	endpoint := vo.NewEndpointStatic("/settings", http.MethodGet)
	s.buildStaticRoute(&endpoint, s.staticController.Settings)
	return &endpoint
}

func (s server) buildStaticRoute(endpointStatic *vo.Endpoint, handler app.HandlerFunc) {
	timeoutHandler := s.timeoutMiddleware.Do
	panicHandler := s.panicRecoveryMiddleware.Do
	limiterHandler := s.limiterMiddleware.Do
	s.router.Handle(s.gopen, endpointStatic, timeoutHandler, panicHandler, limiterHandler, handler)
}

func (s server) buildEndpointHandles() []app.HandlerFunc {
	timeoutHandler := s.timeoutMiddleware.Do
	panicHandler := s.panicRecoveryMiddleware.Do
	securityCorsHandler := s.securityCorsMiddleware.Do
	limiterHandler := s.limiterMiddleware.Do
	cacheHandler := s.cacheMiddleware.Do
	endpointHandler := s.endpointController.Execute
	return []app.HandlerFunc{
		timeoutHandler,
		panicHandler,
		securityCorsHandler,
		limiterHandler,
		cacheHandler,
		endpointHandler,
	}
}
