package app

import (
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/controller"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/middleware"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra"
	"github.com/gin-gonic/gin"
	"net/http"
)

var loggerOptions = logger.Options{
	CustomAfterPrefixText: "CMD",
}

var httpServer *http.Server

type gopen struct {
	gopenVO                vo.GOpen
	writerMiddleware       middleware.Writer
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
	writerMiddleware middleware.Writer,
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
		writerMiddleware:       writerMiddleware,
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
	engine := gin.New()

	// configuramos os middlewares estáticos
	g.buildStaticMiddlewares(engine)

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

		// configuramos o handler do security cors como o middleware
		securityCorsHandler := g.securityCorsMiddleware.Do(endpointVO)
		// configuramos o handler do timeout do endpoint como o middleware
		timeoutHandler := g.buildTimeoutMiddlewareHandler(endpointVO)
		// configuramos o handler do limiter do endpoint como o middleware
		limiterHandler := g.buildLimiterMiddlewareHandler(endpointVO)
		// configuramos o handler de cache do endpoint como o middleware
		cacheHandler := g.buildCacheMiddlewareHandler(endpointVO)

		// configuramos o handler do endpoint como controlador
		endpointHandler := g.endpointController.Execute(endpointVO)

		// cadastramos as rotas no gin engine
		engine.Handle(
			endpointVO.Method(),
			endpointVO.Path(),
			securityCorsHandler,
			timeoutHandler,
			limiterHandler,
			cacheHandler,
			endpointHandler,
		)
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

func (g gopen) buildStaticMiddlewares(engine *gin.Engine) {
	// imprimimos o log cmd
	printInfoLog("Configuring static middlewares...")
	// setamos em ordem os middlewares estáticos
	engine.Use(g.writerMiddleware.Do)
	engine.Use(g.traceMiddleware.Do)
	engine.Use(g.logMiddleware.Do)
}

func (g gopen) buildStaticRoutes(engine *gin.Engine) {
	// imprimimos o log cmd
	printInfoLog("Configuring static routes...")
	// ping route
	engine.GET("/ping", g.staticController.Ping)
	// version
	engine.GET("/version", g.staticController.Version)
	// gopen config infos
	engine.GET("/settings", g.staticController.Settings)
}

func (g gopen) buildTimeoutMiddlewareHandler(endpointVO vo.Endpoint) gin.HandlerFunc {
	// por padrão obtemos o timeout configurado na raiz, caso não informado um valor padrão é retornado
	timeoutDuration := g.gopenVO.Timeout()
	// se o timeout foi informado no endpoint damos prioridade a ele
	if endpointVO.HasTimeout() {
		timeoutDuration = endpointVO.Timeout()
	}
	// retornamos o manipulador com o timeout configura
	return g.timeoutMiddleware.Do(endpointVO, timeoutDuration)
}

func (g gopen) buildLimiterMiddlewareHandler(endpointVO vo.Endpoint) gin.HandlerFunc {
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
	return g.limiterMiddleware.Do(endpointVO, rateLimiterProvider, sizeLimiterProvider)
}

func (g gopen) buildCacheMiddlewareHandler(endpointVO vo.Endpoint) gin.HandlerFunc {
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
	return g.cacheMiddleware.Do(g.gopenVO.CacheStore(), endpointVO, cacheVO)
}

func printInfoLog(msg ...any) {
	logger.InfoOpts(loggerOptions, msg...)
}

func printInfoLogf(format string, msg ...any) {
	logger.InfoOptsf(format, loggerOptions, msg...)
}
