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
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/joho/godotenv"
)

// LoadGopenDefaultEnvs loads the default environment variables for Gopen.
func LoadGopenDefaultEnvs() {
	// carregamos as envs padrões do Gopen
	PrintInfoLogCmd("Loading Gopen envs default...")
	if err := godotenv.Load("./.env"); helper.IsNotNil(err) {
		panic(errors.New("Error load Gopen envs default:", err))
	}
}

// LoadGopenEnvs loads the environment variables indicated by the argument 'env'.
// It retrieves the file URI based on the 'env' and loads the environment variables from that file.
// It logs the process and prints a warning if there is an error loading the environment variables.
func LoadGopenEnvs(env string) {
	// carregamos as envs indicada no arg
	fileEnvUri := getFileEnvUri(env)
	PrintInfoLogCmdf("Loading Gopen envs from uri: %s...", fileEnvUri)
	if err := godotenv.Load(fileEnvUri); helper.IsNotNil(err) {
		PrintWarningLogCmd("Error load Gopen envs from uri:", fileEnvUri, "err:", err)
	}
}

// ReloadGopenEnvs loads the environment variables indicated by the env argument.
//
// It performs the following steps:
// 1. Retrieves the file URI for the specified environment using getFileEnvUri.
// 2. Prints an information log message indicating the URI being loaded.
// 3. Uses godotenv.Overload to load the environment variables from the file URI.
// 4. If an error occurs during the loading process, prints a warning log message.
//
// The function receives a string 'env' as the required environment parameter.
// It does not return any value.
func ReloadGopenEnvs(env string) {
	// carregamos as envs indicada no arg
	fileEnvUri := getFileEnvUri(env)
	PrintInfoLogCmdf("Reloading Gopen envs from uri: %s...", fileEnvUri)
	if err := godotenv.Overload(fileEnvUri); helper.IsNotNil(err) {
		PrintWarningLogCmd("Error reload Gopen envs from uri:", fileEnvUri, "err:", err)
	}
}
