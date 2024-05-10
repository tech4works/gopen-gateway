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

package infra

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

// runtimeFolder is a constant string representing the filepath of the runtime folder for the Gopen application.
const runtimeFolder = "./runtime"

// jsonRuntimeUri is a constant string representing the filepath of the Gopen JSON result file.
const jsonRuntimeUri = runtimeFolder + "/.json"

// jsonProvider represents a JSON provider that implements the JsonProvider interface.
// It provides functionalities to read, validate, write, and remove JSON files.
type jsonProvider struct {
}

// NewJsonProvider creates a new instance of the JsonProvider interface.
func NewJsonProvider() interfaces.JsonProvider {
	return jsonProvider{}
}

// Read reads the contents of a file specified by the URI and returns them as a byte slice.
// It returns an error if the file cannot be read.
func (j jsonProvider) Read(uri string) ([]byte, error) {
	return os.ReadFile(uri)
}

// ValidateJsonBySchema validates a JSON document against a JSON schema specified by the URI.
// It takes the URI of the JSON schema and the JSON document as parameters.
// It returns an error if the validation fails, indicating that the JSON document does not conform to the schema.
// The error message includes a description of each validation error encountered.
// If the validation is successful, the error returned will be nil.
// The function uses the gojsonschema library to load the schema and document, perform the validation,
// and collect the validation errors, if any.
func (j jsonProvider) ValidateJsonBySchema(jsonSchemaUri string, jsonBytes []byte) error {
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

// WriteGopenJson writes the content of the GopenJson object to a JSON file.
// It first checks if the runtime folder exists, and if not, it creates the folder using os.MkdirAll.
// Then it marshals the GopenJson object into a byte slice using json.MarshalIndent.
// The encoding uses an empty prefix and a tab as an indentation suffix.
// Finally, it writes the byte slice to the JSON file using os.WriteFile.
// It returns an error if any of the above operations fail.
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

// RemoveGopenJson removes the GopenJson file from the runtime folder.
// It first tries to remove the runtime folder using os.Remove.
// If the folder does not exist, it returns immediately without any error.
// If the folder exists, it removes it and any file or folder inside it.
// It returns an error if any error occurs during the removal process, except when the folder does not exist.
// The error returned will be nil if the removal process is successful or if the folder does not exist.
func (j jsonProvider) RemoveGopenJson() error {
	err := os.Remove(runtimeFolder)
	if errors.IsNot(err, os.ErrNotExist) {
		err = nil
	}
	return nil
}
