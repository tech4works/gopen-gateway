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
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/mapper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/middleware"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/service"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/api"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/boot"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

// NewGopen initializes a new Gopen struct based on the provided GopenJson object.
// It populates the fields of the Gopen struct with the corresponding values from the GopenJson object.
// It also populates the endpoints field by iterating over the EndpointJson objects in the Endpoints slice of the
// GopenJson object, and converting each EndpointJson object to an Endpoint object using the newEndpoint function.
// The newly created Gopen struct is returned as a pointer.
func NewGopen(gopenJson *vo.GopenJson, cacheStore domain.CacheStore) Gopen {
	boot.PrintInfo("Building value objects...")
	gopen := vo.NewGopen(gopenJson)

	boot.PrintInfo("Building infra...")
	loggerProvider := infra.NewLoggerProvider()
	restTemplate := infra.NewRestTemplate(loggerProvider)

	boot.PrintInfo("Building domain...")
	backendService := service.NewBackend(restTemplate)
	endpointService := service.NewEndpoint(backendService)
	cacheService := service.NewCache(cacheStore)

	boot.PrintInfo("Building middlewares...")
	panicRecoveryMiddleware := middleware.NewPanicRecovery(loggerProvider)
	logMiddleware := middleware.NewLog(loggerProvider)
	securityCorsMiddleware := middleware.NewSecurityCors(gopen.SecurityCors())
	timeoutMiddleware := middleware.NewTimeout()
	limiterMiddleware := middleware.NewLimiter()
	cacheMiddleware := middleware.NewCache(cacheService, loggerProvider)

	boot.PrintInfo("Building controllers...")
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

	httpServer = &http.Server{
		Addr:    address,
		Handler: engine,
	}
	_ = httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
// It waits until the context is canceled, all requests are done, or until the timeout is reached.
// If the HTTP server is nil, the method will return nil.
// However, if the server is active, it returns an error resulted from http.Server's Shutdown method.
//
// Returns an error if any occurred during the server shutdown. Returns nil if the server was already nil or shutdown
// executed without errors.
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
	boot.PrintInfo("Configuring static routes...")
	formatLog := "Registered route with 5 handles: %s --> \"%s\""

	pingEndpoint := g.registerStaticPingRoute(engine)
	boot.PrintInfof(formatLog, pingEndpoint.Method(), pingEndpoint.Path())

	versionEndpoint := g.registerStaticVersionRoute(engine)
	boot.PrintInfof(formatLog, versionEndpoint.Method(), versionEndpoint.Path())

	settingsEndpoint := g.registerStaticSettingsRoute(engine)
	boot.PrintInfof(formatLog, settingsEndpoint.Method(), settingsEndpoint.Path())
}

// registerStaticPingRoute is a method of the gopenApp type that registers a static route for the "/ping" path with the
// HTTP method "GET".
// It takes an engine parameter of type *gin.Engine and registers the "/ping" route with the gopen.staticController.Ping
// handler function. It returns a vo.Endpoint pointer representing the registered route.
func (g gopenApp) registerStaticPingRoute(engine *gin.Engine) *vo.Endpoint {
	endpoint := vo.NewEndpointStatic("/ping", http.MethodGet)
	g.registerStaticRoute(engine, &endpoint, g.staticController.Ping)
	return &endpoint
}

// registerStaticVersionRoute is a method of the gopenApp type that registers a static route for the "/version" path
// with the HTTP method "GET".
// It takes an engine parameter of type *gin.Engine and registers the "/version" route with the
// gopen.staticController.Version handler function. It returns a vo.Endpoint pointer representing the registered route.
func (g gopenApp) registerStaticVersionRoute(engine *gin.Engine) *vo.Endpoint {
	endpoint := vo.NewEndpointStatic("/version", http.MethodGet)
	g.registerStaticRoute(engine, &endpoint, g.staticController.Version)
	return &endpoint
}

// registerStaticSettingsRoute is a method of the gopenApp type that registers a static route for the "/settings" path
// with the HTTP method "GET". It takes an engine parameter of type *gin.Engine and registers the "/settings" route with
// the gopen.staticController.Settings handler function. It returns a vo.Endpoint pointer representing the registered route.
func (g gopenApp) registerStaticSettingsRoute(engine *gin.Engine) *vo.Endpoint {
	endpoint := vo.NewEndpointStatic("/settings", http.MethodGet)
	g.registerStaticRoute(engine, &endpoint, g.staticController.Settings)
	return &endpoint
}

// registerStaticRoute is a method of the gopenApp type that registers a static route for a given endpoint with the
// provided handler function. It takes an engine parameter of type *gin.Engine, an endpointStatic parameter of type
// *vo.Endpoint, and a handler parameter of type api.HandlerFunc. The method sets up middleware functions including
// timeoutHandler, panicHandler, logHandler, limiterHandler, and the provided handler function.
// Finally, it calls the api.Handle function passing the engine, g.gopen, endpointStatic, and the middleware functions
// and handler as arguments. This method doesn't return any values.
func (g gopenApp) registerStaticRoute(engine *gin.Engine, endpointStatic *vo.Endpoint, handler api.HandlerFunc) {
	timeoutHandler := g.timeoutMiddleware.Do
	panicHandler := g.panicRecoveryMiddleware.Do
	logHandler := g.logMiddleware.Do
	limiterHandler := g.limiterMiddleware.Do
	api.Handle(engine, g.gopen, endpointStatic, timeoutHandler, panicHandler, logHandler, limiterHandler, handler)
}

// buildEndpointHandles is a method of the gopenApp type that returns a slice of api.HandlerFunc.
// It constructs the slice by assigning each middleware function, as well as the endpointHandler function,
// to a corresponding api.HandlerFunc variable. The slice is then returned.
//
// Returns a slice of api.HandlerFunc, representing the ordered sequence of middleware functions
// and the endpointHandler function.
func (g gopenApp) buildEndpointHandles() []api.HandlerFunc {
	timeoutHandler := g.timeoutMiddleware.Do
	panicHandler := g.panicRecoveryMiddleware.Do
	logHandler := g.logMiddleware.Do
	securityCorsHandler := g.securityCorsMiddleware.Do
	limiterHandler := g.limiterMiddleware.Do
	cacheHandler := g.cacheMiddleware.Do
	endpointHandler := g.endpointController.Execute
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
