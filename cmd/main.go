package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/controller"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/service"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/middleware"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/xeipuuv/gojsonschema"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"
)

const gopenJsonResultName = "./gopen.json"
const gopenJsonSchema = "file://./json-schema.json"

var loggerOptions = logger.Options{
	CustomAfterPrefixText: "CMD",
}

var gopenApp app.GOpen

var countListerAndServer = 1

func main() {
	printInfoLog("Starting..")

	// inicializamos o valor env para obter como argumento de aplicação
	var env string
	if helper.IsLessThanOrEqual(os.Args, 1) {
		panic(errors.New("Please enter ENV as second argument! ex: dev, prd"))
	}
	env = os.Args[1]

	// carregamos as variáveis de ambiente padrão da app
	loadGOpenDefaultEnvs()

	// carregamos as variáveis de ambiente indicada
	loadGOpenEnvs(env)

	// construímos o dto de configuração do GOpen
	gopenDTO := loadGOpenJson(env)

	// salvamos o gopenDTO resultante
	writeGOpenJsonResult(gopenDTO)

	// inicializamos a aplicação
	go startApp(env, gopenDTO)

	// seguramos a goroutine principal esperando que aplicação seja interrompida
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		// removemos o arquivo json que foi usado
		removeGOpenJsonResult()
		// imprimimos que a aplicação foi interrompida
		logger.Info("GOpen Stopped!")
	}
}

func loadGOpenDefaultEnvs() {
	// carregamos as envs padrões do GOpen
	printInfoLog("Loading GOpen envs default...")
	if err := godotenv.Load("./internal/infra/config/.env"); helper.IsNotNil(err) {
		panic(errors.New("Error load GOpen envs default:", err))
	}
}

func loadGOpenEnvs(env string) {
	// carregamos as envs indicada no arg
	fileEnvUri := getFileEnvUri(env)
	printInfoLogf("Loading GOpen envs from uri: %s...", fileEnvUri)
	if err := godotenv.Load(fileEnvUri); helper.IsNotNil(err) {
		printWarningLog("Error load GOpen envs from uri:", fileEnvUri, "err:", err)
	}
}

func loadGOpenJson(env string) dto.GOpen {
	// carregamos o arquivo de json de configuração do GOpen
	fileJsonUri := getFileJsonUri(env)
	printInfoLogf("Loading GOpen json from file: %s...", fileJsonUri)
	fileJsonBytes, err := os.ReadFile(fileJsonUri)
	if helper.IsNotNil(err) {
		panic(errors.New("Error read martini config from file json:", fileJsonUri, "err:", err))
	}

	// preenchemos os valores de variável de ambiente com a sintaxe pre-definida
	fileJsonBytes = fillEnvValues(fileJsonBytes)

	// validamos o schema
	if err = validateJsonSchema(fileJsonUri, fileJsonBytes); helper.IsNotNil(err) {
		panic(err)
	}

	// convertemos o valor em bytes em DTO
	var gopenDTO dto.GOpen
	err = helper.ConvertToDest(fileJsonBytes, &gopenDTO)
	if helper.IsNotNil(err) {
		panic(errors.New("Error parse GOpen json file to DTO:", err))
	}

	// temos um double-check de validação da estrutura
	if err = helper.Validate().Struct(gopenDTO); helper.IsNotNil(err) {
		panic(errors.New("Error validate GOpenDTO:", err))
	}

	// retornamos o DTO que é a configuração do GOpen
	return gopenDTO
}

func fillEnvValues(gopenBytesJson []byte) []byte {
	printInfoLog("Filling environment variables with $word syntax..")

	// convertemos os bytes do gopen json em string
	gopenStrJson := helper.SimpleConvertToString(gopenBytesJson)

	// compilamos o regex indicando um valor de env $API_KEY por exemplo
	regex := regexp.MustCompile(`\$\w+`)
	// damos o find pelo regex
	words := regex.FindAllString(gopenStrJson, -1)

	// imprimimos todas as palavras encontradas a ser preenchidas
	printInfoLog(len(words), "environment variable values were found to fill in!")

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
	printInfoLog(count, "environment variables successfully filled!")

	// convertemos esse novo
	return helper.SimpleConvertToBytes(gopenStrJson)
}

