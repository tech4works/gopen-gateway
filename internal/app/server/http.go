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

package server

import (
	"context"
	"fmt"
	"net"
	nethttp "net/http"
	"os"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/controller"
	"github.com/tech4works/gopen-gateway/internal/app/factory"
	"github.com/tech4works/gopen-gateway/internal/app/interceptor"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/app/usecase"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"github.com/tech4works/gopen-gateway/internal/domain/service"
)

type http struct {
	net                      *nethttp.Server
	gopen                    *vo.GopenConfig
	log                      app.BootLog
	router                   app.Router
	panicRecoveryInterceptor interceptor.PanicRecovery
	logInterceptor           interceptor.Log
	securityCorsInterceptor  interceptor.SecurityCors
	timeoutInterceptor       interceptor.Timeout
	limiterInterceptor       interceptor.Limiter
	staticController         controller.Static
	endpointController       controller.Endpoint
}

type HTTP interface {
	ListenAndServe()
	Shutdown(ctx context.Context) error
}

func New(
	gopen *dto.Gopen,
	log app.BootLog,
	router app.Router,
	httpClient app.HTTPClient,
	publisherClient app.PublisherClient,
	middlewareLog app.MiddlewareLog,
	endpointLog app.EndpointLog,
	backendLog app.BackendLog,
	httpLog app.HTTPLog,
	jsonPath domain.JSONPath,
	converter domain.Converter,
	store domain.Store,
	nomenclature domain.Nomenclature,
) HTTP {
	log.PrintInfo("Building domain...")
	dynamicValueService := service.NewDynamicValue(jsonPath)
	mapperService := service.NewMapper(jsonPath, dynamicValueService)
	projectorService := service.NewProjector(jsonPath, dynamicValueService)
	modifierService := service.NewModifier(jsonPath, dynamicValueService)
	joinService := service.NewJoin(jsonPath, dynamicValueService)
	omitterService := service.NewOmitter(jsonPath)
	nomenclatureService := service.NewNomenclature(jsonPath, nomenclature)
	contentService := service.NewContent(converter)
	aggregatorService := service.NewAggregator(jsonPath)

	buildPipelineService := service.NewBuildPipeline(modifierService, joinService, mapperService, projectorService,
		omitterService, nomenclatureService, contentService, aggregatorService, dynamicValueService)
	limiterService := service.NewLimiter()
	securityCorsService := service.NewSecurityCors(dynamicValueService)
	cacheService := service.NewCache(dynamicValueService, store)

	log.PrintInfo("Building factories...")
	backendRequestFactory := factory.NewBackendRequest(buildPipelineService)
	backendResponseFactory := factory.NewBackendResponse(buildPipelineService)
	endpointResponseFactory := factory.NewEndpointResponse(aggregatorService, buildPipelineService)

	log.PrintInfo("Building use cases...")
	endpointUseCase := usecase.NewEndpoint(dynamicValueService, cacheService, backendRequestFactory,
		backendResponseFactory, endpointResponseFactory, httpClient, publisherClient, endpointLog, backendLog)

	log.PrintInfo("Building middlewares...")
	panicRecoveryInterceptor := interceptor.NewPanicRecovery(middlewareLog)
	logInterceptor := interceptor.NewLog(httpLog)
	securityCorsInterceptor := interceptor.NewSecurityCors(securityCorsService)
	timeoutInterceptor := interceptor.NewTimeout()
	limiterInterceptor := interceptor.NewLimiter(limiterService)

	log.PrintInfo("Building controllers...")
	staticController := controller.NewStatic(gopen)
	endpointController := controller.NewEndpoint(endpointUseCase)

	log.PrintInfo("Building value objects...")
	return &http{
		gopen:                    factory.BuildGopen(gopen),
		log:                      log,
		router:                   router,
		panicRecoveryInterceptor: panicRecoveryInterceptor,
		logInterceptor:           logInterceptor,
		timeoutInterceptor:       timeoutInterceptor,
		limiterInterceptor:       limiterInterceptor,
		securityCorsInterceptor:  securityCorsInterceptor,
		staticController:         staticController,
		endpointController:       endpointController,
	}
}

