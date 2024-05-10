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
const jsonSchemaUri = "https://raw.githubusercontent.com/GabrielHCataldo/gopen-gateway/main/json-schema.json"

// bootService is a type that represents the boot service for initializing the application.
// It holds the environment, cache store provider, logger provider, and JSON provider.
type bootService struct {
	// env represents the environment in which the application is running.
	// It is used for initializing and loading configuration files specific to the environment.
	env string
	// cacheStoreProvider is an interface that defines methods for creating cache store instances.
	// It provides methods for creating an in-memory cache or a Redis cache.
	cacheStoreProvider interfaces.CacheStoreProvider
	// loggerProvider is an interface that provides methods for printing different types of log messages
	// such as logo, titles, info, and warnings. It is used by the bootService for logging purposes.
	loggerProvider interfaces.CmdLoggerProvider
	// jsonProvider represents an interface for handling JSON data.
	jsonProvider interfaces.JsonProvider
}

// Boot is an interface that represents the boot service for initializing the application.
// It defines the methods for loading default environments, loading environments, reloading environments,
// loading JSON configuration, creating a cache store, and creating a Watcher for file system events.
type Boot interface {
	// LoadDefaultEnvs sets the default environment variables required for bootstrapping the application.
	LoadDefaultEnvs()
	// LoadEnvs loads the environment variables customized by the environment indicated at application startup.
	LoadEnvs()
	// ReloadEnvs reloads the environment variables customized by the environment indicated at application startup.
	ReloadEnvs()
	// LoadJson loads the GopenJson configuration for the Gopen application.
	// It returns a pointer to a GopenJson object.
	LoadJson() *vo.GopenJson
	// CacheStore is a method that takes a *vo.StoreJson as input and returns an interfaces.CacheStore.
	// It represents the cache store for the Gopen application.
	// The CacheStore interface provides methods for interacting with a cache store,
	// such as setting a value, deleting a value, getting a value, and closing the store.
	// The storeJson parameter represents the store configuration for the Gopen application,
	// which includes the Redis configuration.
	// The method returns an interfaces.CacheStore instance that can be used to interact
	// with the cache store.
	CacheStore(storeJson *vo.StoreJson) interfaces.CacheStore
	// Watcher is a method that creates a file system watcher and returns a pointer to it.
	// The watcher is used to monitor file system events.
	// The method takes an eventCallback function as input, which will be called
	// when any file system event occurs.
	// The function has no return value.
	Watcher(callback func()) *fsnotify.Watcher
}

// NewBoot returns a new instance of the bootService structure that implements the Boot interface.
// It takes the environment string, cacheStoreProvider, loggerProvider, and jsonProvider as parameters.
// The returned bootService instance will have the provided environment, cacheStoreProvider, loggerProvider, and
// jsonProvider.
func NewBoot(env string, cacheStoreProvider interfaces.CacheStoreProvider, loggerProvider interfaces.CmdLoggerProvider,
	jsonProvider interfaces.JsonProvider) Boot {
	return bootService{
		env:                env,
		cacheStoreProvider: cacheStoreProvider,
		loggerProvider:     loggerProvider,
		jsonProvider:       jsonProvider,
	}
}

// LoadDefaultEnvs loads the default environment variables from the "./.env" file.
// It uses the godotenv package to load the environment variables and panics if an error occurs.
func (b bootService) LoadDefaultEnvs() {
	if err := godotenv.Load("./.env"); helper.IsNotNil(err) {
		panic(errors.New("Error load Gopen envs default:", err))
	}
}

// LoadEnvs loads the environment variables from the file specified by the 'gopenEnvUri' variable.
// It uses the godotenv package to load the environment variables and panics if an error occurs.
func (b bootService) LoadEnvs() {
	gopenEnvUri := b.getGopenEnvUri()

	b.loggerProvider.PrintInfof("Loading Gopen envs from uri: %s...", gopenEnvUri)

	if err := godotenv.Load(gopenEnvUri); helper.IsNotNil(err) {
		panic(errors.New("Error load Gopen envs from uri:", gopenEnvUri, "err:", err))
	}
}

// ReloadEnvs reloads the environment variables from the URI specified by the 'gopenEnvUri' variable.
// It uses the godotenv package to reload the environment variables and panics if an error occurs.
func (b bootService) ReloadEnvs() {
	gopenEnvUri := b.getGopenEnvUri()

	b.loggerProvider.PrintInfof("Reloading Gopen envs from uri: %s...", gopenEnvUri)

	if err := godotenv.Overload(gopenEnvUri); helper.IsNotNil(err) {
		panic(errors.New("Error reload Gopen envs from uri:", gopenEnvUri, "err:", err))
	}
}

// LoadJson loads the Gopen configuration JSON file, validates it against a schema, and parses
// it into a GopenJson object. It first gets the URI of the JSON file using the getGopenJsonUri()
// method. It then prints a log message indicating the file loading process. Next, it reads the JSON
// file and panics if there is an error. After that, it fills the environment variable values using
// the fillEnvValues() method. It then validates the JSON file against a schema and panics if there
// is an error. Finally, it parses the JSON file into a GopenJson object and returns it.
func (b bootService) LoadJson() *vo.GopenJson {
	gopenJsonUri := b.getGopenJsonUri()

	b.loggerProvider.PrintInfof("Loading Gopen json from file: %s...", gopenJsonUri)

	gopenJsonBytes, err := b.jsonProvider.Read(gopenJsonUri)
	if helper.IsNotNil(err) {
		panic(errors.New("Error read Gopen config from file json:", gopenJsonUri, "err:", err))
	}
	gopenJsonBytes = b.fillEnvValues(gopenJsonBytes)

	b.loggerProvider.PrintInfof("Validating the %s file by schema %s...", gopenJsonUri, jsonSchemaUri)

	if err = b.jsonProvider.ValidateJsonBySchema(jsonSchemaUri, gopenJsonBytes); helper.IsNotNil(err) {
		panic(err)
	}

	b.loggerProvider.PrintInfo("Parsing json file to object value...")

	gopenJson, err := vo.NewGopenJson(gopenJsonBytes)
	if helper.IsNotNil(err) {
		panic(err)
	}
	return gopenJson
}

