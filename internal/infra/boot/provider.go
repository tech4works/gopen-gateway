/*
 * Copyright 2024 Tech4Works
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package boot

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	aws "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/xeipuuv/gojsonschema"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"

	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/dto"
	"github.com/tech4works/gopen-gateway/internal/app/server"
	"github.com/tech4works/gopen-gateway/internal/infra/api"
	"github.com/tech4works/gopen-gateway/internal/infra/cache"
	"github.com/tech4works/gopen-gateway/internal/infra/convert"
	"github.com/tech4works/gopen-gateway/internal/infra/http"
	"github.com/tech4works/gopen-gateway/internal/infra/jsonpath"
	"github.com/tech4works/gopen-gateway/internal/infra/log"
	"github.com/tech4works/gopen-gateway/internal/infra/nomenclature"
	"github.com/tech4works/gopen-gateway/internal/infra/publisher"
	"github.com/tech4works/gopen-gateway/internal/infra/telemetry"
)

const runtimeFolder = "./runtime"
const jsonRuntimeUri = runtimeFolder + "/.json"
const jsonSchemaUri = "file://./json-schema.json"

type provider struct {
	log               app.BootLog
	telemetryShutdown func(context.Context) error
}

func New() app.Boot {
	return &provider{
		log: log.NewBoot(),
	}
}

func (p provider) Init() *dto.Gopen {
	if checker.IsEmpty(os.Getenv("ENV")) || checker.IsEmpty(os.Getenv("PORT")) {
		panic("Please fill in the mandatory environment variables which are ENV and PORT!")
	}

	if !checker.IsInt(os.Getenv("PORT")) {
		panic("Please fill the environment variable PORT with numbers only!")
	}

	errors.SetPolicy(errors.PolicyNative) // todo: preciso adicionar a opcao de trace=DEBUG,INFO,WARN,ERROR,PANIC

	p.log.PrintLogo()

	p.log.PrintInfo("Loading Gopen envs...")
	err := p.loadEnvs()
	if checker.NonNil(err) {
		p.log.PrintWarn(err)
	}

	p.log.PrintInfo("Loading Gopen json...")
	gopen, err := p.loadJSON()
	if checker.NonNil(err) {
		panic(err)
	}

	os.Setenv("OTEL_SERVICE_NAME", "gopen-gateway")

	shutdown, err := telemetry.Setup(
		context.Background(),
		"gopen-gateway",
		gopen.Version,
		os.Getenv("ENV"),
	)
	if checker.NonNil(err) {
		p.log.PrintWarn("Error setting up OpenTelemetry:", err)
	} else {
		p.telemetryShutdown = shutdown
	}

	return gopen
}

func (p provider) Start(gopen *dto.Gopen) {
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Minute)
	defer cancel()

	p.log.PrintInfo("Configuring cache store...")
	store := cache.NewMemoryStore()
	if checker.NonNil(gopen.Store) {
		store = cache.NewRedisStore(gopen.Store.Redis.Address, gopen.Store.Redis.Password)
	}
	defer store.Close()

	p.log.PrintInfo("Configuring publishers clients...")
	var sqsClient *sqs.Client
	var snsClient *sns.Client

	awsConfig, _ := aws.LoadDefaultConfig(ctx)
	if checker.NonNil(awsConfig) {
		sqsClient = sqs.NewFromConfig(awsConfig)
		snsClient = sns.NewFromConfig(awsConfig)
	}

	err := p.writeRuntimeJSON(gopen)
	if checker.NonNil(err) {
		p.log.PrintWarn(err)
	}

	p.log.PrintInfo("Building log providers...")
	middlewareLog := log.NewMiddleware()
	endpointLog := log.NewEndpoint()
	backendLog := log.NewBackend()
	httpLog := log.NewHTTPLog()

	p.log.PrintInfo("Building server...")
	router := api.NewRouter()
	httpClient := http.NewClient()
	publisherClient := publisher.NewClient(sqsClient, snsClient)
	jsonPath := jsonpath.New()
	nConverter := convert.New()
	nNomenclature := nomenclature.New()

	httpServer := server.New(gopen, p.log, router, httpClient, publisherClient, middlewareLog, endpointLog, backendLog,
		httpLog, jsonPath, nConverter, store, nNomenclature)

	if gopen.HotReload {
		p.log.PrintInfo("Configuring watcher...")
		watcher, err := p.initWatcher(gopen, httpServer)
		if checker.NonNil(err) {
			p.log.PrintWarn("Error configure watcher:", err)
		} else {
			defer watcher.Close()
		}
	}

	httpServer.ListenAndServe()
}

func (p *provider) Stop() {
	p.log.SkipLine()

	p.log.PrintInfo("Removing runtime json...")

	err := p.removeRuntimeJSON()
	if checker.NonNil(err) {
		p.log.PrintWarn("Error to remove runtime json!")
	}

	if checker.NonNil(p.telemetryShutdown) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err = p.telemetryShutdown(ctx); checker.NonNil(err) {
			p.log.PrintWarn("Error shutting down OpenTelemetry:", err)
		}
	}

	p.log.PrintTitle("STOPPED")
}

func (p provider) restart(oldGopen *dto.Gopen, oldServer server.HTTP) {
	defer func() {
		if r := recover(); checker.NonNil(r) {
			errorDetails := errors.Wrap(r.(error))
			p.log.PrintError("Error restart server:", errorDetails.Cause())

			p.recovery(oldGopen)
		}
	}()

	p.log.SkipLine()
	p.log.PrintTitle("RESTART")

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
	defer cancel()

	p.log.PrintInfo("Shutting down current server...")
	err := oldServer.Shutdown(ctx)
	if checker.NonNil(err) {
		p.log.PrintWarnf("Error shutdown current server: %s!", err)
		return
	}

	p.log.PrintInfo("Reloading Gopen envs...")
	err = p.reloadEnvs()
	if checker.NonNil(err) {
		p.log.PrintWarn(err)
	}

	p.log.PrintInfo("Reloading Gopen json...")
	gopen, err := p.loadJSON()
	if checker.NonNil(err) {
		panic(err)
	}

	p.Start(gopen)
}

func (p provider) recovery(oldGopen *dto.Gopen) {
	p.log.SkipLine()
	p.log.PrintTitle("RECOVERY")

	go p.Start(oldGopen)
}

func (p provider) initWatcher(oldGopen *dto.Gopen, oldServer server.HTTP) (*fsnotify.Watcher, error) {
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
					continue
				}
				p.restart(oldGopen, oldServer)
			}
		}
	}()

	for _, path := range []string{p.buildEnvUri(), p.buildJSONUri()} {
		err = watcher.Add(path)
		if checker.NonNil(err) {
			return nil, err
		}
	}

	return watcher, nil
}

func (p provider) loadEnvs() (err error) {
	gopenEnvUri := p.buildEnvUri()

	if err = godotenv.Load(gopenEnvUri); checker.NonNil(err) && !os.IsNotExist(err) {
		err = errors.New("Error load Gopen envs from uri:", gopenEnvUri, "err:", err)
	}

	return err
}

func (p provider) reloadEnvs() (err error) {
	gopenEnvUri := p.buildEnvUri()

	if err = godotenv.Overload(gopenEnvUri); checker.NonNil(err) {
		err = errors.New("Error load Gopen envs from uri:", gopenEnvUri, "err:", err)
	}

	return err
}

func (p provider) loadJSON() (*dto.Gopen, error) {
	gopenJSONUri := p.buildJSONUri()

	gopenJSONBytes, err := os.ReadFile(gopenJSONUri)
	if checker.NonNil(err) {
		return nil, errors.New("Error read Gopen config from file json:", gopenJSONUri, "err:", err)
	}
	gopenJSONBytes = p.fillEnvValues(gopenJSONBytes)

	if err = p.validateJSONBySchema(gopenJSONBytes); checker.NonNil(err) {
		return nil, err
	}

	var gopen dto.Gopen
	err = converter.ToDestWithErr(gopenJSONBytes, &gopen)
	if checker.NonNil(err) {
		return nil, err
	}

	return &gopen, nil
}

func (p provider) fillEnvValues(gopenJSONBytes []byte) []byte {
	// todo: aceitar campos não string receber variável de ambiente também
	//  foi pensado que talvez utilizar campos string e any para isso, convertendo para o tipo desejado apenas
	//  quando objeto de valor for montado
	gopenJSONStr := converter.ToString(gopenJSONBytes)

	regex := regexp.MustCompile(`\$\w+`)
	words := regex.FindAllString(gopenJSONStr, -1)

	count := 0
	for _, word := range words {
		envKey := strings.ReplaceAll(word, "$", "")
		envValue := os.Getenv(envKey)
		if checker.IsNotEmpty(envValue) {
			gopenJSONStr = strings.ReplaceAll(gopenJSONStr, word, envValue)
			count++
		}
	}

	return converter.ToBytes(gopenJSONStr)
}

func (p provider) validateJSONBySchema(jsonBytes []byte) error {
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

func (p provider) writeRuntimeJSON(gopen *dto.Gopen) error {
	if _, err := os.Stat(runtimeFolder); os.IsNotExist(err) {
		err = os.MkdirAll(runtimeFolder, 0755)
		if checker.NonNil(err) {
			return err
		}
	}

	gopenJSONBytes, err := json.MarshalIndent(gopen, "", "\t")
	if checker.IsNil(err) {
		err = os.WriteFile(jsonRuntimeUri, gopenJSONBytes, 0644)
	}

	return err
}

func (p provider) removeRuntimeJSON() error {
	err := os.Remove(runtimeFolder)
	if errors.IsNot(err, os.ErrNotExist) {
		err = nil
	}
	return err
}

func (p provider) buildEnvUri() string {
	path := fmt.Sprintf("./gopen/%s/.env", os.Getenv("ENV"))
	if _, err := os.Stat(path); checker.IsNil(err) {
		return path
	}
	return "./gopen/.env"
}

func (p provider) buildJSONUri() string {
	path := fmt.Sprintf("./gopen/%s/.json", os.Getenv("ENV"))
	if _, err := os.Stat(path); checker.IsNil(err) {
		return path
	}
	return "./gopen/.json"
}
