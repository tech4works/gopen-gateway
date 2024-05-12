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

package interfaces

import (
	"context"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"net/http"
)

// CmdLoggerProvider is an interface that defines methods for printing different types of log messages
// such as logo, titles, info, and warnings.
type CmdLoggerProvider interface {
	// PrintLogo is a method that prints the logo of the software along with the provided version string.
	PrintLogo(version string)
	// PrintTitle prints a title message using the provided string.
	PrintTitle(title string)
	// PrintInfo prints an informational message with the provided arguments.
	// The message can be of any type, and multiple arguments can be provided.
	PrintInfo(msg ...any)
	// PrintInfof prints an informational message with the provided arguments, using a format string.
	// Format string can contain placeholders for the arguments.
	// Arguments can be of any type, and multiple arguments can be provided.
	PrintInfof(format string, msg ...any)
	// PrintWarning prints a warning message with the provided arguments.
	// The message can be of any type, and multiple arguments can be provided.
	PrintWarning(msg ...any)
	// PrintWarningf prints a warning message with the provided arguments and a format string.
	// Format string can contain placeholders for the arguments.
	// Arguments can be of any type, and multiple arguments can be provided.
	PrintWarningf(format string, msg ...any)
}

// JsonProvider is an interface that provides methods for reading, validating, writing, and removing JSON data.
type JsonProvider interface {
	// Read reads the content of a file specified by the given URI and returns it as a byte slice.
	// It takes the URI of the file as a parameter and returns the content of the file as a byte slice.
	// If there is an error while reading the file, it returns an error.
	// The URI should be a string representing the location of the file.
	// The returned byte slice represents the content of the file.
	Read(uri string) ([]byte, error)
	// ValidateJsonBySchema validates a JSON object against a JSON schema.
	// It takes the URI of the JSON schema and the JSON object as parameters.
	// If the JSON object does not conform to the schema, it returns an error.
	// The JSON schema URI should be a string representing the location of the schema.
	// The JSON bytes represent the content of the JSON object.
	ValidateJsonBySchema(jsonSchemaUri string, jsonBytes []byte) error
	// WriteGopenJson writes the content of the GopenJson object to a runtime JSON file.
	// It takes a pointer to the GopenJson object as a parameter.
	// If there is an error while writing the file, it returns an error.
	WriteGopenJson(gopenJson *vo.GopenJson) error
	// RemoveGopenJson removes the runtime JSON file.
	// It reads the URI of the file from the jsonProvider and deletes the file from the file system.
	// If there is an error while removing the file, it returns an error.
	// The URI should be a string representing the location of the file.
	RemoveGopenJson() error
}

// CacheStore is an interface that defines methods for interacting with a cache store.
type CacheStore interface {
	Set(ctx context.Context, key string, value *vo.CacheResponse) error
	// Del deletes the cache entry with the specified key from the cache store.
	// The key is a string that serves as the identifier of the cache entry.
	// The error returned indicates any issues encountered while deleting the cache entry.
	// This method is used to remove a key-value pair from the cache store.
	// Implementing the CacheStore interface, this method is responsible for removing the cache entry using the underlying implementation.
	Del(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (*vo.CacheResponse, error)
	// Close is a method defined in the CacheStore interface. It is used to close the cache store and release any resources
	// associated with it. The method returns an error if there was a problem closing the store.
	Close() error
}

// CacheStoreProvider is an interface that defines methods for creating cache store instances.
type CacheStoreProvider interface {
	// Memory returns a new instance of the MemoryStore structure that implements the CacheStore interface.
	// This implementation uses an in-memory cache with a time-to-live (TTL).
	Memory() CacheStore
	// Redis returns a new instance of the RedisStore structure that implements
	// the CacheStore interface. This implementation uses a Redis cache with the given address and password.
	Redis(address, password string) CacheStore
}

// RestTemplate is an interface that represents a template for making HTTP requests.
// It provides a method MakeRequest for sending an HTTP request and returning the corresponding
// HTTP response and an error, if any.
type RestTemplate interface {
	// MakeRequest sends an HTTP request and returns the corresponding HTTP response and error.
	// It takes an HTTP request object as a parameter.
	// The function's steps are as follows:
	//
	//  1. The function sends the HTTP request using a REST client.
	//  2. If the operation fails, the function returns an error.
	//  3. Otherwise, it returns the HTTP response.
	//
	// Parameters:
	// httpRequest: the HTTP request to be sent.
	//
	// Returns:
	// An HTTP response object and an error.
	MakeRequest(netHttpRequest *http.Request) (*http.Response, error)
}