// CacheStore returns an implementation of the CacheStore interface based on the provided StoreJson.
// In the implementation, it configures the cache store by printing an info log message.
// If the provided StoreJson is not nil, it creates a Redis cache store using the address and password from the StoreJson.
// If the storeJson is nil, it creates an in-memory cache store.
// The method then returns the created cache store.
func (b bootService) CacheStore(storeJson *vo.StoreJson) interfaces.CacheStore {
	b.loggerProvider.PrintInfo("Configuring cache store...")
	if helper.IsNotNil(storeJson) {
		return b.cacheStoreProvider.Redis(storeJson.Redis.Address, storeJson.Redis.Password)
	}
	return b.cacheStoreProvider.Memory()
}

// Watcher configures a fsnotify.Watcher to listen for file modification events.
// It creates a new watcher instance and starts a goroutine to listen for events.
// The watcher is set to watch two files: the gopenEnvUri and the gopenJsonUri.
// If there is an error configuring the watcher, a warning message is printed and nil is returned.
// The callback function is called whenever a file modification event occurs.
// The watcher instance is returned.
func (b bootService) Watcher(callback func()) *fsnotify.Watcher {
	b.loggerProvider.PrintInfo("Configuring watcher...")

	watcher, err := fsnotify.NewWatcher()
	if helper.IsNotNil(err) {
		b.loggerProvider.PrintWarning("Error configure watcher:", err)
		return nil
	}

	go b.watchEvents(watcher, callback)

	gopenEnvUri := b.getGopenEnvUri()
	gopenJsonUri := b.getGopenJsonUri()

	err = watcher.Add(gopenEnvUri)
	if helper.IsNotNil(err) {
		b.loggerProvider.PrintWarningf("Error add watcher on file: %s err: %s", gopenEnvUri, err)
	}
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

// fillEnvValues replaces environment variable placeholders in the gopen JSON string with their corresponding values.
// It searches for values in the format of $ENV_VARIABLE_NAME and gets the value from the environment variables.
// The function uses a regular expression to find all the environment variable placeholders in the gopen JSON string.
// It then replaces each placeholder with its corresponding environment variable value.
// If a value is found for a placeholder, it replaces the placeholder in the gopen JSON string and increments a counter.
// Finally, the modified gopen JSON string is converted back to bytes and returned.
// Note: Non-string fields and environment variables can also be accepted by passing them as any type.
// This feature is not fully implemented yet.
// This function also prints debug information such as the number of environment variable values found and replaced.
// The function utilizes the loggerProvider to print the debug information.
func (b bootService) fillEnvValues(gopenJsonBytes []byte) []byte {
	// todo: aceitar campos não string receber variável de ambiente também
	//  foi pensado que talvez utilizar campos string e any para isso, convertendo para o tipo desejado apenas
	//  quando objeto de valor for montado

	b.loggerProvider.PrintInfo("Filling environment variables with $word syntax...")

	gopenJsonStr := helper.SimpleConvertToString(gopenJsonBytes)

	regex := regexp.MustCompile(`\$\w+`)
	words := regex.FindAllString(gopenJsonStr, -1)

	b.loggerProvider.PrintInfo(len(words), "environment variable values were found to fill in!")

	count := 0
	for _, word := range words {
		envKey := strings.ReplaceAll(word, "$", "")
		envValue := os.Getenv(envKey)
		if helper.IsNotEmpty(envValue) {
			gopenJsonStr = strings.ReplaceAll(gopenJsonStr, word, envValue)
			count++
		}
	}

	b.loggerProvider.PrintInfo(count, "environment variables successfully filled!")

	return helper.SimpleConvertToBytes(gopenJsonStr)
}

// watchEvents continuously listens for file system events and executes the callback function
// when an event occurs. It takes a *fsnotify.Watcher as input parameter and a callback function.
// The function calls the executeEvent function when a file system event occurs and calls the
// executeErrorEvent function when an error event occurs.
func (b bootService) watchEvents(watcher *fsnotify.Watcher, callback func()) {
	for {
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

// executeEvent executes a callback function when a file modification event is triggered.
// It takes in three parameters: event of type fsnotify.Event, ok of type bool, and callback of type func().
// If ok is false, the function returns immediately.
// The function prints an informational message using the loggerProvider, indicating the type of event and the file name.
// Finally, it calls the provided callback function.
func (b bootService) executeEvent(event fsnotify.Event, ok bool, callback func()) {
	if !ok {
		return
	}
	b.loggerProvider.PrintInfof("File modification event %s on file %s triggered!", event.Op.String(), event.Name)
	callback()
}

// executeErrorEvent handles the error event by logging a warning message.
// It prints the error message using the loggerProvider's PrintWarningf method.
// It only executes if the ok parameter is true, otherwise it does nothing.
func (b bootService) executeErrorEvent(err error, ok bool) {
	if !ok {
		return
	}
	b.loggerProvider.PrintWarningf("Watcher event error triggered! err: %s", err)
}