func validateJsonSchema(fileJsonUri string, fileJsonBytes []byte) error {
	printInfoLogf("Validating the %s file schema...", fileJsonUri)

	// carregamos o schema e o documento
	schemaLoader := gojsonschema.NewReferenceLoader(gopenJsonSchema)
	documentLoader := gojsonschema.NewBytesLoader(fileJsonBytes)

	// chamamos o validate
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if helper.IsNotNil(err) {
		panic(errors.New("Error validate schema:", err))
	}

	// checamos se valido, caso nao seja formatamos a mensagem
	if !result.Valid() {
		errorMsg := fmt.Sprintf("Json %s poorly formatted!\n", fileJsonUri)
		for _, desc := range result.Errors() {
			errorMsg += fmt.Sprintf("- %s\n", desc)
		}
		return errors.New(errorMsg)
	}
	// se tudo ocorrem bem retornamos nil
	return nil
}

func buildCacheStore(storeDTO *dto.Store) infra.CacheStore {
	printInfoLog("Configuring cache store...")

	if helper.IsNotNil(storeDTO) {
		return infra.NewRedisStore(storeDTO.Redis.Address, storeDTO.Redis.Password)
	}

	return infra.NewMemoryStore()
}

func listerAndServer(cacheStore infra.CacheStore, gopenVO vo.GOpen) {
	printInfoLog("Building infra..")
	restTemplate := infra.NewRestTemplate()
	traceProvider := infra.NewTraceProvider()
	logProvider := infra.NewLogProvider()

	printInfoLog("Building domain..")
	modifierService := service.NewModifier()
	backendService := service.NewBackend(modifierService, restTemplate)
	endpointService := service.NewEndpoint(backendService)

	printInfoLog("Building middlewares..")
	traceMiddleware := middleware.NewTrace(traceProvider)
	logMiddleware := middleware.NewLog(logProvider)
	securityCorsMiddleware := middleware.NewSecurityCors(gopenVO.SecurityCors())
	limiterMiddleware := middleware.NewLimiter()
	timeoutMiddleware := middleware.NewTimeout()
	cacheMiddleware := middleware.NewCache(cacheStore)

	printInfoLog("Building controllers..")
	staticController := controller.NewStatic(gopenVO)
	endpointController := controller.NewEndpoint(endpointService)

	printInfoLog("Building application..")
	// inicializamos a aplicação
	gopenApp = app.NewGOpen(
		gopenVO,
		traceMiddleware,
		logMiddleware,
		securityCorsMiddleware,
		timeoutMiddleware,
		limiterMiddleware,
		cacheMiddleware,
		staticController,
		endpointController,
	)

	// chamamos o lister and server da aplicação
	gopenApp.ListerAndServer()

	printInfoLogf("Lister and server on %c goroutine stopped!", countListerAndServer)

	countListerAndServer++
}

func writeGOpenJsonResult(gopenDTO dto.GOpen) {
	gopenBytes, err := json.MarshalIndent(gopenDTO, "", "\t")
	if helper.IsNil(err) {
		err = os.WriteFile(gopenJsonResultName, gopenBytes, 0644)
	}
	if helper.IsNotNil(err) {
		printWarningLog("Error write file gopen.json result:", err)
	}
}

func configureWatcher(env string, gopenDTO dto.GOpen) *fsnotify.Watcher {
	if !gopenDTO.HotReload {
		return nil
	}

	printInfoLog("Configuring watcher...")

	// instânciamos o novo watcher
	watcher, err := fsnotify.NewWatcher()
	if helper.IsNotNil(err) {
		printWarningLog("Error configure watcher:", err)
	}

	// inicializamos o novo goroutine de ouvinte de eventos
	go watchEvents(env, watcher)

	// adicionamos os arquivos a serem observados
	fileEnvUri := getFileEnvUri(env)
	fileJsonUri := getFileJsonUri(env)

	// primeiro tentamos adicionar o .env
	err = watcher.Add(fileEnvUri)
	if helper.IsNotNil(err) {
		printWarningLogf("Error add watcher on file: %s err: %s", fileEnvUri, err)
	}
	// depois tentamos adicionar o .json
	err = watcher.Add(fileJsonUri)
	if helper.IsNotNil(err) {
		printWarningLogf("Error add watcher on file: %s err: %s", fileJsonUri, err)
	}

	return watcher
}

