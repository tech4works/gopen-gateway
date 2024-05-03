package config

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

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/fsnotify/fsnotify"
)

func NewWatcher(env string, eventCallback func(env string)) *fsnotify.Watcher {
	PrintInfoLogCmd("Configuring watcher...")

	// instânciamos o novo watcher
	watcher, err := fsnotify.NewWatcher()
	if helper.IsNotNil(err) {
		PrintWarningLogCmd("Error configure watcher:", err)
	}

	// inicializamos o novo goroutine de ouvinte de eventos
	go watchEvents(env, watcher, eventCallback)

	// adicionamos os arquivos a serem observados
	fileEnvUri := getFileEnvUri(env)
	fileJsonUri := getFileJsonUri(env)

	// primeiro tentamos adicionar o .env
	err = watcher.Add(fileEnvUri)
	if helper.IsNotNil(err) {
		PrintWarningLogCmdf("Error add watcher on file: %s err: %s", fileEnvUri, err)
	}
	// depois tentamos adicionar o .json
	err = watcher.Add(fileJsonUri)
	if helper.IsNotNil(err) {
		PrintWarningLogCmdf("Error add watcher on file: %s err: %s", fileJsonUri, err)
	}

	return watcher
}

func CloseWatcher(watcher *fsnotify.Watcher) {
	err := watcher.Close()
	if helper.IsNotNil(err) {
		PrintWarningLogCmdf("Error close watcher: %s", err)
	}
}

// watchEvents listens for file events using the provided fsnotify.Watcher and invokes the callback function for each event.
// It uses a select statement to handle both file events and errors.
// If an event occurs, it calls the executeEvent function passing the event, environment, and callback.
// If an error occurs, it calls the executeErrorEvent function passing the error.
func watchEvents(env string, watcher *fsnotify.Watcher, callback func(env string)) {
	// abrimos um for infinito para sempre ouvir os eventos do watcher
	for {
		// prendemos o loop atual aguardando o canal ser notificado de watcher
		select {
		case event, ok := <-watcher.Events:
			// chamamos a função que executa o evento
			executeEvent(env, event, ok, callback)
		case err, ok := <-watcher.Errors:
			// chamamos a função que executa o evento de erro
			executeErrorEvent(err, ok)
		}
	}
}

// executeEvent logs the file modification event and triggers the callback.
// If ok is false, it returns immediately.
// It logs an information message with the event type and file name.
// Then, it calls the callback function with the environment.
//
// Parameters:
// - env: The environment.
// - event: The fsnotify.Event structure containing event details.
// - ok: Indicates whether the event was received successfully or an error occurred.
// - callback: The function to be called with the environment.
func executeEvent(env string, event fsnotify.Event, ok bool, callback func(env string)) {
	if !ok {
		return
	}
	PrintInfoLogCmdf("File modification event %s on file %s triggered!", event.Op.String(), event.Name)
	callback(env)
}

// executeErrorEvent logs a warning message with the given error.
// If ok is false, it returns immediately.
//
// Parameters:
// - err: The error.
// - ok: Indicates whether the error was received successfully or an error occurred.
func executeErrorEvent(err error, ok bool) {
	if !ok {
		return
	}
	PrintWarningLogCmdf("Watcher event error triggered! err: %s", err)
}
