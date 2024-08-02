package boot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/app/server"
	"github.com/tech4works/gopen-gateway/internal/infra/api"
	"github.com/tech4works/gopen-gateway/internal/infra/cache"
	"github.com/tech4works/gopen-gateway/internal/infra/converter"
	"github.com/tech4works/gopen-gateway/internal/infra/http"
	"github.com/tech4works/gopen-gateway/internal/infra/jsonpath"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/xeipuuv/gojsonschema"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

const runtimeFolder = "./runtime"
const jsonRuntimeUri = runtimeFolder + "/.json"
const jsonSchemaUri = "./json-schema.json"

type provider struct {
	logger app.Logger
}

func New() app.Boot {
	return provider{
		logger: newLogger(),
	}
}

func (p provider) Init() string {
	if helper.IsLessThanOrEqual(os.Args, 1) {
		panic("Please enter ENV as second argument! ex: dev, prd")
	}

	err := p.loadDefaultEnvs()
	if helper.IsNotNil(err) {
		panic(err)
	}

	tracerHost := os.Getenv("TRACER_HOST")
	if helper.IsNotEmpty(tracerHost) {
		p.logger.PrintInfo("Booting Tracer...")
		tracer, closer, err := p.initTracer(tracerHost)
		if helper.IsNotNil(err) {
			p.logger.PrintWarn(err)
		} else {
			defer closer.Close()
			opentracing.SetGlobalTracer(tracer)
		}
	}

	p.logger.PrintLogo()

	return os.Args[1]
}

func (p provider) Start(env string) {
	p.logger.PrintInfo("Loading Gopen envs...")
	err := p.loadEnvs(env)
	if helper.IsNotNil(err) {
		p.logger.PrintWarn(err)
	}

	p.logger.PrintInfo("Loading Gopen json...")
	gopen, err := p.loadJson(env)
	if helper.IsNotNil(err) {
		panic(err)
	}

	p.logger.PrintInfo("Configuring cache store...")
	store := cache.NewMemoryStore()
	if helper.IsNotNil(gopen.Store) {
		store = cache.NewRedisStore(gopen.Store.Redis.Address, gopen.Store.Redis.Password)
	}
	defer store.Close()

	err = p.writeRuntimeJson(gopen)
	if helper.IsNotNil(err) {
		p.logger.PrintWarn(err)
	}

	p.logger.PrintInfo("Building application...")
	router := api.NewRouter()
	httpClient := http.NewClient()
	jsonPath := jsonpath.New()
	nConverter := converter.New()

	httpServer := server.New(gopen, p.logger, router, httpClient, jsonPath, nConverter, store)

	if gopen.HotReload {
		p.logger.PrintInfo("Configuring watcher...")
		watcher, err := p.initWatcher(env, p.restart(env, httpServer))
		if helper.IsNotNil(err) {
			p.logger.PrintWarn("Error configure watcher:", err)
		} else {
			defer watcher.Close()
		}
	}

	p.logger.PrintInfo("Starting application...")
	httpServer.ListerAndServe()
}

func (p provider) Stop() {
	fmt.Println()
	err := p.removeRuntimeJson()
	if helper.IsNotNil(err) {
		p.logger.PrintWarn("Error to remove runtime json!")
	}
	p.logger.PrintTitle("STOPPED")
}

func (p provider) restart(env string, oldServer server.HTTP) func() {
	return func() {
		defer func() {
			if r := recover(); helper.IsNotNil(r) {
				errorDetails := errors.Details(r.(error))
				p.logger.PrintError("Error restart server:", errorDetails.GetCause())

				p.recovery(oldServer)
			}
		}()

		fmt.Println()
		fmt.Println()
		p.logger.PrintTitle("RESTART")

		ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
		defer cancel()

		p.logger.PrintInfo("Shutting down current server...")
		err := oldServer.Shutdown(ctx)
		if helper.IsNotNil(err) {
			p.logger.PrintWarnf("Error shutdown current server: %s!", err)
			return
		}

		go p.Start(env)
	}
}

func (p provider) recovery(oldServer server.HTTP) {
	fmt.Println()
	p.logger.PrintTitle("RECOVERY")

	go oldServer.ListerAndServe()
}

func (p provider) initTracer(host string) (opentracing.Tracer, io.Closer, error) {
	jaegerConfig := &jaegercfg.Configuration{
		ServiceName: "gopen-gateway",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           false,
			LocalAgentHostPort: host,
		},
	}
	return jaegerConfig.NewTracer(jaegercfg.Logger(&noop{}))
}

