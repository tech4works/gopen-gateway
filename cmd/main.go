package main

import (
	"encoding/json"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/martini-gateway/internal/application"
	middleware2 "github.com/GabrielHCataldo/martini-gateway/internal/application/middleware"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/model/dto"
	"github.com/GabrielHCataldo/martini-gateway/internal/application/usecase"
	"github.com/GabrielHCataldo/martini-gateway/internal/domain/factory"
	"github.com/GabrielHCataldo/martini-gateway/internal/domain/service"
	"github.com/GabrielHCataldo/martini-gateway/internal/infra"
	"github.com/GabrielHCataldo/martini-gateway/internal/infra/geolocalization"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"
)

func main() {
	env := "dev"
	if len(os.Args) > 1 {
		env = os.Args[1]
	}
	if env != "prod" {
		if err := godotenv.Load("martini/" + env + ".env"); err != nil {
			logger.Errorf("Error load .env file: %s", err)
			return
		}
	}

	var configDto dto.Config
	bytes, err := os.ReadFile("martini/" + env + ".martini.json")
	if err != nil {
		logger.Errorf("Error read martini.json file: %s", err)
		return
	}
	strJSON := string(bytes)
	bytesReadJSON := []byte(fillEnvValues(strJSON))
	err = json.Unmarshal(bytesReadJSON, &configDto)
	if err != nil {
		logger.Errorf("Error parse martini.json file to Struct: %s", err)
		return
	}
	if err = helper.Validate().Struct(configDto); err != nil {
		logger.Errorf("Error validate martini.json: %s", err)
		return
	}

	redisClient := infra.ConnectRedis(os.Getenv("REDIS_URL"), os.Getenv("REDIS_PASSWORD"))
	defer infra.DisconnectRedis()

	geolocalizationClient := geolocalization.NewClient()

	handlerTimeout, err := time.ParseDuration(configDto.Timeout.Handler)
	if err != nil {
		logger.Errorf("Error parse duration timeout.handler field: %s", err)
		return
	}

	if helper.IsEmpty(configDto.Limiter.MaxSizeRequestBody) {
		configDto.Limiter.MaxSizeRequestBody = "1MB"
	}
	maxSizeRequestBody, err := helper.ConvertByteUnit(configDto.Limiter.MaxSizeRequestBody)
	if err != nil {
		logger.Errorf("Error parse byte unit limiter.maxSizeRequestBody field: %s", err)
		return
	}
	if helper.IsEmpty(configDto.Limiter.MaxSizeMultipartMemory) {
		configDto.Limiter.MaxSizeMultipartMemory = "1MB"
	}
	maxSizeMultipartMemory, err := helper.ConvertByteUnit(configDto.Limiter.MaxSizeMultipartMemory)
	if err != nil {
		logger.Errorf("Error parse byte unit limiter.maxSizeMultipartMemory field: %s", err)
		return
	}

	modifierFactory := service.NewModifier()
	backendFactory := factory.NewBackend()

	localeService := service.NewLocale(geolocalizationClient)
	backendService := service.NewBackend()

	limiterMiddleware := middleware2.NewLimiter(
		maxSizeRequestBody,
		maxSizeMultipartMemory,
		configDto.Limiter.MaxIpRequestPerSeconds,
	)
	timeoutMiddleware := middleware2.NewTimeout(handlerTimeout)
	corsMiddleware := middleware2.NewCors(configDto.ExtraConfig.SecurityCors, localeService)

	endpointUseCase := usecase.NewEndpoint(configDto, backendService, backendFactory, modifierFactory)

	app := application.NewGateway(
		configDto,
		redisClient,
		limiterMiddleware,
		timeoutMiddleware,
		corsMiddleware,
		endpointUseCase,
	)
	go app.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		logger.Info("Stop application!")
	}
}

func fillEnvValues(stringJSON string) string {
	regex := regexp.MustCompile(`\B\$\w+`)
	resultFind := regex.FindAllString(stringJSON, -1)
	for _, word := range resultFind {
		envJsonValue := strings.ReplaceAll(word, "$", "")
		envValue := os.Getenv(envJsonValue)
		if envValue != "" {
			stringJSON = strings.ReplaceAll(stringJSON, word, envValue)
		}
	}
	return stringJSON
}
