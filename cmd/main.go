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

// env variable represents the environment in which the software is running.
var env string

// loggerProvider represents the provider for command line logging functionality. It is used to print
// various log messages, titles, and warnings during the execution of the software.
var loggerProvider interfaces.CmdLoggerProvider

// jsonProvider represents the provider for JSON-related functionalities in the application.
// It provides methods to read JSON from a URI, validate JSON against a JSON schema,
// write JSON to a file, and remove a JSON file.
// The implementation of this provider should satisfy the JsonProvider interface.
// JsonProvider can be initialized using `infra.NewJsonProvider()`.
var jsonProvider interfaces.JsonProvider

// bootService represents the implementation of the Boot interface. It provides
// methods for bootstrapping the application, loading configurations, initializing
// cache store, creating a file watcher, and starting the application.
var bootService service.Boot

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

// init is a function that initializes the software by setting up the environment,
// creating providers, and instantiating the bootService.
// It should be called before any other code in the package.
//
// Constraints:
//
// - Requires at least one command-line argument (environment) provided as os.Args[1].
// - Panics if no command-line argument is provided.
//
// Steps:
//
// 1. Check if the command-line argument count is less than or equal to 1.
//   - If true, panic with the message "Please enter ENV as second argument! ex: dev, prd".
//
// 2. Set the "env" variable to the value of os.Args[1].
// 3. Create a new instance of the "cacheStoreProvider" using infra.NewCacheStoreProvider().
// 4. Create a new instance of the "jsonProvider" using infra.NewJsonProvider().
// 5. Create a new instance of the "loggerProvider" using infra.NewCmdLoggerProvider().
// 6. Create a new instance of the "bootService" using service.NewBoot() with the following parameters:
//   - env: The value of the "env" variable.
//   - cacheStoreProvider: The instance of the "cacheStoreProvider".
//   - loggerProvider: The instance of the "loggerProvider".
//   - jsonProvider: The instance of the "jsonProvider".
func init() {
	if helper.IsLessThanOrEqual(os.Args, 1) {
		panic("Please enter ENV as second argument! ex: dev, prd")
	}
	env = os.Args[1]

	cacheStoreProvider := infra.NewCacheStoreProvider()
	jsonProvider = infra.NewJsonProvider()
	loggerProvider = infra.NewCmdLoggerProvider()

	bootService = service.NewBoot(env, cacheStoreProvider, loggerProvider, jsonProvider)
}

// main is the entry point of the application. It starts the bootService, loads default envs,
// prints the logo and title, loads envs, starts the application in a separate goroutine, and keeps the main goroutine active.
//
// Steps:
// 1. Load the default envs by calling bootService.LoadDefaultEnvs().
// 2. Print the logo and title using loggerProvider.PrintLogo(os.Getenv("VERSION")) and loggerProvider.PrintTitle("START").
// 3. Load envs by calling bootService.LoadEnvs().
// 4. Start the application in a separate goroutine by calling go startApp().
// 5. Keep the main goroutine active indefinitely by calling keepActive().
func main() {
	bootService.LoadDefaultEnvs()

	loggerProvider.PrintLogo(os.Getenv("VERSION"))
	loggerProvider.PrintTitle("START")

	bootService.LoadEnvs()

	startApp()

	keepActive()
}

// startApp is a function that starts the Gopen application with the loaded configuration and cache store.
// It performs the following steps:
//  1. Load the gopenJson configuration by calling bootService.LoadJson().
//  2. Get the cache store instance by calling bootService.CacheStore(gopenJson.Store).
//  3. Close the cache store when the function returns by calling closeCacheStore(cacheStore) using a defer statement.
//  4. Start a file watcher if the hot-reload flag is enabled in the gopenJson configuration, and close the watcher when the function returns
//     by calling bootService.Watcher(restartApp) and using a defer statement to call closeWatcher(watcher).
//  5. Write the gopenJson configuration to the runtime JSON file by calling jsonProvider.WriteGopenJson(gopenJson).
//  6. If there is an error while writing the runtime JSON, log a warning message using loggerProvider.PrintWarning().
//  7. Log an informational message indicating that the application is being built using loggerProvider.PrintInfo().
//  8. Create a new instance of the Gopen application by calling app.NewGopen(gopenJson, cacheStore).
//  9. Call the ListerAndServer method on the Gopen application instance to start the application.
func startApp() {
	gopenJson = bootService.LoadJson()

	cacheStore := bootService.CacheStore(gopenJson.Store)
	defer closeCacheStore(cacheStore)

	if gopenJson.HotReload {
		watcher := bootService.Watcher(restartApp)
		defer closeWatcher(watcher)
	}

	err := jsonProvider.WriteGopenJson(gopenJson)
	if helper.IsNotNil(err) {
		loggerProvider.PrintWarning("Error to write runtime json! err:", err)
	}

	loggerProvider.PrintInfo("Building application...")
	gopenApp = app.NewGopen(gopenJson, cacheStore)
	go gopenApp.ListerAndServer()
}

