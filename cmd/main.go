package main

import (
	"encoding/json"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/martini-gateway/internal/application"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/controller"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/middleware"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/model/dto"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/usecase"
	"github.com/GabrielHCataldo/martini-gateway/internal/domain/service"
	"github.com/GabrielHCataldo/martini-gateway/internal/infra"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"
)

func main() {
	logger.Info("Start application!")
	env := "dev"
	if helper.IsGreaterThan(os.Args, 1) {
		env = os.Args[1]
	}
	envUri := fmt.Sprint("config/", env, ".martini.env")
	logger.Info("Loading envs from file:", envUri)
	if err := godotenv.Load(envUri); helper.IsNotNil(err) {
		logger.Error("Error load env from file:", envUri, "err:", err)
		return
	}

	martiniFileJsonUri := fmt.Sprint("config/", env, ".martini.json")
	logger.Info("Loading martini config from file json:", martiniFileJsonUri)
	var martini dto.Martini
	martiniBytes, err := os.ReadFile(martiniFileJsonUri)
	if helper.IsNotNil(err) {
		logger.Error("Error read martini config from file json:", martiniFileJsonUri, "err:", err)
		return
	}
	logger.Info("Start filling environment variables with $word syntax!")
	martiniBytes = fillEnvValues(martiniBytes)
	err = json.Unmarshal(martiniBytes, &martini)
	if helper.IsNotNil(err) {
		logger.Error("Error parse martini config file to DTO:", err)
		return
	} else if err = helper.Validate().Struct(martini); helper.IsNotNil(err) {
		logger.Error("Error validate martini config file:", err)
		return
	}

	redisClient := infra.ConnectRedis(os.Getenv("REDIS_URL"), os.Getenv("REDIS_PASSWORD")) //todo: passar isso para martini.json como opcional
	defer infra.DisconnectRedis()
	memoryStore := persist.NewRedisStore(redisClient) //todo: isso tem que vim a configuração config json

	logger.Info("Converting martini timeout and limiter!")
	handlerTimeout, err := time.ParseDuration(martini.Timeout.Handler)
	if err != nil {
		logger.Error("Error parse duration timeout.handler field:", err)
		return
	}
	if helper.IsEmpty(martini.Limiter.MaxSizeRequestBody) {
		martini.Limiter.MaxSizeRequestBody = "3MB"
	}
	maxSizeRequestBody, err := helper.ConvertByteUnit(martini.Limiter.MaxSizeRequestBody)
	if err != nil {
		logger.Error("Error parse byte unit limiter.max-size-request-body field:", err)
		return
	}
	if helper.IsEmpty(martini.Limiter.MaxSizeMultipartMemory) {
		martini.Limiter.MaxSizeMultipartMemory = "5MB"
	}
	maxSizeMultipartMemory, err := helper.ConvertByteUnit(martini.Limiter.MaxSizeMultipartMemory)
	if err != nil {
		logger.Error("Error parse byte unit limiter.max-size-multipart-memory field:", err)
		return
	}

	modifierService := service.NewModifier()
	backendService := service.NewBackend()

	timeoutUseCase := usecase.NewTimeout(handlerTimeout)
	limiterUseCase := usecase.NewLimiter(
		maxSizeRequestBody,
		maxSizeMultipartMemory,
		martini.Limiter.MaxIpRequestPerSeconds,
	)
	corsUseCase := usecase.NewCors(martini.ExtraConfig.SecurityCors)
	endpointUseCase := usecase.NewEndpoint(backendService, modifierService)

	limiterMiddleware := middleware.NewLimiter(limiterUseCase)
	timeoutMiddleware := middleware.NewTimeout(timeoutUseCase)
	corsMiddleware := middleware.NewCors(corsUseCase)

	endpointController := controller.NewEndpoint(martini, endpointUseCase)

	app := application.NewGateway(
		martini,
		memoryStore,
		limiterMiddleware,
		timeoutMiddleware,
		corsMiddleware,
		endpointController,
	)
	go app.Run()

	fileBytes, _ := json.MarshalIndent(martini, "", "\t")
	_ = os.WriteFile("martini.json", fileBytes, 0644)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		logger.Info("Stop application!")
	}
}

func fillEnvValues(bytesJson []byte) []byte {
	stringJson := helper.SimpleConvertToString(bytesJson)
	regex := regexp.MustCompile(`\$\w+`)
	resultFind := regex.FindAllString(stringJson, -1)
	logger.Info(len(resultFind), "environment variable values were found to fill in!")
	countProcessed := 0
	for _, word := range resultFind {
		envJsonValue := strings.ReplaceAll(word, "$", "")
		envValue := os.Getenv(envJsonValue)
		if helper.IsNotEmpty(envValue) {
			stringJson = strings.ReplaceAll(stringJson, word, envValue)
			countProcessed++
		}
	}
	logger.Info(countProcessed, "environment variables successfully filled!")
	return helper.SimpleConvertToBytes(stringJson)
}
