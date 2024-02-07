package application

import (
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/controller"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/middleware"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/model/dto"
	"github.com/GabrielHCataldo/martini-gateway/internal/infra"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type gateway struct {
	martini            dto.Martini
	redisStore         *persist.RedisStore
	timeoutMiddleware  middleware.Timeout
	limiterMiddleware  middleware.Limiter
	corsMiddleware     middleware.Cors
	endpointController controller.Endpoint
}

type Gateway interface {
	Run()
}

func NewGateway(
	martini dto.Martini,
	redisStore *persist.RedisStore,
	timeoutMiddleware middleware.Timeout,
	limiterMiddleware middleware.Limiter,
	corsMiddleware middleware.Cors,
	endpointController controller.Endpoint,
) Gateway {
	return gateway{
		martini:            martini,
		redisStore:         redisStore,
		timeoutMiddleware:  timeoutMiddleware,
		limiterMiddleware:  limiterMiddleware,
		corsMiddleware:     corsMiddleware,
		endpointController: endpointController,
	}
}

func (a gateway) Run() {
	logger.Info("Starting gateway application!")

	ginEngine := gin.New()

	logger.Info("Configuring middleware!")
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(gin.Logger())
	ginEngine.Use(a.timeoutMiddleware.PreHandlerRequest)
	ginEngine.Use(a.limiterMiddleware.PreHandlerRequest)
	ginEngine.Use(a.corsMiddleware.PreHandlerRequest)

	logger.Info("Setting cache duration by martini.cache!")
	cacheDuration, err := time.ParseDuration(a.martini.Cache)
	if helper.IsNotNil(err) {
		panic(errors.New("Error parse config.cache duration:", err))
	}

	logger.Info("Starting route configuration!")
	ginEngine.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "%s", "Pong!")
	})
	for _, endpoint := range a.martini.Endpoints {
		if helper.NotContains(a.martini.ExtraConfig.SecurityCors.AllowMethods, endpoint.Method) {
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
		cacheGinHandler := infra.CacheHandler(a.redisStore, cacheDurationEndpoint)
		switch endpoint.Method {
		case http.MethodPost:
			ginEngine.POST(endpoint.Endpoint, cacheGinHandler, a.endpointController.Execute)
			break
		case http.MethodGet:
			ginEngine.GET(endpoint.Endpoint, cacheGinHandler, a.endpointController.Execute)
			break
		case http.MethodPut:
			ginEngine.PUT(endpoint.Endpoint, cacheGinHandler, a.endpointController.Execute)
			break
		case http.MethodPatch:
			ginEngine.PATCH(endpoint.Endpoint, cacheGinHandler, a.endpointController.Execute)
			break
		case http.MethodDelete:
			ginEngine.DELETE(endpoint.Endpoint, cacheGinHandler, a.endpointController.Execute)
			break
		default:
			panic(errors.New("Error method: ", endpoint.Method, "not found on valid methods (POST, GET, PUT, PATCH, DELETE)"))
		}
	}

	if helper.IsEmpty(a.martini.Limiter.MaxSizeRequestHeader) {
		a.martini.Limiter.MaxSizeRequestHeader = "1MB"
	}
	if helper.IsEmpty(a.martini.Limiter.MaxSizeMultipartMemory) {
		a.martini.Limiter.MaxSizeMultipartMemory = "5MB"
	}
	maxSizeRequestHeader, err := helper.ConvertMegaByteUnit(a.martini.Limiter.MaxSizeRequestHeader)
	if err != nil {
		panic(errors.New("Error parse megabyte unit limiter.maxSizeRequestHeader field:", err))
		return
	}

	var readHeaderTimeout time.Duration
	var readTimeout time.Duration
	var writeTimeout time.Duration
	if helper.IsNotEmpty(a.martini.Timeout.ReadHeader) {
		readHeaderTimeout, err = time.ParseDuration(a.martini.Timeout.ReadHeader)
		if helper.IsNotNil(err) {
			panic(errors.New("Error parse duration timeout.readHeader field:", err))
		}
	}
	if helper.IsNotEmpty(a.martini.Timeout.Read) {
		readTimeout, err = time.ParseDuration(a.martini.Timeout.Read)
		if helper.IsNotNil(err) {
			panic(errors.New("Error parse duration timeout.read field:", err))
		}
	}
	if helper.IsNotEmpty(a.martini.Timeout.Write) {
		writeTimeout, err = time.ParseDuration(a.martini.Timeout.Write)
		if helper.IsNotNil(err) {
			panic(errors.New("Error parse duration timeout.write field:", err))
		}
	}

	address := fmt.Sprint(":", a.martini.Port)
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