func (h *http) ListenAndServe() {
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	h.log.PrintInfo("Configuring routes...")

	h.buildStaticRoutes()
	h.buildRoutes()

	h.net = &nethttp.Server{
		Handler: h.router.Engine(),
	}

	var listener net.Listener
	var err error

	// todo: revisar
	//if h.gopen.HasProxy() {
	//	h.log.PrintInfo("Configuring proxy...")
	//
	//	var opts []config.HTTPEndpointOption
	//	for _, d := range h.gopen.Proxy().Domains() {
	//		opts = append(opts, config.WithDomain(d))
	//	}
	//
	//	listener, err = ngrok.Listen(
	//		ctx,
	//		config.HTTPEndpoint(opts...),
	//		ngrok.WithAuthtoken(h.gopen.Proxy().Token()),
	//	)
	//} else {
	//	listener, err = net.Listen("tcp", fmt.Sprint(":", os.Getenv("PORT")))
	//}

	listener, err = net.Listen("tcp", fmt.Sprint(":", os.Getenv("PORT")))

	if checker.NonNil(err) {
		panic(err)
	}

	h.log.SkipLine()
	h.log.PrintTitle(fmt.Sprintf("LISTEN AND SERVE %s", listener.Addr().String()))

	h.net.Serve(listener)
}

func (h *http) Shutdown(ctx context.Context) error {
	if checker.IsNil(h.net) {
		return nil
	}
	return h.net.Shutdown(ctx)
}

func (h *http) buildRoutes() {
	for _, endpoint := range h.gopen.Endpoints() {
		handles := h.buildEndpointHandles()
		h.router.Handle(h.gopen, &endpoint, handles...)

		h.log.PrintInfof("Registered route with %s handles: %s", converter.ToString(len(handles)), endpoint.Resume())
	}
}

func (h *http) buildStaticRoutes() {
	formatLog := "Registered route with 5 handles: %s --> \"%s\""

	pingEndpoint := h.buildStaticPingRoute()
	h.log.PrintInfof(formatLog, pingEndpoint.Method(), pingEndpoint.Path())

	versionEndpoint := h.buildStaticVersionRoute()
	h.log.PrintInfof(formatLog, versionEndpoint.Method(), versionEndpoint.Path())
}

func (h *http) buildStaticPingRoute() *vo.EndpointConfig {
	endpoint := vo.NewEndpointConfigStatic("/ping", nethttp.MethodGet)
	h.buildStaticRoute(&endpoint, h.staticController.Ping)
	return &endpoint
}

func (h *http) buildStaticVersionRoute() *vo.EndpointConfig {
	endpoint := vo.NewEndpointConfigStatic("/version", nethttp.MethodGet)
	h.buildStaticRoute(&endpoint, h.staticController.Version)
	return &endpoint
}

func (h *http) buildStaticRoute(endpointStatic *vo.EndpointConfig, handler app.HandlerFunc) {
	timeoutHandler := h.timeoutInterceptor.Do
	panicHandler := h.panicRecoveryInterceptor.Do
	logHandler := h.logInterceptor.Do
	limiterHandler := h.limiterInterceptor.Do
	h.router.Handle(h.gopen, endpointStatic, timeoutHandler, panicHandler, logHandler, limiterHandler, handler)
}

func (h *http) buildEndpointHandles() []app.HandlerFunc {
	return []app.HandlerFunc{
		h.panicRecoveryInterceptor.Do,
		h.timeoutInterceptor.Do,
		h.logInterceptor.Do,
		h.securityCorsInterceptor.Do,
		h.limiterInterceptor.Do,
		h.endpointController.Do,
	}
}
