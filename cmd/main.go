/*
 * Copyright 2024 Gabriel Cataldo
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

package main

import (
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/interfaces"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/service"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra"
	"github.com/fsnotify/fsnotify"
	"os"
	"os/signal"
	"time"
)

var env string

var loggerProvider interfaces.CmdLoggerProvider

var jsonProvider interfaces.JsonProvider

var bootService service.Boot

var gopenApp app.Gopen

var gopenJson *vo.GopenJson

func init() {
	// inicializamos o valor env para obter como argumento de aplicação
	if helper.IsLessThanOrEqual(os.Args, 1) {
		panic("Please enter ENV as second argument! ex: dev, prd")
	}
	env = os.Args[1]
	// inicializamos as depêndencias para o serviço
	cacheStoreProvider := infra.NewCacheStoreProvider()
	jsonProvider = infra.NewJsonProvider()
	loggerProvider = infra.NewCmdLoggerProvider()
	// inicializamos o serviço de domínio de inicialização
	bootService = service.NewBoot(env, cacheStoreProvider, loggerProvider, jsonProvider)
}

func main() {
	// carregamos as variáveis de ambiente padrão da app
	bootService.LoadDefaultEnvs()

	// imprimimos os logs de start
	loggerProvider.PrintLogo(os.Getenv("VERSION"))
	loggerProvider.PrintTitle("START")

	// carregamos as variáveis de ambiente indicada
	bootService.LoadEnvs()

	// carregamos o json de configuração pelo ambiente indicado
	gopenJson = bootService.LoadJson()

	// inicializamos a aplicação
	go startApp()

	// seguramos a goroutine principal esperando que aplicação seja interrompida
	keepActive()
}

func startApp() {
	// configuramos o store interface
	cacheStore := bootService.CacheStore(gopenJson.Store)
	defer closeCacheStore(cacheStore)

	// configuramos o watch para ouvir mudanças do json de configuração caso hot-reload for true
	if gopenJson.HotReload {
		watcher := bootService.Watcher(restartApp)
		defer closeWatcher(watcher)
	}

	// deu tudo certo escrevemos o json de resposta salvamos o vo resultante
	err := jsonProvider.WriteGopenJson(gopenJson)
	if helper.IsNotNil(err) {
		loggerProvider.PrintWarning("Error to write runtime json! err:", err)
	}

	// construímos o gopen app passando gopenJson e cacheStore
	loggerProvider.PrintInfo("Building application...")
	gopenApp = app.NewGopen(gopenJson, cacheStore)

	// chamamos o lister and server da aplicação
	gopenApp.ListerAndServer()
}

func restartApp() {
	// damos um recovery para não parar a aplicação caso de um erro panic ao restartar
	defer restartPanicRecovery()

	// print log de restart
	fmt.Println()
	fmt.Println()
	loggerProvider.PrintTitle("RESTART")

	// inicializamos um contexto de timeout para ter um tempo de limite de tentativa
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	// paramos a aplicação, para começar com o novo DTO e as novas envs
	loggerProvider.PrintInfo("Shutting down current server...")
	err := gopenApp.Shutdown(ctx)
	if helper.IsNotNil(err) {
		loggerProvider.PrintWarningf("Error shutdown current server: %s!", err)
		return
	}

	// carregamos as variáveis de ambiente indicada
	bootService.ReloadEnvs()

	// carregamos o json de configuração pelo ambiente indicado
	gopenJson = bootService.LoadJson()

	// iniciamos a aplicação com as informações alteradas
	go startApp()
}

func restartPanicRecovery() {
	// caso dê algum erro de panic para reiniciar a aplicação, damos recovery
	if r := recover(); helper.IsNotNil(r) {
		// damos o recovery do app antigo caso tenha ocorrido um erro no app, caso contrário paramos o app
		loggerProvider.PrintWarning("Error restart server:", r)
		// damos o recovery na aplicação
		recoveryApp()
	}
}

func recoveryApp() {
	fmt.Println()
	loggerProvider.PrintTitle("RECOVERY")
	go startApp()
}

func keepActive() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		stopApp()
	}
}

func stopApp() {
	// removemos o arquivo json que foi usado
	err := jsonProvider.RemoveGopenJson()
	if helper.IsNotNil(err) {
		loggerProvider.PrintWarning("Error to remove runtime json!")
	}
	// imprimimos que a aplicação foi interrompida
	loggerProvider.PrintTitle("STOPPED")
}

func closeWatcher(watcher *fsnotify.Watcher) {
	if helper.IsNotNil(watcher) {
		err := watcher.Close()
		if helper.IsNotNil(err) {
			loggerProvider.PrintWarningf("Error close watcher: %s", err)
		}
	}
}

func closeCacheStore(cacheStore interfaces.CacheStore) {
	err := cacheStore.Close()
	if helper.IsNotNil(err) {
		logger.Warning("Error close cache store:", err)
	}
}
