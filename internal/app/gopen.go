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
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/controller"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/interfaces"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/service"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

// loggerOptions is a variable that holds the options for the logger package.
// It is of type logger.Options.
// The CustomAfterPrefixText field is used to set a custom text that will be
// printed after the log prefix. In this example, the value is set to "APP".
var loggerOptions = logger.Options{
	HideArgCaller:         true,
	CustomAfterPrefixText: "APP",
}

// httpServer is a variable that holds an instance of the http.Server struct.
// The http.Server struct represents an HTTP server that listens for incoming connections and handles them
// using the specified Handler. It is used to configure and start an HTTP server.
// The httpServer variable is used to register the Gin engine as the handler for the HTTP server in the
// ListerAndServer method of the gopenApp type.
// It can be used to access properties and methods of the HTTP server, such as Shutdown, which gracefully shuts down
// the server without interrupting active connections.
// If the HTTP server is nil, the Shutdown method will return nil.
// Otherwise, it will return an error resulting from the http.Server's Shutdown method.
var httpServer *http.Server

// gopenApp is a struct that holds various components and controllers required for running a Gopen server.
// It contains a gopenApp field that represents the configuration and settings for the Gopen server.
// It also includes middleware implementations such as panicRecoveryMiddleware, traceMiddleware, logMiddleware,
// securityCorsMiddleware, timeoutMiddleware, limiterMiddleware, cacheMiddleware, as well as static and endpoint
// controllers to handle requests.
type gopenApp struct {
	gopen                   *vo.Gopen
	panicRecoveryMiddleware middleware.PanicRecovery
	logMiddleware           middleware.Log
	securityCorsMiddleware  middleware.SecurityCors
	timeoutMiddleware       middleware.Timeout
	limiterMiddleware       middleware.Limiter
	cacheMiddleware         middleware.Cache
	staticController        controller.Static
	endpointController      controller.Endpoint
}

// Gopen is an interface that represents the functionality of a Gopen server.
// It contains the methods ListerAndServer() and Shutdown(ctx context.Context) error
type Gopen interface {
	// ListerAndServer starts the Gopen application with the initialized cache store
	// and Gopen configuration. It builds necessary infrastructures, services,
	// middlewares, controllers, and the Gopen application. The Gopen application is
	// then started by calling its ListAndServer() method.
	ListerAndServer()
	// Shutdown stops the Gopen application gracefully within the given context.
	// It shuts down the server, closes any resources, and cancels any ongoing operations.
	// It returns an error if the shutdown is not successful.
	Shutdown(ctx context.Context) error
}

func NewGopen(gopenJson *vo.GopenJson, cacheStore interfaces.CacheStore) Gopen {
	printInfo("Building value objects...")
	gopen := vo.NewGopen(gopenJson)

	printInfo("Building infra...")
	restTemplate := infra.NewRestTemplate()
	logProvider := infra.NewHttpLoggerProvider()

	printInfo("Building domain...")
	backendService := service.NewBackend(restTemplate)
	endpointService := service.NewEndpoint(backendService)
	cacheService := service.NewCache(cacheStore)

	printInfo("Building middlewares...")
	panicRecoveryMiddleware := middleware.NewPanicRecovery()
	logMiddleware := middleware.NewLog(logProvider)
	securityCorsMiddleware := middleware.NewSecurityCors(gopen.SecurityCors())
	timeoutMiddleware := middleware.NewTimeout()
	limiterMiddleware := middleware.NewLimiter()
	cacheMiddleware := middleware.NewCache(cacheService)

	printInfo("Building controllers...")
	staticController := controller.NewStatic(gopenJson.Version, mapper.BuildSettingViewDTO(gopenJson, gopen))
	endpointController := controller.NewEndpoint(endpointService)

	return gopenApp{
		gopen:                   gopen,
		panicRecoveryMiddleware: panicRecoveryMiddleware,
		logMiddleware:           logMiddleware,
		timeoutMiddleware:       timeoutMiddleware,
		limiterMiddleware:       limiterMiddleware,
		cacheMiddleware:         cacheMiddleware,
		securityCorsMiddleware:  securityCorsMiddleware,
		staticController:        staticController,
		endpointController:      endpointController,
	}
}

