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
	"encoding/json"
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/xeipuuv/gojsonschema"
	"os"
	"regexp"
	"strings"
)

// jsonResultUri is a constant string representing the filepath of the Gopen JSON result file.
const jsonResultUri = "./runtime/.json"

// jsonSchemaUri is a constant string representing the URI of the JSON schema file.
const jsonSchemaUri = "https://raw.githubusercontent.com/GabrielHCataldo/gopen-gateway/main/json-schema.json"

// LoadGopenJson loads the Gopen configuration JSON file based on the specified environment.
// The function takes an environment string as input and returns a pointer to a GopenJson object.
// It first gets the file URI for the JSON file using the getFileJsonUri function.
// Then, it reads the file contents and fills the environment variable values using the fillEnvValues function.
// After that, it validates the JSON file against a schema using the validateJsonBySchema function.
// If the JSON is valid, it creates a new GopenJson object using the NewGopenJson function.
// If any error occurs during the process, it panics and logs an error message.
// Finally, it returns the created GopenJson object.
// Note: This function assumes the existence of the getFileJsonUri, PrintInfoLogCmdf, fillEnvValues,
// validateJsonBySchema, and NewGopenJson functions.
func LoadGopenJson(env string) *vo.GopenJson {
	// carregamos o arquivo de json de configuração do Gopen
	fileJsonUri := getFileJsonUri(env)
	PrintInfoLogCmdf("Loading Gopen json from file: %s...", fileJsonUri)
	fileJsonBytes, err := os.ReadFile(fileJsonUri)
	if helper.IsNotNil(err) {
		panic(errors.New("Error read Gopen config from file json:", fileJsonUri, "err:", err))
	}

	// preenchemos os valores de variável de ambiente com a sintaxe pre-definida
	fileJsonBytes = fillEnvValues(fileJsonBytes)

	// validamos o schema
	if err = validateJsonBySchema(fileJsonUri, fileJsonBytes); helper.IsNotNil(err) {
		panic(err)
	}

	// construímos o objeto de valor relacionado ao json de configuração
	gopenJson, err := vo.NewGopenJson(fileJsonBytes)
	if helper.IsNotNil(err) {
		panic(err)
	}

	// se tudo ocorreu bem, retornamos
	return gopenJson
}

// fillEnvValues fills the environment variables in a JSON string using the $word syntax.
// The function takes a byte array of the JSON string as input and returns the modified JSON string as a byte array.
// It searches for all occurrences of $word in the JSON string and replaces them with the corresponding environment variable value.
// The function uses regular expressions to find all $word occurrences and os.Getenv() to get the environment variable value.
// If a valid value is found, it replaces the $word with the value in the JSON string.
// The function prints the number of environment variable values found and successfully filled during the process.
// It also uses the helper functions 'SimpleConvertToString' and 'SimpleConvertToBytes' to convert the byte array to string and vice versa.
func fillEnvValues(gopenBytesJson []byte) []byte {
	// todo: aceitar campos não string receber variável de ambiente também
	//  foi pensado que talvez utilizar campos string e any para isso, convertendo para o tipo desejado apenas
	//  quando objeto de valor for montado

	PrintInfoLogCmd("Filling environment variables with $word syntax..")

	// convertemos os bytes do gopen json em string
	gopenStrJson := helper.SimpleConvertToString(gopenBytesJson)

	// compilamos o regex indicando um valor de env $API_KEY por exemplo
	regex := regexp.MustCompile(`\$\w+`)
	// damos o find pelo regex
	words := regex.FindAllString(gopenStrJson, -1)

	// imprimimos todas as palavras encontradas a ser preenchidas
	PrintInfoLogCmd(len(words), "environment variable values were found to fill in!")

	// inicializamos o contador de valores processados
	count := 0
	for _, word := range words {
		// replace do valor padrão $
		envKey := strings.ReplaceAll(word, "$", "")
		// obtemos o valor da env pela chave indicada
		envValue := os.Getenv(envKey)
		// caso valor encontrado, damos o replace da palavra encontrada pelo valor
		if helper.IsNotEmpty(envValue) {
			gopenStrJson = strings.ReplaceAll(gopenStrJson, word, envValue)
			count++
		}
	}
	// imprimimos a quantidade de envs preenchidas
	PrintInfoLogCmd(count, "environment variables successfully filled!")

	// convertemos esse novo
	return helper.SimpleConvertToBytes(gopenStrJson)
}

