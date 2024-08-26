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
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
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
	"github.com/xeipuuv/gojsonschema"
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

func (p provider) Init() *dto.Gopen {
	if checker.IsEmpty(os.Getenv("GOPEN_ENV")) || checker.IsEmpty(os.Getenv("GOPEN_PORT")) {
		panic("Please fill in the mandatory environment variables which are GOPEN_ENV and GOPEN_PORT!")
	}

	if !checker.IsInt(os.Getenv("GOPEN_PORT")) {
		panic("Please fill the environment variable GOPEN_PORT with numbers only!")
	}

	p.log.PrintLogo()

	p.log.PrintInfo("Loading Gopen envs...")
	err := p.loadEnvs()
	if checker.NonNil(err) {
		p.log.PrintWarn(err)
	}

	p.log.PrintInfo("Loading Gopen json...")
	gopen, err := p.loadJson()
	if checker.NonNil(err) {
		panic(err)
	}

	os.Setenv("ELASTIC_APM_ENVIRONMENT", os.Getenv("GOPEN_ENV"))
	os.Setenv("ELASTIC_APM_SERVICE_VERSION", gopen.Version)

	return gopen
}

func (p provider) Start(gopen *dto.Gopen) {
	p.log.PrintInfo("Configuring cache store...")
	store := cache.NewMemoryStore()
	if checker.NonNil(gopen.Store) {
		store = cache.NewRedisStore(gopen.Store.Redis.Address, gopen.Store.Redis.Password)
	}
	defer store.Close()

	err := p.writeRuntimeJson(gopen)
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
	nConverter := convert.New()
	nNomenclature := nomenclature.New()

	httpServer := server.New(gopen, p.log, router, httpClient, endpointLog, backendLog, httpLog, jsonPath, nConverter,
		store, nNomenclature)

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

func (p provider) Stop() {
	p.log.SkipLine()

	p.log.PrintInfo("Removing runtime json...")

	err := p.removeRuntimeJson()
	if checker.NonNil(err) {
		p.log.PrintWarn("Error to remove runtime json!")
	}

	p.log.PrintTitle("STOPPED")
}

func (p provider) restart(oldGopen *dto.Gopen, oldServer server.HTTP) {
	defer func() {
		if r := recover(); checker.NonNil(r) {
			errorDetails := errors.Details(r.(error))
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
	err = p.loadEnvs()
	if checker.NonNil(err) {
		p.log.PrintWarn(err)
	}

	p.log.PrintInfo("Reloading Gopen json...")
	gopen, err := p.loadJson()
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

	for _, path := range []string{p.buildEnvUri(), p.buildJsonUri()} {
		err = watcher.Add(path)
		if checker.NonNil(err) {
			return nil, err
		}
	}

	return watcher, nil
}

func (p provider) loadEnvs() (err error) {
	gopenEnvUri := p.buildEnvUri()

	if err = godotenv.Overload(gopenEnvUri); checker.NonNil(err) {
		err = errors.New("Error load Gopen envs from uri:", gopenEnvUri, "err:", err)
	}

	return err
}

func (p provider) loadJson() (*dto.Gopen, error) {
	gopenJsonUri := p.buildJsonUri()

	gopenJsonBytes, err := os.ReadFile(gopenJsonUri)
	if checker.NonNil(err) {
		return nil, errors.New("Error read Gopen config from file json:", gopenJsonUri, "err:", err)
	}
	gopenJsonBytes = p.fillEnvValues(gopenJsonBytes)

	if err = p.validateJsonBySchema(gopenJsonBytes); checker.NonNil(err) {
		return nil, err
	}

	var gopen dto.Gopen
	err = converter.ToDestWithErr(gopenJsonBytes, &gopen)
	if checker.NonNil(err) {
		return nil, err
	}

	return &gopen, nil
}

func (p provider) fillEnvValues(gopenJsonBytes []byte) []byte {
	// todo: aceitar campos não string receber variável de ambiente também
	//  foi pensado que talvez utilizar campos string e any para isso, convertendo para o tipo desejado apenas
	//  quando objeto de valor for montado
	gopenJsonStr := converter.ToString(gopenJsonBytes)

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

	return converter.ToBytes(gopenJsonStr)
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

func (p provider) buildEnvUri() string {
	return fmt.Sprintf("./gopen/%s/.env", os.Getenv("GOPEN_ENV"))
}

func (p provider) buildJsonUri() string {
	return fmt.Sprintf("./gopen/%s/.json", os.Getenv("GOPEN_ENV"))
}