// restartApp is a function that restarts the Gopen application by performing the following steps:
//
// 1. Recover from any panic that occurs during the execution of this function by calling the restartPanicRecovery()
// function.
// 2. Print two empty lines to separate the restart log messages from previous messages.
// 3. Print a title "RESTART" using the loggerProvider.PrintTitle() method.
// 4. Create a new context with a timeout of 30 seconds using the context.WithTimeout() function.
// 5. Defer the cancel() function to cancel the context when this function returns.
// 6. Print an info log message "Shutting down current server..." using the loggerProvider.PrintInfo() method.
// 7. Call the Shutdown() method on the gopenApp instance to shut down the current server.
// 8. Check if an error occurred during the server shutdown.
//   - If true, print a warning log message with the error using the loggerProvider.PrintWarningf() method and return
//     from the function.
//
// 9. Reload the environment variables by calling the bootService.ReloadEnvs() method.
// 10. Call the startApp() function in a new goroutine to start the Gopen application again.
func restartApp() {
	defer restartPanicRecovery()

	fmt.Println()
	fmt.Println()
	loggerProvider.PrintTitle("RESTART")

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	loggerProvider.PrintInfo("Shutting down current server...")
	err := gopenApp.Shutdown(ctx)
	if helper.IsNotNil(err) {
		loggerProvider.PrintWarningf("Error shutdown current server: %s!", err)
		return
	}

	bootService.ReloadEnvs()

	startApp()
}

// restartPanicRecovery is a function that handles panics and performs recovery operations.
// It checks if a panic has occurred and if so, it logs a warning message using the loggerProvider
// and calls the recoveryApp function to perform the recovery process.
func restartPanicRecovery() {
	if r := recover(); helper.IsNotNil(r) {
		loggerProvider.PrintWarning("Error restart server:", r)
		recoveryApp()
	}
}

// recoveryApp is a function that prints a title "RECOVERY" using the loggerProvider.PrintTitle() method,
// and then starts the application by calling the startApp() function as a goroutine.
func recoveryApp() {
	fmt.Println()
	loggerProvider.PrintTitle("RECOVERY")
	startApp()
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

// stopApp stops the application by removing the runtime JSON and printing a stop message.
// It performs the following steps:
//
// 1. Remove the runtime JSON using jsonProvider.RemoveGopenJson().
//   - If an error occurs during the removal, print a warning message using loggerProvider.PrintWarning().
//
// 2. Print a "STOPPED" title message using loggerProvider.PrintTitle().
//
// This function should be called to gracefully stop the application and perform any necessary cleanup or shutdown operations.
//
// Note: The behavior of this function may depend on the implementations of jsonProvider and loggerProvider.
// Please refer to their documentation for more details.
func stopApp() {
	err := jsonProvider.RemoveGopenJson()
	if helper.IsNotNil(err) {
		loggerProvider.PrintWarning("Error to remove runtime json!")
	}
	loggerProvider.PrintTitle("STOPPED")
}

// closeWatcher is a function that closes the file watcher by performing the following steps:
//
// 1. Check if the "watcher" parameter is not nil.
//   - If true, continue to the next step. Otherwise, return from the function.
//
// 2. Call the Close() method on the "watcher" instance to close the file watcher.
// 3. Check if an error occurred during the watcher close operation.
//   - If true, print a warning log message with the error using the loggerProvider.PrintWarningf() method.
func closeWatcher(watcher *fsnotify.Watcher) {
	if helper.IsNotNil(watcher) {
		err := watcher.Close()
		if helper.IsNotNil(err) {
			loggerProvider.PrintWarningf("Error close watcher: %s", err)
		}
	}
}

// closeCacheStore closes the given cacheStore. If an error occurs during the close operation, it logs a warning message.
func closeCacheStore(cacheStore interfaces.CacheStore) {
	err := cacheStore.Close()
	if helper.IsNotNil(err) {
		logger.Warning("Error close cache store:", err)
	}
}
