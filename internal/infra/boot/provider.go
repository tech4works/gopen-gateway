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
	"github.com/tech4works/checker"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/app/server"
	"github.com/tech4works/gopen-gateway/internal/infra/api"
	"github.com/tech4works/gopen-gateway/internal/infra/cache"
	"github.com/tech4works/gopen-gateway/internal/infra/converter"
	"github.com/tech4works/gopen-gateway/internal/infra/http"
	"github.com/tech4works/gopen-gateway/internal/infra/jsonpath"
	"github.com/tech4works/gopen-gateway/internal/infra/log"
	"github.com/tech4works/gopen-gateway/internal/infra/nomenclature"
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
const jsonSchemaUri = "file://./json-schema.json"

type provider struct {
	log app.BootLog
}

func New() app.Boot {
	return provider{
		log: log.NewBoot(),
	}
}

func (p provider) Init() string {
	if checker.IsLengthLessThanOrEqual(os.Args, 1) {
		panic("Please enter ENV as second argument! ex: dev, prd")
	}

	err := p.loadDefaultEnvs()
	if checker.NonNil(err) {
		panic(err)
	}

	tracerHost := os.Getenv("TRACER_HOST")
	if checker.IsNotEmpty(tracerHost) {
		p.log.PrintInfo("Booting Tracer...")
		tracer, closer, err := p.initTracer(tracerHost)
		if checker.NonNil(err) {
			p.log.PrintWarn(err)
		} else {
			defer closer.Close()
			opentracing.SetGlobalTracer(tracer)
		}
	}

	p.log.PrintLogo()

	return os.Args[1]
}

func (p provider) Start(env string) {
	p.log.PrintInfo("Loading Gopen envs...")
	err := p.loadEnvs(env)
	if checker.NonNil(err) {
		p.log.PrintWarn(err)
	}

	p.log.PrintInfo("Loading Gopen json...")
	gopen, err := p.loadJson(env)
	if checker.NonNil(err) {
		panic(err)
	}

	p.log.PrintInfo("Configuring cache store...")
	store := cache.NewMemoryStore()
	if checker.NonNil(gopen.Store) {
		store = cache.NewRedisStore(gopen.Store.Redis.Address, gopen.Store.Redis.Password)
	}
	defer store.Close()

	err = p.writeRuntimeJson(gopen)
	if checker.NonNil(err) {
		p.log.PrintWarn(err)
	}

	p.log.PrintInfo("Building log providers...")
	endpointLog := log.NewEndpoint()
	backendLog := log.NewBackend()
	httpLog := log.NewHTTPLog()

	p.log.PrintInfo("Building server...")
	router := api.NewRouter()
	httpClient := http.NewClient()
	jsonPath := jsonpath.New()
	nConverter := converter.New()
	nNomenclature := nomenclature.New()

	httpServer := server.New(gopen, p.log, router, httpClient, endpointLog, backendLog, httpLog, jsonPath, nConverter,
		store, nNomenclature)

	if gopen.HotReload {
		p.log.PrintInfo("Configuring watcher...")
		watcher, err := p.initWatcher(env, httpServer)
		if checker.NonNil(err) {
			p.log.PrintWarn("Error configure watcher:", err)
		} else {
			defer watcher.Close()
		}
	}

	httpServer.ListenAndServe()
}

func (p provider) Stop() {
	p.log.SkipLine()

	p.log.PrintInfo("Removing runtime json...")

	err := p.removeRuntimeJson()
	if checker.NonNil(err) {
		p.log.PrintWarn("Error to remove runtime json!")
	}

	p.log.PrintTitle("STOPPED")
}

func (p provider) restart(env string, oldServer server.HTTP) {
	defer func() {
		if r := recover(); checker.NonNil(r) {
			errorDetails := errors.Details(r.(error))
			p.log.PrintError("Error restart server:", errorDetails.GetCause())

			p.recovery(oldServer)
		}
	}()

	p.log.SkipLine()
	p.log.PrintTitle("RESTART")

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	p.log.PrintInfo("Shutting down current server...")
	err := oldServer.Shutdown(ctx)
	if checker.NonNil(err) {
		p.log.PrintWarnf("Error shutdown current server: %s!", err)
		return
	}

	go p.Start(env)
}

func (p provider) recovery(oldServer server.HTTP) {
	p.log.SkipLine()
	p.log.PrintTitle("RECOVERY")

	go oldServer.ListenAndServe()
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
	return jaegerConfig.NewTracer(jaegercfg.Logger(log.NewNoop()))
}

func (p provider) initWatcher(env string, oldServer server.HTTP) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if checker.NonNil(err) {
		return nil, err
	}

	defer func() {
		if checker.NonNil(err) {
			watcher.Close()
		}
	}()

	go func() {
		for {
			select {
			case ev, ok := <-watcher.Events:
				if !ok || checker.NotEquals(ev.Op, fsnotify.Chmod) {
					return
				}
				p.restart(env, oldServer)
			}
		}
	}()

	for _, path := range []string{p.buildEnvUri(env), p.buildJsonUri(env)} {
		err = watcher.Add(path)
		if checker.NonNil(err) {
			return nil, err
		}
	}

	return watcher, nil
}

func (p provider) loadDefaultEnvs() (err error) {
	if err = godotenv.Load("./.env"); checker.NonNil(err) {
		err = errors.New("Error load Gopen envs default:", err)
	}
	return err
}

func (p provider) loadEnvs(env string) (err error) {
	gopenEnvUri := p.buildEnvUri(env)

	if err = godotenv.Overload(gopenEnvUri); checker.NonNil(err) {
		err = errors.New("Error load Gopen envs from uri:", gopenEnvUri, "err:", err)
	}

	return err
}

func (p provider) loadJson(env string) (*dto.Gopen, error) {
	gopenJsonUri := p.buildJsonUri(env)

	gopenJsonBytes, err := os.ReadFile(gopenJsonUri)
	if checker.NonNil(err) {
		return nil, errors.New("Error read Gopen config from file json:", gopenJsonUri, "err:", err)
	}
	gopenJsonBytes = p.fillEnvValues(gopenJsonBytes)

	if err = p.validateJsonBySchema(gopenJsonBytes); checker.NonNil(err) {
		return nil, err
	}

	var gopen dto.Gopen
	err = helper.ConvertToDest(gopenJsonBytes, &gopen)
	if checker.NonNil(err) {
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
		if checker.IsNotEmpty(envValue) {
			gopenJsonStr = strings.ReplaceAll(gopenJsonStr, word, envValue)
			count++
		}
	}

	return helper.SimpleConvertToBytes(gopenJsonStr)
}

func (p provider) validateJsonBySchema(jsonBytes []byte) error {
	schemaLoader := gojsonschema.NewReferenceLoader(jsonSchemaUri)
	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if checker.NonNil(err) {
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
		if checker.NonNil(err) {
			return err
		}
	}

	gopenJsonBytes, err := json.MarshalIndent(gopen, "", "\t")
	if checker.IsNil(err) {
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