// ListerAndServer is a method of the gopen type that sets up and runs an HTTP server.
// It starts by setting the Gin framework's mode to Release and initializes a new Gin engine.
// Then it configures static routes and begins the process of registering routes for each endpoint.
// If an endpoint is already registered, it raises an error.
// After route registration, it constructs the server's address using a configured port,
// and starts an HTTP server listening on the constructed address.
// The server uses the Gin engine as its handler.
// This method doesn't accept parameters or return values.
func (g gopenApp) ListerAndServer() {
	printInfo("Starting lister and server...")

	// instanciamos o gin engine
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	// configuramos rotas estáticas
	g.buildStaticRoutes(engine)

	printInfo("Starting to read endpoints to register routes...")
	// iteramos os endpoints para cadastrar as rotas
	for _, endpoint := range g.gopen.Endpoints() {
		handles := g.buildEndpointHandles()

		api.Handle(engine, g.gopen, &endpoint, handles...)

		lenString := helper.SimpleConvertToString(len(handles))
		printInfof("Registered route with %s handles: %s", lenString, endpoint.Resume())
	}

	// montamos o endereço com a porta configurada
	address := fmt.Sprint(":", g.gopen.Port())

	// rodamos o gin engine
	printInfof("Listening and serving HTTP on %s!", address)

	// construímos o http server do go com o handler gin
	httpServer = &http.Server{
		Addr:    address,
		Handler: engine,
	}

	// chamamos o lister and server
	_ = httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
// It waits until the context is canceled, all requests are done, or until the timeout is reached.
// If the HTTP server is nil, the method will return nil.
// However, if the server is active, it returns an error resulted from http.Server's Shutdown method.
//
// Returns an error if any occurred during the server shutdown. Returns nil if the server was already nil or shutdown executed without errors.
func (g gopenApp) Shutdown(ctx context.Context) error {
	if helper.IsNil(httpServer) {
		return nil
	}
	return httpServer.Shutdown(ctx)
}

// buildStaticRoutes is a method of the gopenApp type that configures static routes for the Gin engine.
// It takes an engine parameter of type *gin.Engine and configures the following routes:
// - "/ping" with the HTTP method "GET" that maps to gopen.staticController.Ping
// - "/version" with the HTTP method "GET" that maps to gopen.staticController.Version
// - "/settings" with the HTTP method "GET" that maps to gopen.staticController.Settings
func (g gopenApp) buildStaticRoutes(engine *gin.Engine) {
	// imprimimos o log cmd
	printInfo("Configuring static routes...")

	// format
	formatLog := "Registered route with 5 handles: %s --> \"%s\""

	// ping route
	pingEndpoint := g.registerStaticPingRoute(engine)
	printInfof(formatLog, pingEndpoint.Method(), pingEndpoint.Path())

	// version
	versionEndpoint := g.registerStaticVersionRoute(engine)
	printInfof(formatLog, versionEndpoint.Method(), versionEndpoint.Path())

	// gopen config infos
	settingsEndpoint := g.registerStaticSettingsRoute(engine)
	printInfof(formatLog, settingsEndpoint.Method(), settingsEndpoint.Path())
}

func (g gopenApp) registerStaticPingRoute(engine *gin.Engine) *vo.Endpoint {
	endpoint := vo.NewEndpointStatic("/ping", http.MethodGet)
	g.registerStaticRoute(engine, &endpoint, g.staticController.Ping)
	return &endpoint
}

func (g gopenApp) registerStaticVersionRoute(engine *gin.Engine) *vo.Endpoint {
	endpoint := vo.NewEndpointStatic("/version", http.MethodGet)
	g.registerStaticRoute(engine, &endpoint, g.staticController.Version)
	return &endpoint
}

func (g gopenApp) registerStaticSettingsRoute(engine *gin.Engine) *vo.Endpoint {
	endpoint := vo.NewEndpointStatic("/settings", http.MethodGet)
	g.registerStaticRoute(engine, &endpoint, g.staticController.Settings)
	return &endpoint
}

func (g gopenApp) registerStaticRoute(engine *gin.Engine, endpointStatic *vo.Endpoint, handler api.HandlerFunc) {
	// configuramos o handler do timeout do endpoint como o middleware
	timeoutHandler := g.timeoutMiddleware.Do
	// configuramos o handler de panic recovery
	panicHandler := g.panicRecoveryMiddleware.Do
	// configuramos o handler do log como o middleware
	logHandler := g.logMiddleware.Do
	// configuramos o handler do limiter do endpoint como o middleware
	limiterHandler := g.limiterMiddleware.Do
	// registramos o endpoint estático
	api.Handle(engine, g.gopen, endpointStatic, timeoutHandler, panicHandler, logHandler, limiterHandler, handler)
}

// buildEndpointHandles is a method of the gopenApp type that builds a list of middleware handlers for a given endpoint.
// It takes an endpointVO of type vo.Endpoint as a parameter and returns a slice of api.HandlerFunc.
// Each middleware handler is configured based on specific middleware instances defined in the gopenApp type.
// The handlers are added to the slice in the following order:
// 1. timeoutHandler: Used to handle timeout requests. The timeout duration is determined based on both the endpointVO
// 2. panicHandler: Used to handle panic errors.
// 3. traceHandler: Used to handle trace requests.
// 4. logHandler: Used to handle logging requests.
// 5. securityCorsHandler: Used to handle security CORS requests.
// and the gopenApp configurations.
// 6. limiterHandler: Used to handle limiter requests. The limiter vo is determined based on both the endpointVO and the
// gopenApp configurations.
// 7. cacheHandler: Used to handle cache requests. The cache duration, cache strategy headers, and allow cache control
// configurations are determined based on both the endpointVO and gopenApp
func (g gopenApp) buildEndpointHandles() []api.HandlerFunc {
	// configuramos o handler do timeout do endpoint como o middleware
	timeoutHandler := g.timeoutMiddleware.Do
	// configuramos o handler de panic recovery
	panicHandler := g.panicRecoveryMiddleware.Do
	// configuramos o handler do log como o middleware
	logHandler := g.logMiddleware.Do
	// configuramos o handler do security cors como o middleware
	securityCorsHandler := g.securityCorsMiddleware.Do
	// configuramos o handler do limiter do endpoint como o middleware
	limiterHandler := g.limiterMiddleware.Do
	// configuramos o handler de cache do endpoint como o middleware
	cacheHandler := g.cacheMiddleware.Do
	// configuramos o handler do endpoint como controlador
	endpointHandler := g.endpointController.Execute
	// montamos a lista de manipuladores
	return []api.HandlerFunc{
		timeoutHandler,
		panicHandler,
		logHandler,
		securityCorsHandler,
		limiterHandler,
		cacheHandler,
		endpointHandler,
	}
}

// printInfo is a function that logs informational messages using logger.InfoOpts from the logger package.
// It accepts a variadic parameter 'msg' of any type that represents the message to be logged.
// The loggerOptions variable is used as the logger's options, which contains custom configuration options.
func printInfo(msg ...any) {
	logger.InfoOpts(loggerOptions, msg...)
}

// printInfof is a function that logs informational messages using logger.InfoOptsf from the logger package.
// It accepts a parameter 'format' of type string that represents the format string for the log message.
// It also accepts a variadic parameter 'msg' of any type that represents the message to be logged.
// The loggerOptions variable is used as the logger's options, which contains custom configuration options.
func printInfof(format string, msg ...any) {
	logger.InfoOptsf(format, loggerOptions, msg...)
}
