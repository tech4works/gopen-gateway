package service

import (
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/interfaces"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"os"
	"regexp"
	"strings"
)

// jsonSchemaUri is a constant string representing the URI of the JSON schema file.
const jsonSchemaUri = "file://./json-schema.json"

type bootService struct {
	env                string
	cacheStoreProvider interfaces.CacheStoreProvider
	loggerProvider     interfaces.CmdLoggerProvider
	jsonProvider       interfaces.JsonProvider
}

type Boot interface {
	LoadDefaultEnvs()
	LoadEnvs()
	ReloadEnvs()
	LoadJson() *vo.GopenJson
	CacheStore(storeJson *vo.StoreJson) interfaces.CacheStore
	Watcher(eventCallback func()) *fsnotify.Watcher
}

func NewBoot(env string, cacheStoreProvider interfaces.CacheStoreProvider, loggerProvider interfaces.CmdLoggerProvider,
	jsonProvider interfaces.JsonProvider) Boot {
	return bootService{
		env:                env,
		cacheStoreProvider: cacheStoreProvider,
		loggerProvider:     loggerProvider,
		jsonProvider:       jsonProvider,
	}
}

func (b bootService) LoadDefaultEnvs() {
	if err := godotenv.Load("./.env"); helper.IsNotNil(err) {
		panic(errors.New("Error load Gopen envs default:", err))
	}
}

func (b bootService) LoadEnvs() {
	// montamos a uri com base no parâmetro env
	gopenEnvUri := b.getGopenEnvUri()

	// imprimimos o log de carregamento
	b.loggerProvider.PrintInfof("Loading Gopen envs from uri: %s...", gopenEnvUri)

	// carregamos as envs pela uri
	if err := godotenv.Load(gopenEnvUri); helper.IsNotNil(err) {
		panic(errors.New("Error load Gopen envs from uri:", gopenEnvUri, "err:", err))
	}
}

func (b bootService) ReloadEnvs() {
	// montamos a uri com base no parâmetro env
	gopenEnvUri := b.getGopenEnvUri()

	// imprimimos o log de recarregamento
	b.loggerProvider.PrintInfof("Reloading Gopen envs from uri: %s...", gopenEnvUri)

	// carregamos as envs pela uri
	if err := godotenv.Overload(gopenEnvUri); helper.IsNotNil(err) {
		panic(errors.New("Error reload Gopen envs from uri:", gopenEnvUri, "err:", err))
	}
}

func (b bootService) LoadJson() *vo.GopenJson {
	// carregamos o arquivo de json de configuração do Gopen a partir da env informada
	gopenJsonUri := b.getGopenJsonUri()

	// imprimimos o log de carregamento do arquivo json de config
	b.loggerProvider.PrintInfof("Loading Gopen json from file: %s...", gopenJsonUri)

	// carregamos o arquivo json de config
	gopenJsonBytes, err := b.jsonProvider.Read(gopenJsonUri)
	if helper.IsNotNil(err) {
		panic(errors.New("Error read Gopen config from file json:", gopenJsonUri, "err:", err))
	}

	// preenchemos os valores de variável de ambiente com a sintaxe pre-definida
	gopenJsonBytes = b.fillEnvValues(gopenJsonBytes)

	// imprimimos o log de validação pelo json schema
	b.loggerProvider.PrintInfof("Validating the %s file by schema %s...", gopenJsonUri, jsonSchemaUri)

	// validamos o schema do gopen json
	if err = b.jsonProvider.ValidateJsonBySchema(jsonSchemaUri, gopenJsonBytes); helper.IsNotNil(err) {
		panic(err)
	}

	// imprimimos o log de informação de parsing para VO
	b.loggerProvider.PrintInfo("Parsing json file to object value...")

	// construímos o objeto de valor relacionado ao json de configuração
	gopenJson, err := vo.NewGopenJson(gopenJsonBytes)
	if helper.IsNotNil(err) {
		panic(err)
	}

	// se tudo ocorrer bem retornamos
	return gopenJson
}

func (b bootService) CacheStore(storeJson *vo.StoreJson) interfaces.CacheStore {
	b.loggerProvider.PrintInfo("Configuring cache store...")
	if helper.IsNotNil(storeJson) {
		return b.cacheStoreProvider.Redis(storeJson.Redis.Address, storeJson.Redis.Password)
	}
	return b.cacheStoreProvider.Memory()
}

