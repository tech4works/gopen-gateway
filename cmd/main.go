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
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/GabrielHCataldo/gopen-gateway/internal/infra/boot"
	"github.com/opentracing/opentracing-go"
	"os"
	"os/signal"
	"time"
)

// env variable represents the environment in which the software is running.
var env string

// gopenApp variable represents an instance of the app.Gopen interface.
// It is used to start and shutdown the Gopen application.
var gopenApp app.Gopen

// gopenJson variable is a pointer to `vo.GopenJson`. It represents the configuration JSON for the Gopen application.
// It contains various fields representing different aspects of the application configuration,
// such as comments, version, port, hot-reload flag, store configuration, timeout, cache configuration, limiter
// configuration, security CORS configuration, middlewares configuration, and endpoints configuration.
// It is used in the `startApp` function to load the JSON configuration, initialize cache store,
// write the runtime JSON, build the application, and start the listener and server.
var gopenJson *vo.GopenJson

func init() {
	if helper.IsLessThanOrEqual(os.Args, 1) {
		panic("Please enter ENV as second argument! ex: dev, prd")
	}
	env = os.Args[1]

	err := boot.LoadDefaultEnvs()
	if helper.IsNotNil(err) {
		panic(err)
	}
}

func main() {
	boot.PrintLogo(os.Getenv("VERSION"))
	boot.PrintTitle("START")

	boot.PrintInfo("Loading Gopen envs...")
	err := boot.LoadEnvs(env)
	if helper.IsNotNil(err) {
		boot.PrintWarn(err)
	}

	boot.PrintInfo("Loading Gopen json...")
	gopenJson, err = boot.LoadJson(env)
	if helper.IsNotNil(err) {
		panic(err)
	}

	jaegerHost := os.Getenv("JAEGER_HOST")
	if helper.IsNotEmpty(jaegerHost) {
		boot.PrintInfo("Booting Jaeger...")
		jaeger, closer, err := boot.InitJaeger(jaegerHost)
		if helper.IsNotNil(err) {
			boot.PrintWarn(err)
		} else {
			defer closer.Close()
			opentracing.SetGlobalTracer(jaeger)
		}
	}

	go startApp()

	keepActive()
}

func startApp() {
	boot.PrintInfo("Configuring cache store...")

	cacheStore := boot.NewCacheStore(gopenJson.Store)
	defer cacheStore.Close()

	if gopenJson.HotReload {
		boot.PrintInfo("Configuring watcher...")
		watcher, err := boot.NewWatcher(env, restartApp)
		if helper.IsNotNil(err) {
			boot.PrintWarn("Error configure watcher:", err)
		} else {
			defer watcher.Close()
		}
	}

	err := boot.WriteRuntimeJson(gopenJson)
	if helper.IsNotNil(err) {
		boot.PrintWarn(err)
	}

	boot.PrintInfo("Building application...")
	gopenApp = app.NewGopen(gopenJson, cacheStore)
	gopenApp.ListerAndServer()
}

func restartApp() {
	defer restartPanicRecovery()

	fmt.Println()
	fmt.Println()
	boot.PrintTitle("RESTART")

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	boot.PrintInfo("Shutting down current server...")
	err := gopenApp.Shutdown(ctx)
	if helper.IsNotNil(err) {
		boot.PrintWarnf("Error shutdown current server: %s!", err)
		return
	}

	err = boot.ReloadEnvs(env)
	if helper.IsNotNil(err) {
		panic(err)
	}

	restartJson, err := boot.LoadJson(env)
	if helper.IsNotNil(err) {
		panic(err)
	}
	gopenJson = restartJson

	go startApp()
}

// restartPanicRecovery is a function that handles panics and performs recovery operations.
// It checks if a panic has occurred and if so, it logs a warning message using the loggerProvider
// and calls the recoveryApp function to perform the recovery process.
func restartPanicRecovery() {
	if r := recover(); helper.IsNotNil(r) {
		errorDetails := errors.Details(r.(error))
		boot.PrintError("Error restart server:", errorDetails.GetCause())

		recoveryApp()
	}
}

// recoveryApp is a function that prints a title "RECOVERY" using the loggerProvider.PrintTitle() method,
// and then starts the application by calling the startApp() function as a goroutine.
func recoveryApp() {
	fmt.Println()
	boot.PrintTitle("RECOVERY")

	go startApp()
}

// keepActive is a function that keeps the main goroutine active by waiting for an OS interrupt signal.
// It performs the following steps:
//
// 1. Create a buffered channel of os.Signal with a capacity of 1 using make().
// 2. Notify the channel when an interrupt signal is received using signal.Notify().
// 3. Block the execution until a signal is received using a select statement.
//   - If an interrupt signal is received, call stopApp() to stop the application.
//
// This function is typically called at the end of the main() function to prevent the program from exiting
// immediately and to keep it running until an interrupt signal is received.
//
// The stopApp() function is responsible for performing any cleanup or shutdown operations before exiting the program.
func keepActive() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		stopApp()
	}
}

func stopApp() {
	err := boot.RemoveRuntimeJson()
	if helper.IsNotNil(err) {
		boot.PrintWarn("Error to remove runtime json!")
	}
	boot.PrintTitle("STOPPED")
}
