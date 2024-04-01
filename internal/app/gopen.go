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

var loggerOptions = logger.Options{
	CustomAfterPrefixText: "APP",
}

var httpServer *http.Server

type gopen struct {
	gopenVO                vo.GOpen
	traceMiddleware        middleware.Trace
	logMiddleware          middleware.Log
	securityCorsMiddleware middleware.SecurityCors
	timeoutMiddleware      middleware.Timeout
	limiterMiddleware      middleware.Limiter
	cacheMiddleware        middleware.Cache
	staticController       controller.Static
	endpointController     controller.Endpoint
}

type GOpen interface {
	ListerAndServer()
	Shutdown(ctx context.Context) error
}

func NewGOpen(
	gopenVO vo.GOpen,
	traceMiddleware middleware.Trace,
	logMiddleware middleware.Log,
	securityCorsMiddleware middleware.SecurityCors,
	timeoutMiddleware middleware.Timeout,
	limiterMiddleware middleware.Limiter,
	cacheMiddleware middleware.Cache,
	staticController controller.Static,
	endpointController controller.Endpoint,
) GOpen {
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
	address := fmt.Sprint(":", g.gopenVO.Port())

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
	timeoutHandler := g.buildTimeoutMiddlewareHandler(endpointVO)
	// configuramos o handler do limiter do endpoint como o middleware
	limiterHandler := g.buildLimiterMiddlewareHandler(endpointVO)
	// configuramos o handler de cache do endpoint como o middleware
	cacheHandler := g.buildCacheMiddlewareHandler(endpointVO)
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

// buildTimeoutMiddlewareHandler is a method of the gopen type that builds a timeout middleware handler for a given endpoint.
// It accepts an endpointVO of type vo.Endpoint and returns an api.HandlerFunc.
// The method starts by obtaining the timeout duration configured in the gopenVO.
// If the endpointVO has its own timeout duration, it overrides the default value.
// Finally, it returns the timeout middleware handler with the configured timeout duration.
func (g gopen) buildTimeoutMiddlewareHandler(endpointVO vo.Endpoint) api.HandlerFunc {
	// por padrão obtemos o timeout configurado na raiz, caso não informado um valor padrão é retornado
	timeoutDuration := g.gopenVO.Timeout()
	// se o timeout foi informado no endpoint damos prioridade a ele
	if endpointVO.HasTimeout() {
		timeoutDuration = endpointVO.Timeout()
	}
	// retornamos o manipulador com o timeout configura
	return g.timeoutMiddleware.Do(timeoutDuration)
}

// buildLimiterMiddlewareHandler is a method of the gopen type that constructs and returns a limiter middleware handler
// for the given endpoint. It first retrieves the limiter rate and capacity values from the gopenVO configuration. If
// these values are specified in the endpoint, they take priority. Next, it retrieves the max header size, max body size,
// and max multipart form size values from the gopenVO configuration. If these values are specified in the endpoint,
// they take.
func (g gopen) buildLimiterMiddlewareHandler(endpointVO vo.Endpoint) api.HandlerFunc {
	// por padrão obtemos o limiter.rate.every configurado na raiz, caso não informado um valor padrão é retornado
	rateEvery := g.gopenVO.LimiterRateEvery()
	// caso informado no endpoint damos prioridade
	if endpointVO.HasLimiterRateEvery() {
		rateEvery = endpointVO.LimiterRateEvery()
	}

	// por padrão obtemos o limiter.rate.capacity configurado na raiz, caso não informado um valor padrão é retornado
	rateCapacity := g.gopenVO.LimiterRateCapacity()
	// caso informado no endpoint damos prioridade
	if endpointVO.HasLimiterRateCapacity() {
		rateCapacity = endpointVO.LimiterRateCapacity()
	}

	// por padrão obtemos o limiter.max-header-size configurado na raiz, caso não informado um valor padrão é retornado
	maxHeaderSize := g.gopenVO.LimiterMaxHeaderSize()
	// caso informado no endpoint damos prioridade
	if endpointVO.HasLimiterMaxHeaderSize() {
		maxHeaderSize = endpointVO.LimiterMaxHeaderSize()
	}

	// por padrão obtemos o limiter.max-body-size configurado na raiz, caso não informado um valor padrão é retornado
	maxBodySize := g.gopenVO.LimiterMaxBodySize()
	// caso informado no endpoint damos prioridade
	if endpointVO.HasLimiterMaxBodySize() {
		maxBodySize = endpointVO.LimiterMaxBodySize()
	}

	// por padrão obtemos o limiter.max-multipart-form-size configurado na raiz, caso não informado um valor padrão é retornado
	maxMultipartForm := g.gopenVO.LimiterMaxMultipartMemorySize()
	// caso informado no endpoint damos prioridade
	if endpointVO.HasLimiterMaxMultipartFormSize() {
		maxMultipartForm = endpointVO.LimiterMaxMultipartMemorySize()
	}

	// inicializamos o limitador de taxa
	rateLimiterProvider := infra.NewRateLimiterProvider(rateEvery, rateCapacity)
	// inicializamos o limitador de tamanho
	sizeLimiterProvider := infra.NewSizeLimiterProvider(maxHeaderSize, maxBodySize, maxMultipartForm)

	// construímos a chamada limiter
	return g.limiterMiddleware.Do(rateLimiterProvider, sizeLimiterProvider)
}

// buildCacheMiddlewareHandler is a method of the gopen type that builds and returns a cache middleware handler for the
// given endpointVO. It first obtains the cache duration value from the parent gopenVO. If the endpointVO has a cache
// duration value, it overrides the parent's value. Next, it obtains the cache strategy headers value from the parent
// gopenVO. If the endpointVO has a cache strategy headers value, it overrides the parent's value. Then, it obtains the
// allow cache control value from the parent gopenVO. If the endpointVO has an allow cache control value, it overrides
// the parent's value. Using these values, it creates a cacheVO value object using the vo.NewCacheFromEndpoint function.
// Finally, it returns the cache middleware handler by calling g.cacheMiddleware.Do with the cacheVO object as the argument.
func (g gopen) buildCacheMiddlewareHandler(endpointVO vo.Endpoint) api.HandlerFunc {
	// obtemos o valor do pai
	cacheDuration := g.gopenVO.CacheDuration()
	// caso seja informado no endpoint, damos prioridade
	if endpointVO.HasCacheDuration() {
		cacheDuration = endpointVO.CacheDuration()
	}
	// obtemos o valor do pai
	cacheStrategyHeaders := g.gopenVO.CacheStrategyHeaders()
	// caso seja informado no endpoint, damos prioridade
	if endpointVO.HasCacheStrategyHeaders() {
		cacheStrategyHeaders = endpointVO.CacheStrategyHeaders()
	}

	// obtemos o valor do pai
	allowCacheControl := g.gopenVO.AllowCacheControl()
	// caso seja informado no endpoint, damos prioridade
	if endpointVO.HasAllowCacheControl() {
		allowCacheControl = endpointVO.AllowCacheControl()
	}

	// com esses valores, construímos o objeto de valor
	cacheVO := vo.NewCacheFromEndpoint(cacheDuration, cacheStrategyHeaders, allowCacheControl)

	// construímos a chamada de cache middleware para o endpoint
	return g.cacheMiddleware.Do(cacheVO)
}

func printInfoLog(msg ...any) {
	logger.InfoOpts(loggerOptions, msg...)
}

func printInfoLogf(format string, msg ...any) {
	logger.InfoOptsf(format, loggerOptions, msg...)
}
