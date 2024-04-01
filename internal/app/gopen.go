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

func (g gopen) Shutdown(ctx context.Context) error {
	if helper.IsNil(httpServer) {
		return nil
	}
	return httpServer.Shutdown(ctx)
}

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
