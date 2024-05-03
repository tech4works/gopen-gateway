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
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/controller"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/service"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/config"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/middleware"
	"os"
	"os/signal"
	"time"
)

// gopenApp is an instance of the app.Gopen interface that represents the functionality of a Gopen server.
// It is used to start and shutdown the Gopen application by invoking its ListerAndServer() and
// Shutdown(ctx context.Context) error methods.
var gopenApp app.Gopen

var gopenJsonVO *vo.GopenJson

// main is the entry point of the application.
// It prints an informational log message to indicate that the application is starting.
// It initializes the 'env' variable by retrieving it from the command-line arguments.
// If there is no 'env' argument provided, it panics with an error message.
// It loads the default environment variables for Gopen.
// It loads the environment variables indicated by the 'env' argument.
// It starts the 'startApp' function as a goroutine.
// It waits for an interrupt signal and removes the JSON result file when the signal is received.
func main() {
	config.PrintInfoLogCmd("Starting...")

	// inicializamos o valor env para obter como argumento de aplicação
	var env string
	if helper.IsLessThanOrEqual(os.Args, 1) {
		panic(errors.New("Please enter ENV as second argument! ex: dev, prd"))
	}
	env = os.Args[1]

	// carregamos as variáveis de ambiente padrão da app
	config.LoadGopenDefaultEnvs()

	// carregamos as variáveis de ambiente indicada
	config.LoadGopenEnvs(env)

	// carregamos o json de configuração pelo ambiente indicado
	gopenJsonVO = config.LoadGopenJson(env)

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
		config.PrintInfoLogCmd("Gopen stopped!")
	}
}

func startApp(env string) {
	// configuramos o store interface
	cacheStore := config.NewCacheStore(gopenJsonVO.Store)
	defer config.CloseCacheStore(cacheStore)

	// configuramos o watch para ouvir mudanças do json de configuração caso hot-reload for true
	if gopenJsonVO.HotReload {
		watcher := config.NewWatcher(env, restartApp)
		defer config.CloseWatcher(watcher)
	}

	// deu tudo certo escrevemos o json de resposta salvamos o gopenDTO resultante
	config.WriteGopenJsonResult(gopenJsonVO)

	// chamamos o lister and server, ele ira segurar a goroutine do app
	listerAndServer(cacheStore, gopenJsonVO)
}

// listerAndServer builds the value objects, infra, domain, middlewares, controllers, and application
// necessary for running the Gopen application. It then calls the ListerAndServer method of the gopenApp instance.
//
// Parameters:
// cacheStore: An implementation of the CacheStore interface for interacting with a cache store.
// gopenJsonVO: A pointer to a GopenJson struct that represents the configuration json for the Gopen application.
func listerAndServer(cacheStore infra.CacheStore, gopenJsonVO *vo.GopenJson) {
	config.PrintInfoLogCmd("Building value objects...")
	gopenVO := vo.NewGopen(gopenJsonVO)

	config.PrintInfoLogCmd("Building infra...")
	restTemplate := infra.NewRestTemplate()
	traceProvider := infra.NewTraceProvider()
	logProvider := infra.NewLogProvider()

	config.PrintInfoLogCmd("Building domain...")
	modifierService := service.NewModifier()
	backendService := service.NewBackend(modifierService, restTemplate)
	endpointService := service.NewEndpoint(backendService)

	config.PrintInfoLogCmd("Building middlewares...")
	panicRecoveryMiddleware := middleware.NewPanicRecovery()
	traceMiddleware := middleware.NewTrace(traceProvider)
	logMiddleware := middleware.NewLog(logProvider)
	securityCorsMiddleware := middleware.NewSecurityCors(gopenVO.SecurityCors())
	limiterMiddleware := middleware.NewLimiter()
	timeoutMiddleware := middleware.NewTimeout()
	cacheMiddleware := middleware.NewCache(cacheStore)

	config.PrintInfoLogCmd("Building controllers...")
	staticController := controller.NewStatic(gopenJsonVO)
	endpointController := controller.NewEndpoint(endpointService)

	config.PrintInfoLogCmd("Building application...")
	gopenApp = app.NewGopen(
		gopenVO,
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
	config.PrintInfoLogCmd("----------------------- RESTART -----------------------")

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
	gopenJsonVO = config.LoadGopenJson(env)

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
	config.PrintInfoLogCmd("----------------------- RECOVERY -----------------------")

	// iniciamos a aplicação com o config antiga
	go startApp(env)
}
