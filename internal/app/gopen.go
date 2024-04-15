package app

import (
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/controller"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
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
	CustomAfterPrefixText: "APP",
}

// httpServer is a variable that holds an instance of the http.Server struct.
// The http.Server struct represents an HTTP server that listens for incoming connections and handles them
// using the specified Handler. It is used to configure and start an HTTP server.
// The httpServer variable is used to register the Gin engine as the handler for the HTTP server in the
// ListerAndServer method of the gopen type.
// It can be used to access properties and methods of the HTTP server, such as Shutdown, which gracefully shuts down
// the server without interrupting active connections.
// If the HTTP server is nil, the Shutdown method will return nil.
// Otherwise, it will return an error resulting from the http.Server's Shutdown method.
var httpServer *http.Server

// gopen is a struct that holds various components and controllers required for running a Gopen server.
// It contains a gopenVO field that represents the configuration and settings for the Gopen server.
// It also includes middleware implementations such as traceMiddleware, logMiddleware, securityCorsMiddleware,
// timeoutMiddleware, limiterMiddleware, cacheMiddleware, as well as static and endpoint controllers to handle requests.
type gopen struct {
	gopenVO                vo.Gopen
	traceMiddleware        middleware.Trace
	logMiddleware          middleware.Log
	securityCorsMiddleware middleware.SecurityCors
	timeoutMiddleware      middleware.Timeout
	limiterMiddleware      middleware.Limiter
	cacheMiddleware        middleware.Cache
	staticController       controller.Static
	endpointController     controller.Endpoint
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

// NewGopen creates and returns a new `Gopen` object.
// It returns a `Gopen` interface, which represents the Gopen object that stores the provided configuration and middleware.
func NewGopen(
	gopenVO vo.Gopen,
	traceMiddleware middleware.Trace,
	logMiddleware middleware.Log,
	securityCorsMiddleware middleware.SecurityCors,
	timeoutMiddleware middleware.Timeout,
	limiterMiddleware middleware.Limiter,
	cacheMiddleware middleware.Cache,
	staticController controller.Static,
	endpointController controller.Endpoint,
) Gopen {
	return gopen{
		gopenVO:                gopenVO,
		traceMiddleware:        traceMiddleware,
		logMiddleware:          logMiddleware,
		timeoutMiddleware:      timeoutMiddleware,
		limiterMiddleware:      limiterMiddleware,
		cacheMiddleware:        cacheMiddleware,
		securityCorsMiddleware: securityCorsMiddleware,
		staticController:       staticController,
		endpointController:     endpointController,
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
func (g gopen) ListerAndServer() {
	printInfoLog("Starting lister and server...")

	// instanciamos o gin engine
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// configuramos rotas estáticas
	g.buildStaticRoutes(engine)

	printInfoLog("Starting to read endpoints to register routes...")
	// iteramos os endpoints para cadastrar as rotas
	for _, endpointVO := range g.gopenVO.Endpoints() {
		// verificamos se ja existe esse endpoint cadastrado
		for _, route := range engine.Routes() {
			if err := endpointVO.Equals(route); helper.IsNotNil(err) {
				panic(err)
			}
		}

		// configuramos os handles do endpoint
		handles := g.buildEndpointHandles(endpointVO)

		// cadastramos as rotas no nosso wrapper
		api.Handle(engine, g.gopenVO, endpointVO, handles...)

		// imprimimos a informação dos endpoints cadastrado
		printInfoLogf("registered route %s", endpointVO.Resume())
	}

	// montamos o endereço com a porta configurada
	address := fmt.Sprint("127.0.0.1:", g.gopenVO.Port())

	// rodamos o gin engine
	printInfoLogf("Listening and serving HTTP on %s!", address)

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
func (g gopen) Shutdown(ctx context.Context) error {
	if helper.IsNil(httpServer) {
		return nil
	}
	return httpServer.Shutdown(ctx)
}

// buildStaticRoutes is a method of the gopen type that configures static routes for the Gin engine.
// It takes an engine parameter of type *gin.Engine and configures the following routes:
// - "/ping" with the HTTP method "GET" that maps to gopen.staticController.Ping
// - "/version" with the HTTP method "GET" that maps to gopen.staticController.Version
// - "/settings" with the HTTP method "GET" that maps to gopen.staticController.Settings
func (g gopen) buildStaticRoutes(engine *gin.Engine) {
	// imprimimos o log cmd
	printInfoLog("Configuring static routes...")

	// format
	formatLog := "registered route %s -> \"%s\""

	// ping route
	pingMethod := http.MethodGet
	pingPath := "/ping"
	engine.Handle(pingMethod, pingPath, g.staticController.Ping)
	printInfoLogf(formatLog, pingMethod, pingPath)

	// version
	versionMethod := http.MethodGet
	versionPath := "/version"
	engine.Handle(versionMethod, versionPath, g.staticController.Version)
	printInfoLogf(formatLog, versionMethod, versionPath)

	// gopen config infos
	settingsMethod := http.MethodGet
	settingsPath := "/settings"
	engine.Handle(settingsMethod, settingsPath, g.staticController.Settings)
	printInfoLogf(formatLog, settingsMethod, settingsPath)
}

// buildEndpointHandles is a method of the gopen type that builds a list of middleware handlers for a given endpoint.
// It takes an endpointVO of type vo.Endpoint as a parameter and returns a slice of api.HandlerFunc.
// Each middleware handler is configured based on specific middleware instances defined in the gopen type.
// The handlers are added to the slice in the following order:
// 1. traceHandler: Used to handle trace requests.
// 2. logHandler: Used to handle logging requests.
// 3. securityCorsHandler: Used to handle security CORS requests.
// 4. timeoutHandler: Used to handle timeout requests. The timeout duration is determined based on both the endpointVO
// and the gopenVO configurations.
// 5. limiterHandler: Used to handle limiter requests. The limiter vo is determined based on both the endpointVO and the
// gopenVO configurations.
// 6. cacheHandler: Used to handle cache requests. The cache duration, cache strategy headers, and allow cache control
// configurations are determined based on both the endpointVO and gopenVO
func (g gopen) buildEndpointHandles(endpointVO vo.Endpoint) []api.HandlerFunc {
	// configuramos o handler do log como o middleware
	logHandler := g.logMiddleware.Do
	// configuramos o handler do trace como o middleware
	traceHandler := g.traceMiddleware.Do
	// configuramos o handler do security cors como o middleware
	securityCorsHandler := g.securityCorsMiddleware.Do
	// configuramos o handler do timeout do endpoint como o middleware
	timeoutHandler := g.timeoutMiddleware.Do(endpointVO.Timeout())
	// configuramos o handler do limiter do endpoint como o middleware
	limiterHandler := g.buildLimiterMiddlewareHandler(endpointVO)
	// configuramos o handler de cache do endpoint como o middleware
	cacheHandler := g.cacheMiddleware.Do(endpointVO.Cache())
	// configuramos o handler do endpoint como controlador
	endpointHandler := g.endpointController.Execute
	// montamos a lista de manipuladores
	return []api.HandlerFunc{
		traceHandler,
		logHandler,
		securityCorsHandler,
		timeoutHandler,
		limiterHandler,
		cacheHandler,
		endpointHandler,
	}
}

// buildLimiterMiddlewareHandler is a method of the gopen type that constructs and returns a limiter middleware handler
// for the given endpoint. It first retrieves the limiter rate and capacity values from the gopenVO configuration. If
// these values are specified in the endpoint, they take priority. Next, it retrieves the max header size, max body size,
// and max multipart form size values from the gopenVO configuration. If these values are specified in the endpoint,
// they take.
func (g gopen) buildLimiterMiddlewareHandler(endpointVO vo.Endpoint) api.HandlerFunc {
	// inicializamos o limitador de taxa
	rateLimiterProvider := infra.NewRateLimiterProvider(endpointVO.Limiter().Rate())
	// inicializamos o limitador de tamanho
	sizeLimiterProvider := infra.NewSizeLimiterProvider(endpointVO.Limiter())

	// construímos a chamada limiter
	return g.limiterMiddleware.Do(rateLimiterProvider, sizeLimiterProvider)
}

// printInfoLog is a function that logs informational messages using logger.InfoOpts from the logger package.
// It accepts a variadic parameter 'msg' of any type that represents the message to be logged.
// The loggerOptions variable is used as the logger's options, which contains custom configuration options.
func printInfoLog(msg ...any) {
	logger.InfoOpts(loggerOptions, msg...)
}

// printInfoLogf is a function that logs informational messages using logger.InfoOptsf from the logger package.
// It accepts a parameter 'format' of type string that represents the format string for the log message.
// It also accepts a variadic parameter 'msg' of any type that represents the message to be logged.
// The loggerOptions variable is used as the logger's options, which contains custom configuration options.
func printInfoLogf(format string, msg ...any) {
	logger.InfoOptsf(format, loggerOptions, msg...)
}
