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

package boot

import (
	"encoding/json"
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/joho/godotenv"
	"github.com/xeipuuv/gojsonschema"
	"os"
	"regexp"
	"strings"
)

// runtimeFolder is a constant string representing the filepath of the runtime folder for the Gopen application.
const runtimeFolder = "./runtime"

// jsonRuntimeUri is a constant string representing the filepath of the Gopen JSON result file.
const jsonRuntimeUri = runtimeFolder + "/.json"

// jsonSchemaUri is a constant string representing the URI of the JSON schema file.
const jsonSchemaUri = "https://raw.githubusercontent.com/GabrielHCataldo/gopen-gateway/main/json-schema.json"

func LoadDefaultEnvs() (err error) {
	if err = godotenv.Load("./.env"); helper.IsNotNil(err) {
		err = errors.New("Error load Gopen envs default:", err)
	}
	return err
}

func LoadEnvs(env string) (err error) {
	gopenEnvUri := getEnvUri(env)

	if err = godotenv.Load(gopenEnvUri); helper.IsNotNil(err) {
		err = errors.New("Error load Gopen envs from uri:", gopenEnvUri, "err:", err)
	}

	return err
}

func ReloadEnvs(env string) (err error) {
	gopenEnvUri := getEnvUri(env)

	if err = godotenv.Overload(gopenEnvUri); helper.IsNotNil(err) {
		err = errors.New("Error reload envs from uri:", gopenEnvUri, "err:", err)
	}

	return err
}

func LoadJson(env string) (*vo.GopenJson, error) {
	gopenJsonUri := getJsonUri(env)

	gopenJsonBytes, err := os.ReadFile(gopenJsonUri)
	if helper.IsNotNil(err) {
		return nil, errors.New("Error read Gopen config from file json:", gopenJsonUri, "err:", err)
	}
	gopenJsonBytes = fillEnvValues(gopenJsonBytes)

	if err = validateJsonBySchema(jsonSchemaUri, gopenJsonBytes); helper.IsNotNil(err) {
		return nil, err
	}

	return vo.NewGopenJson(gopenJsonBytes)
}

func WriteRuntimeJson(gopenJson *vo.GopenJson) error {
	if _, err := os.Stat(runtimeFolder); os.IsNotExist(err) {
		err = os.MkdirAll(runtimeFolder, 0755)
		if helper.IsNotNil(err) {
			return err
		}
	}

	gopenJsonBytes, err := json.MarshalIndent(gopenJson, "", "\t")
	if helper.IsNil(err) {
		err = os.WriteFile(jsonRuntimeUri, gopenJsonBytes, 0644)
	}

	return err
}

func RemoveRuntimeJson() error {
	err := os.Remove(runtimeFolder)
	if errors.IsNot(err, os.ErrNotExist) {
		err = nil
	}
	return err
}

func fillEnvValues(gopenJsonBytes []byte) []byte {
	// todo: aceitar campos não string receber variável de ambiente também
	//  foi pensado que talvez utilizar campos string e any para isso, convertendo para o tipo desejado apenas
	//  quando objeto de valor for montado
	gopenJsonStr := helper.SimpleConvertToString(gopenJsonBytes)

	regex := regexp.MustCompile(`\$\w+`)
	words := regex.FindAllString(gopenJsonStr, -1)

	count := 0
	for _, word := range words {
		envKey := strings.ReplaceAll(word, "$", "")
		envValue := os.Getenv(envKey)
		if helper.IsNotEmpty(envValue) {
			gopenJsonStr = strings.ReplaceAll(gopenJsonStr, word, envValue)
			count++
		}
	}

	return helper.SimpleConvertToBytes(gopenJsonStr)
}

func validateJsonBySchema(jsonSchemaUri string, jsonBytes []byte) error {
	schemaLoader := gojsonschema.NewReferenceLoader(jsonSchemaUri)
	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if helper.IsNotNil(err) {
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

func getEnvUri(env string) string {
	return fmt.Sprintf("./gopen/%s/.env", env)
}

func getJsonUri(env string) string {
	return fmt.Sprintf("./gopen/%s/.json", env)
}
