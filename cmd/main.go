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
	"github.com/GabrielHCataldo/gopen-gateway/internal/app"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/controller"
	configVO "github.com/GabrielHCataldo/gopen-gateway/internal/domain/config/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/main/service"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/config"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/middleware"
	"os"
	"os/signal"
	"time"
)

var gopenApp app.Gopen

var gopenJson *configVO.GopenJson

func main() {
	config.PrintInfoLogCmd("Starting...")

	// inicializamos o valor env para obter como argumento de aplicação
	var env string
	if helper.IsLessThanOrEqual(os.Args, 1) {
		panic("Please enter ENV as second argument! ex: dev, prd")
	}
	env = os.Args[1]

	// carregamos as variáveis de ambiente padrão da app
	config.LoadGopenDefaultEnvs()

	// carregamos as variáveis de ambiente indicada
	config.LoadGopenEnvs(env)

	// carregamos o json de configuração pelo ambiente indicado
	gopenJson = config.LoadGopenJson(env)

	// inicializamos a aplicação
	go startApp(env)

	// seguramos a goroutine principal esperando que aplicação seja interrompida
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		// removemos o arquivo json que foi usado
		config.RemoveGopenJsonResult()
		// imprimimos que a aplicação foi interrompida
		config.PrintInfoLogCmd("----------------------- <STOPPED> -----------------------")
	}
}

func startApp(env string) {
	// configuramos o store interface
	cacheStore := config.NewCacheStore(gopenJson.Store)
	defer config.CloseCacheStore(cacheStore)

	// configuramos o watch para ouvir mudanças do json de configuração caso hot-reload for true
	if gopenJson.HotReload {
		watcher := config.NewWatcher(env, restartApp)
		defer config.CloseWatcher(watcher)
	}

	// deu tudo certo escrevemos o json de resposta salvamos o gopenDTO resultante
	config.WriteGopenJsonResult(gopenJson)

	// chamamos o lister and server, ele ira segurar a goroutine do app
	listerAndServer(cacheStore, gopenJson)
}

func listerAndServer(cacheStore infra.CacheStore, gopenJson *configVO.GopenJson) {
	config.PrintInfoLogCmd("Building value objects...")
	gopen := configVO.NewGopen(gopenJson)

	config.PrintInfoLogCmd("Building infra...")
	restTemplate := infra.NewRestTemplate()
	traceProvider := infra.NewTraceProvider()
	logProvider := infra.NewLogProvider()

	config.PrintInfoLogCmd("Building domain...")
	backendService := service.NewBackend(restTemplate)
	endpointService := service.NewEndpoint(backendService)

	config.PrintInfoLogCmd("Building middlewares...")
	panicRecoveryMiddleware := middleware.NewPanicRecovery()
	traceMiddleware := middleware.NewTrace(traceProvider)
	logMiddleware := middleware.NewLog(logProvider)
	securityCorsMiddleware := middleware.NewSecurityCors(gopen.SecurityCors())
	limiterMiddleware := middleware.NewLimiter()
	timeoutMiddleware := middleware.NewTimeout()
	cacheMiddleware := middleware.NewCache(cacheStore)

	config.PrintInfoLogCmd("Building controllers...")
	staticController := controller.NewStatic(gopenJson)
	endpointController := controller.NewEndpoint(endpointService)

	config.PrintInfoLogCmd("Building application...")
	gopenApp = app.NewGopen(
		gopen,
		panicRecoveryMiddleware,
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
}

func restartApp(env string) {
	// damos um recovery para não parar a aplicação
	defer restartPanicRecovery(env)

	// print log
	fmt.Println()
	fmt.Println()
	config.PrintInfoLogCmd("----------------------- <RESTART> -----------------------")

	// inicializamos um contexto de timeout para ter um tempo de limite de tentativa
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	// paramos a aplicação, para começar com o novo DTO e as novas envs
	config.PrintInfoLogCmd("Shutting down current server...")
	err := gopenApp.Shutdown(ctx)
	if helper.IsNotNil(err) {
		config.PrintWarningLogCmdf("Error shutdown current server: %s!", err)
		return
	}

	// carregamos as variáveis de ambiente indicada
	config.ReloadGopenEnvs(env)

	// carregamos o json de configuração pelo ambiente indicado
	gopenJson = config.LoadGopenJson(env)

	// iniciamos a aplicação com as informações alteradas
	go startApp(env)
}

func restartPanicRecovery(env string) {
	// caso dê algum erro de panic para reiniciar a aplicação, damos recovery
	if r := recover(); helper.IsNotNil(r) {
		// damos o listerAndServer do app antigo caso tenha ocorrido um erro no app, caso contrário paramos o app
		config.PrintWarningLogCmd("Error restart server:", r)
		// damos o recovery na aplicação
		recoveryApp(env)
	}
}

func recoveryApp(env string) {
	fmt.Println()
	config.PrintInfoLogCmd("----------------------- <RECOVERY> -----------------------")

	// iniciamos a aplicação com o config antiga
	go startApp(env)
}
