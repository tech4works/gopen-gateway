package infra

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
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/interfaces"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/xeipuuv/gojsonschema"
	"os"
)

const runtimeFolder = "./runtime"

// jsonRuntimeUri is a constant string representing the filepath of the Gopen JSON result file.
const jsonRuntimeUri = runtimeFolder + "/.json"

type jsonProvider struct {
}

func NewJsonProvider() interfaces.JsonProvider {
	return jsonProvider{}
}

func (j jsonProvider) Read(uri string) ([]byte, error) {
	return os.ReadFile(uri)
}

func (j jsonProvider) ValidateJsonBySchema(jsonSchemaUri string, jsonBytes []byte) error {
	// carregamos o schema e o documento
	schemaLoader := gojsonschema.NewReferenceLoader(jsonSchemaUri)
	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)
	// chamamos o validate
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if helper.IsNotNil(err) {
		return errors.New("Error validate schema:", err)
	}
	// checamos se valido, caso nao seja formatamos a mensagem
	if !result.Valid() {
		errorMsg := fmt.Sprintf("Map poorly formatted!\n")
		for _, desc := range result.Errors() {
			errorMsg += fmt.Sprintf("- %s\n", desc)
		}
		err = errors.New(errorMsg)
	}
	// retornamos o erro, se nao tiver, sera nil
	return err
}

func (j jsonProvider) WriteGopenJson(gopenJson *vo.GopenJson) error {
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

func (j jsonProvider) RemoveGopenJson() error {
	err := os.Remove(runtimeFolder)
	if errors.IsNot(err, os.ErrNotExist) {
		err = nil
	}
	return nil
}
