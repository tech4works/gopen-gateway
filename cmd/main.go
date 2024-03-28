package main

import (
	"encoding/json"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/controller"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/middleware"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/service"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"regexp"
	"strings"
)

func main() {
	logger.Info("Starting GOpen..")

	// inicializamos o valor env para obter como argumento de aplicação
	var env string
	if helper.IsLessThanOrEqual(os.Args, 1) {
		logger.Error("Please enter ENV as second argument! ex: dev, prd")
		return
	}
	env = os.Args[1]

	// carregamos as envs padrões do GOpen
	logger.Info("Loading GOpen envs default...")
	if err := godotenv.Load("internal/infra/config/.env"); helper.IsNotNil(err) {
		logger.Error("Error load GOpen envs default:", err)
		return
	}

	// carregamos as envs indicada no arg
	fileEnvUri := fmt.Sprintf("gopen/%s.env", env)
	logger.Infof("Loading GOpen envs from uri: %s...", fileEnvUri)
	if err := godotenv.Load(fileEnvUri); helper.IsNotNil(err) {
		logger.Error("Error load GOpen envs from uri:", fileEnvUri, "err:", err)
		return
	}

	// carregamos o arquivo de json de configuração do GOpen
	fileJsonUri := fmt.Sprintf("gopen/%s.json", env)
	logger.Infof("Loading GOpen json from file: %s...", fileJsonUri)
	fileJsonBytes, err := os.ReadFile(fileJsonUri)
	if helper.IsNotNil(err) {
		logger.Error("Error read martini config from file json:", fileJsonUri, "err:", err)
		return
	}

	// preenchemos os valores de variável de ambiente com a sintaxe pre-definida
	logger.Info("Filling environment variables with $word syntax..")
	fileJsonBytes = fillEnvValues(fileJsonBytes)

	// convertemos o valor em bytes em DTO
	var gopenDTO dto.GOpen
	err = helper.ConvertToDest(fileJsonBytes, &gopenDTO)
	if helper.IsNotNil(err) {
		logger.Errorf("Error parse GOpen json file to DTO: %s!", err)
		return
	} else if err = helper.Validate().Struct(gopenDTO); helper.IsNotNil(err) {
		logger.Errorf("Error validate GOpen json file: %s!", err)
		return
	}

	// configuramos o cache store
	logger.Info("Configuring cache store...")

	//todo: passar isso para martini.json como opcional
	//todo: isso tem que vim a configuração config json
	redisClient := infra.ConnectRedis(os.Getenv("REDIS_URL"), os.Getenv("REDIS_PASSWORD"))
	defer infra.DisconnectRedis()

	cacheStore := persist.NewRedisStore(redisClient)

	// construímos os objetos de valores a partir do dto gopen
	logger.Info("Instantiating value objects..")
	gopenVO := vo.NewGOpen(env, gopenDTO, cacheStore)

	logger.Info("Instantiating domain services..")
	restTemplate := infra.NewRestTemplate()
	traceProvider := infra.NewTraceProvider()
	logProvider := infra.NewLogProvider()

	modifierService := service.NewModifier()
	backendService := service.NewBackend(modifierService, restTemplate)
	endpointService := service.NewEndpoint(backendService)

	traceMiddleware := middleware.NewTrace(traceProvider)
	logMiddleware := middleware.NewLog(logProvider)
	securityCorsMiddleware := middleware.NewSecurityCors(gopenVO.SecurityCors())
	limiterMiddleware := middleware.NewLimiter()
	timeoutMiddleware := middleware.NewTimeout()
	cacheMiddleware := middleware.NewCache()

	staticController := controller.NewStatic(gopenDTO)
	endpointController := controller.NewEndpoint(gopenVO, endpointService)

	// inicializamos a aplicação
	gopenApplication := app.NewGOpen(gopenDTO, gopenVO, traceMiddleware, logMiddleware, securityCorsMiddleware,
		timeoutMiddleware, limiterMiddleware, cacheMiddleware, staticController, endpointController)
	// chamamos o lister and server para continuar
	go gopenApplication.ListerAndServer()

	// salvamos o gopenDTO resultante
	gopenBytes, err := json.MarshalIndent(gopenDTO, "", "\t")
	if helper.IsNil(err) {
		err = os.WriteFile("gopen.json", gopenBytes, 0644)
	}
	if helper.IsNotNil(err) {
		logger.Warning("Error write file gopen.json result:", err)
	}

	// seguramos a goroutine principal esperando que aplicação seja interrompida
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		logger.ResetOptionsToDefault()
		logger.Info("Stop GOpen!")
	}
}

func fillEnvValues(gopenBytesJson []byte) []byte {
	// convertemos os bytes do gopen json em string
	gopenStrJson := helper.SimpleConvertToString(gopenBytesJson)

	// compilamos o regex indicando um valor de env $API_KEY por exemplo
	regex := regexp.MustCompile(`\$\w+`)
	// damos o find pelo regex
	words := regex.FindAllString(gopenStrJson, -1)

	// imprimimos todas as palavras encontradas a ser preenchidas
	logger.Info(len(words), "environment variable values were found to fill in!")

	// inicializamos o contador de valores processados
	count := 0
	for _, word := range words {
		// replace do valor padrão $
		envKey := strings.ReplaceAll(word, "$", "")
		// obtemos o valor da env pela chave indicada
		envValue := os.Getenv(envKey)
		// caso valor encontrado, damos o replace da palavra encontrada pelo valor
		if helper.IsNotEmpty(envValue) {
			gopenStrJson = strings.ReplaceAll(gopenStrJson, word, envValue)
			count++
		}
	}
	// imprimimos a quantidade de envs preenchidas
	logger.Info(count, "environment variables successfully filled!")

	// convertemos esse novo
	return helper.SimpleConvertToBytes(gopenStrJson)
}