func (p provider) initWatcher(env string, callback func()) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if helper.IsNotNil(err) {
		return nil, err
	}
	defer func() {
		if helper.IsNotNil(err) {
			watcher.Close()
		}
	}()

	go func() {
		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok {
					return
				}
				callback()
			}
		}
	}()

	for _, path := range []string{p.buildEnvUri(env), p.buildJsonUri(env)} {
		err = watcher.Add(path)
		if helper.IsNotNil(err) {
			return nil, err
		}
	}

	return watcher, nil
}

func (p provider) loadDefaultEnvs() (err error) {
	if err = godotenv.Load("./.env"); helper.IsNotNil(err) {
		err = errors.New("Error load Gopen envs default:", err)
	}
	return err
}

func (p provider) loadEnvs(env string) (err error) {
	gopenEnvUri := p.buildEnvUri(env)

	if err = godotenv.Overload(gopenEnvUri); helper.IsNotNil(err) {
		err = errors.New("Error load Gopen envs from uri:", gopenEnvUri, "err:", err)
	}

	return err
}

func (p provider) loadJson(env string) (*dto.Gopen, error) {
	gopenJsonUri := p.buildJsonUri(env)

	gopenJsonBytes, err := os.ReadFile(gopenJsonUri)
	if helper.IsNotNil(err) {
		return nil, errors.New("Error read Gopen config from file json:", gopenJsonUri, "err:", err)
	}
	gopenJsonBytes = p.fillEnvValues(gopenJsonBytes)

	if err = p.validateJsonBySchema(jsonSchemaUri, gopenJsonBytes); helper.IsNotNil(err) {
		return nil, err
	}

	var gopen dto.Gopen
	err = helper.ConvertToDest(gopenJsonBytes, &gopen)
	if helper.IsNotNil(err) {
		return nil, err
	}

	return &gopen, nil
}

func (p provider) fillEnvValues(gopenJsonBytes []byte) []byte {
	// todo: aceitar campos não string receber variável de ambiente também
	//  foi pensado que talvez utilizar campos string e any para isso, convertendo para o tipo desejado apenas
	//  quando objeto de valor for montado
	gopenJsonStr := helper.SimpleConvertToString(gopenJsonBytes)

	regex := regexp.MustCompile(`\$\w+`)
	words := regex.FindAllString(gopenJsonStr, -1)

	count := 0
	for _, word := range words {
		envKey := strings.ReplaceAll(word, "$", "")
		envValue := os.Getenv(envKey)
		if helper.IsNotEmpty(envValue) {
			gopenJsonStr = strings.ReplaceAll(gopenJsonStr, word, envValue)
			count++
		}
	}

	return helper.SimpleConvertToBytes(gopenJsonStr)
}

func (p provider) validateJsonBySchema(jsonSchemaUri string, jsonBytes []byte) error {
	schemaLoader := gojsonschema.NewReferenceLoader(jsonSchemaUri)
	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if helper.IsNotNil(err) {
		return errors.New("Error validate schema:", err)
	} else if !result.Valid() {
		errorMsg := fmt.Sprintf("Map poorly formatted!\n")
		for _, desc := range result.Errors() {
			errorMsg += fmt.Sprintf("- %s\n", desc)
		}
		err = errors.New(errorMsg)
	}

	return err
}

func (p provider) writeRuntimeJson(gopen *dto.Gopen) error {
	if _, err := os.Stat(runtimeFolder); os.IsNotExist(err) {
		err = os.MkdirAll(runtimeFolder, 0755)
		if helper.IsNotNil(err) {
			return err
		}
	}

	gopenJsonBytes, err := json.MarshalIndent(gopen, "", "\t")
	if helper.IsNil(err) {
		err = os.WriteFile(jsonRuntimeUri, gopenJsonBytes, 0644)
	}

	return err
}

func (p provider) removeRuntimeJson() error {
	err := os.Remove(runtimeFolder)
	if errors.IsNot(err, os.ErrNotExist) {
		err = nil
	}
	return err
}

func (p provider) buildEnvUri(env string) string {
	return fmt.Sprintf("./gopen/%s/.env", env)
}

func (p provider) buildJsonUri(env string) string {
	return fmt.Sprintf("./gopen/%s/.json", env)
}