// WriteGopenJsonResult writes the GopenJson object to a JSON file.
// The function takes a pointer to a GopenJson object as input.
// It converts the GopenJson object to a byte array using json.MarshalIndent.
// If no error occurs during the marshaling process, it writes the byte array to the gopenJsonResult file using os.WriteFile.
// If an error occurs during the writing process, it logs a warning message using PrintWarningLogCmdf.
// Note: The global constant gopenJsonResult represents the filepath of the Gopen JSON result file.
func WriteGopenJsonResult(gopenJson *vo.GopenJson) {
	gopenBytes, err := json.MarshalIndent(gopenJson, "", "\t")
	if helper.IsNil(err) {
		err = os.WriteFile(jsonResultUri, gopenBytes, 0644)
	}
	if helper.IsNotNil(err) {
		PrintWarningLogCmdf("Error write file %s result: %s", jsonResultUri, err)
	}
}

// RemoveGopenJsonResult handles the removal of a JSON result file.
// This function will try to delete the file specified by the constant gopenJsonResult.
// If the file does not exist, the function will exit silently.
// If there is an error during removal that is NOT due to the file not existing,
// it logs a warning message with printWarningLogf.
func RemoveGopenJsonResult() {
	err := os.Remove(jsonResultUri)
	if helper.IsNotNil(err) && errors.IsNot(err, os.ErrNotExist) {
		PrintWarningLogCmdf("Error remove %s err: %s", jsonResultUri, err)
		return
	}
}

// getFileEnvUri returns the file URI for the given environment.
// The returned URI follows the format "./gopen/{env}.env".
func getFileEnvUri(env string) string {
	return fmt.Sprintf("./gopen/%s/.env", env)
}

// getFileJsonUri returns the file URI for the specified environment's JSON file.
// The returned URI follows the format "./gopen/{env}.json".
func getFileJsonUri(env string) string {
	return fmt.Sprintf("./gopen/%s/.json", env)
}

// validateJsonBySchema validates a JSON file against a given schema.
// It takes the file JSON URI and the file JSON bytes as inputs.
// The function starts by printing a log message indicating that it is validating the file schema.
// Then, it loads the schema and the document using the gojsonschema package.
// After that, it calls the Validate function to perform the schema validation.
// If there is an error while validating the schema, the function panics and logs an error message.
// If the file JSON is poorly formatted and does not pass the schema validation, the function constructs an error message with the filename and the validation errors.
// The error message is then returned as an error.
// If the JSON is valid and passes the schema validation, the function returns nil.
func validateJsonBySchema(fileJsonUri string, fileJsonBytes []byte) error {
	PrintInfoLogCmdf("Validating the %s file schema...", fileJsonUri)

	// carregamos o schema e o documento
	schemaLoader := gojsonschema.NewReferenceLoader(jsonSchemaUri)
	documentLoader := gojsonschema.NewBytesLoader(fileJsonBytes)

	// chamamos o validate
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if helper.IsNotNil(err) {
		panic(errors.New("Error validate schema:", err))
	}

	// checamos se valido, caso nao seja formatamos a mensagem
	if !result.Valid() {
		errorMsg := fmt.Sprintf("Json %s poorly formatted!\n", fileJsonUri)
		for _, desc := range result.Errors() {
			errorMsg += fmt.Sprintf("- %s\n", desc)
		}
		return errors.New(errorMsg)
	}
	// se tudo ocorrem bem retornamos nil
	return nil
}
