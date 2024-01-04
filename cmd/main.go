package main

import (
	"encoding/json"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/open-gateway/internal/application"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/dto"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/factory"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/infra/config"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/infra/geolocalization"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/middleware"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/service"
	"github.com/GabrielHCataldo/open-gateway/internal/domain/usecase"
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
		if err := godotenv.Load("config/" + env + ".env"); err != nil {
			logger.Errorf("Error load .env file: %s", err)
			return
		}
	}

	var configDto dto.Config
	bytes, err := os.ReadFile("config/" + env + ".config.json")
	if err != nil {
		logger.Errorf("Error read config.json file: %s", err)
		return
	}
	strJSON := string(bytes)
	bytesReadJSON := []byte(fillEnvValues(strJSON))
	err = json.Unmarshal(bytesReadJSON, &configDto)
	if err != nil {
		logger.Errorf("Error parse config.json file to Struct: %s", err)
		return
	}
	if err = helper.Validate().Struct(configDto); err != nil {
		logger.Errorf("Error validate config.json: %s", err)
		return
	}

	redisClient := config.ConnectRedis(os.Getenv("REDIS_URL"), os.Getenv("REDIS_PASSWORD"))
	defer config.DisconnectRedis()

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

	modifierFactory := factory.NewModifier()
	backendFactory := factory.NewBackend()

	localeService := service.NewLocale(geolocalizationClient)
	backendService := service.NewBackend()

	limiterMiddleware := middleware.NewLimiter(
		maxSizeRequestBody,
		maxSizeMultipartMemory,
		configDto.Limiter.MaxIpRequestPerSeconds,
	)
	timeoutMiddleware := middleware.NewTimeout(handlerTimeout)
	corsMiddleware := middleware.NewCors(configDto.ExtraConfig.SecurityCors, localeService)

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
