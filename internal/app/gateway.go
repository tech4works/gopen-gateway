package app

import (
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	middleware2 "github.com/GabrielHCataldo/gopen-gateway/internal/app/middleware"
	"github.com/GabrielHCataldo/gopen-gateway/internal/application/controller"
	"github.com/GabrielHCataldo/gopen-gateway/internal/application/middleware"
	"github.com/GabrielHCataldo/gopen-gateway/internal/application/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type gateway struct {
	martini            dto.Martini
	redisStore         *persist.RedisStore
	headerMiddleware   middleware.Header
	logMiddleware      middleware.Log
	timeoutMiddleware  middleware2.Timeout
	limiterMiddleware  middleware.Limiter
	corsMiddleware     middleware2.SecurityCors
	endpointController controller.Endpoint
}

type Gateway interface {
	Run()
}

func NewGateway(
	martini dto.Martini,
	redisStore *persist.RedisStore,
	headerMiddleware middleware.Header,
	logMiddleware middleware.Log,
	timeoutMiddleware middleware2.Timeout,
	limiterMiddleware middleware.Limiter,
	corsMiddleware middleware2.SecurityCors,
	endpointController controller.Endpoint,
) Gateway {
	return gateway{
		martini:            martini,
		redisStore:         redisStore,
		headerMiddleware:   headerMiddleware,
		logMiddleware:      logMiddleware,
		timeoutMiddleware:  timeoutMiddleware,
		limiterMiddleware:  limiterMiddleware,
		corsMiddleware:     corsMiddleware,
		endpointController: endpointController,
	}
}

func (g gateway) Run() {
	logger.Info("Starting gateway application!")

	ginEngine := gin.New()

	logger.Info("Configuring middlewares!")
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(g.headerMiddleware.PreHandlerRequest)
	ginEngine.Use(g.logMiddleware.PreHandlerRequest)
	ginEngine.Use(g.timeoutMiddleware.PreHandlerRequest)
	ginEngine.Use(g.limiterMiddleware.PreHandlerRequest)
	ginEngine.Use(g.corsMiddleware.PreHandlerRequest)

	logger.Info("Setting cache duration by martini.cache!")
	cacheDuration, err := time.ParseDuration(g.martini.Cache)
	if helper.IsNotNil(err) {
		panic(errors.New("Error parse config.cache duration:", err))
	}

	logger.Info("Starting route configuration!")
	ginEngine.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "%s", "Pong!")
	})
	for _, endpoint := range g.martini.Endpoints {
		if helper.NotContains(g.martini.ExtraConfig.SecurityCors.AllowMethods, endpoint.Method) {
			panic(errors.New("Error method:", endpoint.Method, "not allowed on security-cors allow-methods"))
		}
		for _, route := range ginEngine.Routes() {
			if helper.Equals(route.Path, endpoint.Endpoint) && helper.Equals(route.Method, endpoint.Method) {
				panic(errors.New(
					"Error endpoint:", endpoint.Endpoint, "method:", endpoint.Method, "repeat route endpoint error",
				))
			}
		}

		var cacheDurationEndpoint time.Duration
		if endpoint.Cacheable {
			cacheDurationEndpoint = cacheDuration
		}
		cacheGinHandler := infra.CacheHandler(g.redisStore, cacheDurationEndpoint)
		switch endpoint.Method {
		case http.MethodPost:
			ginEngine.POST(endpoint.Endpoint, cacheGinHandler, g.endpointController.Execute)
			break
		case http.MethodGet:
			ginEngine.GET(endpoint.Endpoint, cacheGinHandler, g.endpointController.Execute)
			break
		case http.MethodPut:
			ginEngine.PUT(endpoint.Endpoint, cacheGinHandler, g.endpointController.Execute)
			break
		case http.MethodPatch:
			ginEngine.PATCH(endpoint.Endpoint, cacheGinHandler, g.endpointController.Execute)
			break
		case http.MethodDelete:
			ginEngine.DELETE(endpoint.Endpoint, cacheGinHandler, g.endpointController.Execute)
			break
		default:
			panic(errors.New("Error method: ", endpoint.Method, "not found on valid methods (POST, GET, PUT, PATCH, DELETE)"))
		}
	}

	if helper.IsEmpty(g.martini.Limiter.MaxSizeRequestHeader) {
		g.martini.Limiter.MaxSizeRequestHeader = "1MB"
	}
	if helper.IsEmpty(g.martini.Limiter.MaxSizeMultipartMemory) {
		g.martini.Limiter.MaxSizeMultipartMemory = "5MB"
	}
	maxSizeRequestHeader, err := helper.ConvertMegaByteUnit(g.martini.Limiter.MaxSizeRequestHeader)
	if err != nil {
		panic(errors.New("Error parse megabyte unit limiter.maxSizeRequestHeader field:", err))
		return
	}

	var readHeaderTimeout time.Duration
	var readTimeout time.Duration
	var writeTimeout time.Duration
	if helper.IsNotEmpty(g.martini.Timeout.ReadHeader) {
		readHeaderTimeout, err = time.ParseDuration(g.martini.Timeout.ReadHeader)
		if helper.IsNotNil(err) {
			panic(errors.New("Error parse duration timeout.readHeader field:", err))
		}
	}
	if helper.IsNotEmpty(g.martini.Timeout.Read) {
		readTimeout, err = time.ParseDuration(g.martini.Timeout.Read)
		if helper.IsNotNil(err) {
			panic(errors.New("Error parse duration timeout.read field:", err))
		}
	}
	if helper.IsNotEmpty(g.martini.Timeout.Write) {
		writeTimeout, err = time.ParseDuration(g.martini.Timeout.Write)
		if helper.IsNotNil(err) {
			panic(errors.New("Error parse duration timeout.write field:", err))
		}
	}

	address := fmt.Sprint(":", g.martini.Port)
	s := &http.Server{
		Addr:              address,
		Handler:           ginEngine,
		MaxHeaderBytes:    maxSizeRequestHeader,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
	}
	logger.Info("Listening and serving HTTP on", address)
	err = s.ListenAndServe()
	if helper.IsNotNil(err) {
		panic(errors.New("Error start gateway listen and serve on address:", address, "err:", err))
	}
	logger.Info("Application started!")
}