func (b bootService) Watcher(eventCallback func()) *fsnotify.Watcher {
	b.loggerProvider.PrintInfo("Configuring watcher...")

	// instânciamos o novo watcher
	watcher, err := fsnotify.NewWatcher()
	if helper.IsNotNil(err) {
		b.loggerProvider.PrintWarning("Error configure watcher:", err)
		return nil
	}

	// inicializamos o novo goroutine de ouvinte de eventos
	go b.watchEvents(watcher, eventCallback)

	// adicionamos os arquivos a serem observados
	gopenEnvUri := b.getGopenEnvUri()
	gopenJsonUri := b.getGopenJsonUri()

	// primeiro tentamos adicionar o .env
	err = watcher.Add(gopenEnvUri)
	if helper.IsNotNil(err) {
		b.loggerProvider.PrintWarningf("Error add watcher on file: %s err: %s", gopenEnvUri, err)
	}
	// depois tentamos adicionar o .json
	err = watcher.Add(gopenJsonUri)
	if helper.IsNotNil(err) {
		b.loggerProvider.PrintWarningf("Error add watcher on file: %s err: %s", gopenJsonUri, err)
	}

	return watcher
}

// getGopenEnvUri returns the file URI for the given environment.
// The returned URI follows the format "./gopen/{env}.env".
func (b bootService) getGopenEnvUri() string {
	return fmt.Sprintf("./gopen/%s/.env", b.env)
}

// getGopenJsonUri returns the file URI for the specified environment's JSON file.
// The returned URI follows the format "./gopen/{env}.json".
func (b bootService) getGopenJsonUri() string {
	return fmt.Sprintf("./gopen/%s/.json", b.env)
}

func (b bootService) fillEnvValues(gopenJsonBytes []byte) []byte {
	// todo: aceitar campos não string receber variável de ambiente também
	//  foi pensado que talvez utilizar campos string e any para isso, convertendo para o tipo desejado apenas
	//  quando objeto de valor for montado

	b.loggerProvider.PrintInfo("Filling environment variables with $word syntax...")

	// convertemos os bytes do gopen json em string
	gopenJsonStr := helper.SimpleConvertToString(gopenJsonBytes)

	// compilamos o regex indicando um valor de env $API_KEY por exemplo
	regex := regexp.MustCompile(`\$\w+`)
	// damos o find pelo regex
	words := regex.FindAllString(gopenJsonStr, -1)

	// imprimimos todas as palavras encontradas a ser preenchidas
	b.loggerProvider.PrintInfo(len(words), "environment variable values were found to fill in!")

	// inicializamos o contador de valores processados
	count := 0
	for _, word := range words {
		// replace do valor padrão $
		envKey := strings.ReplaceAll(word, "$", "")
		// obtemos o valor da env pela chave indicada
		envValue := os.Getenv(envKey)
		// caso valor encontrado, damos o replace da palavra encontrada pelo valor
		if helper.IsNotEmpty(envValue) {
			gopenJsonStr = strings.ReplaceAll(gopenJsonStr, word, envValue)
			count++
		}
	}
	// imprimimos a quantidade de envs preenchidas
	b.loggerProvider.PrintInfo(count, "environment variables successfully filled!")

	// convertemos esse novo
	return helper.SimpleConvertToBytes(gopenJsonStr)
}

func (b bootService) watchEvents(watcher *fsnotify.Watcher, callback func()) {
	// abrimos um for infinito para sempre ouvir os eventos do watcher
	for {
		// prendemos o loop atual aguardando o canal ser notificado de watcher
		select {
		case event, ok := <-watcher.Events:
			// chamamos a função que executa o evento
			b.executeEvent(event, ok, callback)
		case err, ok := <-watcher.Errors:
			// chamamos a função que executa o evento de erro
			b.executeErrorEvent(err, ok)
		}
	}
}

func (b bootService) executeEvent(event fsnotify.Event, ok bool, callback func()) {
	if !ok {
		return
	}
	b.loggerProvider.PrintInfof("File modification event %s on file %s triggered!", event.Op.String(), event.Name)
	callback()
}

func (b bootService) executeErrorEvent(err error, ok bool) {
	if !ok {
		return
	}
	b.loggerProvider.PrintWarningf("Watcher event error triggered! err: %s", err)
}