func watchEvents(env string, watcher *fsnotify.Watcher) {
	// abrimos um for infinito para sempre ouvir os eventos do watcher
	for {
		// prendemos o loop atual aguardando o canal ser notificado de watcher
		select {
		case event, ok := <-watcher.Events:
			// chamamos a função que executa o evento
			executeEvent(env, event, ok)
			break
		case err, ok := <-watcher.Errors:
			// chamamos a função que executa o evento de erro
			executeErrorEvent(err, ok)
		}
		// aguardamos até a próxima
		//time.Sleep(1000)
	}
}

func executeEvent(env string, event fsnotify.Event, ok bool) {
	if !ok {
		return
	}
	printInfoLogf("File modification event %s on file %s triggered!", event.Op.String(), event.Name)
	restartApp(env)
}

func executeErrorEvent(err error, ok bool) {
	if !ok {
		return
	}
	logger.Debug("error:", err)
}

func startApp(env string, gopenDTO dto.GOpen) {
	// configuramos o store interface
	cacheStore := buildCacheStore(gopenDTO.Store)
	defer closeCacheStore(cacheStore)

	// configuramos o watch para ouvir mudanças do json de configuração
	watcher := configureWatcher(env, gopenDTO)
	defer closeWatcher(watcher)

	// construímos os objetos de valores a partir do dto gopen
	printInfoLog("Building value objects..")
	gopenVO := vo.NewGOpen(env, gopenDTO)

	// chamamos o lister and server, ele ira segurar a goroutine, depois que ele é parado, as linhas seguintes vão ser chamados
	listerAndServer(cacheStore, gopenVO)
}

func restartApp(env string) {
	// print log
	printInfoLog("---------- RESTART ----------")

	// inicializamos um contexto de timeout para ter um tempo de limite de tentativa
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	// paramos a aplicação, para começar com o novo DTO e as novas envs
	printInfoLog("Shutting down current server...")
	err := gopenApp.Shutdown(ctx)
	if helper.IsNotNil(err) {
		printWarningLogf("Error shutdown app: %s!", err)
		return
	}

	// carregamos as variáveis de ambiente indicada
	loadGOpenEnvs(env)

	// lemos o novo DTO
	gopenDTO := loadGOpenJson(env)

	// começamos um novo app listener com as informações alteradas
	go startApp(env, gopenDTO)
}

func removeGOpenJsonResult() {
	err := os.Remove(gopenJsonResultName)
	if helper.IsNotNil(err) {
		printWarningLogf("Error remove %s err: %s", gopenJsonResultName, err)
		return
	}
}

func closeWatcher(watcher *fsnotify.Watcher) {
	if helper.IsNotNil(watcher) {
		err := watcher.Close()
		if helper.IsNotNil(err) {
			printWarningLogf("Error close watcher: %s", err)
		}
	}
}

func closeCacheStore(store infra.CacheStore) {
	err := store.Close()
	if helper.IsNotNil(err) {
		printWarningLog("Error close cache store:", err)
	}
}

func getFileEnvUri(env string) string {
	return fmt.Sprintf("gopen/%s.env", env)
}

func getFileJsonUri(env string) string {
	return fmt.Sprintf("gopen/%s.json", env)
}

func printInfoLog(msg ...any) {
	logger.InfoOpts(loggerOptions, msg...)
}

func printInfoLogf(format string, msg ...any) {
	logger.InfoOptsf(format, loggerOptions, msg...)
}

func printWarningLog(msg ...any) {
	logger.WarningOpts(loggerOptions, msg...)
}

func printWarningLogf(format string, msg ...any) {
	logger.WarningOptsf(format, loggerOptions, msg...)
}
