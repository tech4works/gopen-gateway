package application

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/dto"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/middleware"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/usecase"
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"slices"
	"strconv"
	"time"
)

type gateway struct {
	config            dto.Config
	redisClient       *redis.Client
	timeoutMiddleware middleware.Timeout
	limiterMiddleware middleware.Limiter
	corsMiddleware    middleware.Cors
	endpointUseCase   usecase.Endpoint
}

type Gateway interface {
	Run()
}

func NewGateway(
	config dto.Config,
	redisClient *redis.Client,
	timeoutMiddleware middleware.Timeout,
	limiterMiddleware middleware.Limiter,
	corsMiddleware middleware.Cors,
	endpointUseCase usecase.Endpoint,
) Gateway {
	return gateway{
		redisClient:       redisClient,
		config:            config,
		timeoutMiddleware: timeoutMiddleware,
		limiterMiddleware: limiterMiddleware,
		corsMiddleware:    corsMiddleware,
		endpointUseCase:   endpointUseCase,
	}
}

func (a gateway) Run() {
	ginRoute := gin.New()
	ginRoute.Use(gin.Recovery())
	ginRoute.Use(gin.Logger())
	ginRoute.Use(a.limiterMiddleware.PreHandlerRequest)
	ginRoute.Use(a.timeoutMiddleware.PreHandlerRequest)
	ginRoute.Use(a.corsMiddleware.PreHandlerRequest)

	ginRoute.GET("/ping", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	cacheDuration, err := time.ParseDuration(a.config.Cache)
	if err != nil {
		logger.Error(err)
		return
	}
	memoryStore := persist.NewRedisStore(a.redisClient)

	for _, endpoint := range a.config.Endpoints {
		if !slices.Contains(a.config.ExtraConfig.SecurityCors.AllowMethods, endpoint.Method) {
			logger.Error("Error method:", endpoint.Method, "not allowed on security-cors allow-methods")
			return
		}
		for _, item := range ginRoute.Routes() {
			if item.Path == endpoint.Endpoint && item.Method == endpoint.Method {
				logger.Error("Error endpoint:", endpoint.Endpoint, "method:", endpoint.Method,
					"repeat route endpoint error")
				return
			}
		}

		cacheDurationEndpoint := 5 * time.Second
		if endpoint.Cacheable {
			cacheDurationEndpoint = cacheDuration
		}

		cacheGinHandler := cache.Cache(
			memoryStore,
			cacheDurationEndpoint,
			cache.WithCacheStrategyByRequest(func(ctx *gin.Context) (bool, cache.Strategy) {
				path := ctx.Request.RequestURI
				method := ctx.Request.Method
				device := ctx.GetHeader("Device")
				ip := ctx.ClientIP()
				key := ip
				if helper.IsNotEmpty(device) {
					key = key + ":" + device
				}
				key = key + ":" + path + ":" + method
				withoutCacheHeader := ctx.GetHeader("Without-Cache")

				shouldCache := true
				if helper.IsNotEmpty(withoutCacheHeader) {
					withoutCache, _ := strconv.ParseBool(withoutCacheHeader)
					shouldCache = withoutCache && method == "GET"
				}
				return shouldCache, cache.Strategy{CacheKey: key}
			}),
			cache.WithBeforeReplyWithCache(func(ctx *gin.Context, cache *cache.ResponseCache) {
				ctx.Header("X-Gateway-Cache", "true")
			}),
		)

		switch endpoint.Method {
		case http.MethodPost:
			ginRoute.POST(endpoint.Endpoint, cacheGinHandler, a.endpointUseCase.Execute)
			break
		case http.MethodGet:
			ginRoute.GET(endpoint.Endpoint, cacheGinHandler, a.endpointUseCase.Execute)
			break
		case http.MethodPut:
			ginRoute.PUT(endpoint.Endpoint, cacheGinHandler, a.endpointUseCase.Execute)
			break
		case http.MethodPatch:
			ginRoute.PATCH(endpoint.Endpoint, cacheGinHandler, a.endpointUseCase.Execute)
			break
		case http.MethodDelete:
			ginRoute.DELETE(endpoint.Endpoint, cacheGinHandler, a.endpointUseCase.Execute)
			break
		default:
			logger.Error("Error method: ", endpoint.Method, "not found on valid methods (POST,GET,PUT,PATCH,DELETE)")
			return
		}
	}

	if helper.IsEmpty(a.config.Limiter.MaxSizeRequestHeader) {
		a.config.Limiter.MaxSizeRequestHeader = "1MB"
	}
	if helper.IsEmpty(a.config.Limiter.MaxSizeMultipartMemory) {
		a.config.Limiter.MaxSizeMultipartMemory = "5MB"
	}
	maxSizeRequestHeader, err := helper.ConvertMegaByteUnit(a.config.Limiter.MaxSizeRequestHeader)
	if err != nil {
		logger.Errorf("Error parse megabyte unit limiter.maxSizeRequestHeader field: %s", err)
		return
	}

	readHeaderTimeout := 0 * time.Second
	readTimeout := 0 * time.Second
	writeTimeout := 0 * time.Second
	if helper.IsNotEmpty(a.config.Timeout.ReadHeader) {
		readHeaderTimeout, err = time.ParseDuration(a.config.Timeout.ReadHeader)
		if err != nil {
			logger.Errorf("Error parse duration timeout.readHeader field: %s", err)
			return
		}
	}
	if helper.IsNotEmpty(a.config.Timeout.Read) {
		readTimeout, err = time.ParseDuration(a.config.Timeout.Read)
		if err != nil {
			logger.Errorf("Error parse duration timeout.read field: %s", err)
			return
		}
	}
	if helper.IsNotEmpty(a.config.Timeout.Write) {
		writeTimeout, err = time.ParseDuration(a.config.Timeout.Write)
		if err != nil {
			logger.Errorf("Error parse duration timeout.write field: %s", err)
			return
		}
	}

	address := ":" + strconv.Itoa(a.config.Port)
	s := &http.Server{
		Addr:              address,
		Handler:           ginRoute,
		MaxHeaderBytes:    maxSizeRequestHeader,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
	}
	logger.Info("Listening and serving HTTP on", address)
	err = s.ListenAndServe()
	if err != nil {
		logger.Errorf("Error start gateway listen and serve on address %s err: %v", address, err)
	}
}
