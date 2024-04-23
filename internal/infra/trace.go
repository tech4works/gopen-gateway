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
	"fmt"
	"github.com/google/uuid"
	"time"
)

type trace struct {
}

// TraceProvider is an interface that defines the behavior for generating a trace ID.
type TraceProvider interface {
	// GenerateTraceId generates a trace ID for tracking requests.
	// It returns a string representing the generated trace ID.
	// The implementation of this method may vary depending on the requirements and system configuration.
	GenerateTraceId() string
}

// NewTraceProvider creates a new instance of trace that implements the TraceProvider interface.
func NewTraceProvider() TraceProvider {
	return trace{}
}

// GenerateTraceId generates a unique trace ID by combining a UUID and the current UnixNano timestamp.
// It uses the UUID package to generate a random UUID. It then retrieves the current UnixNano timestamp
// using the time package. The UUID string is concatenated with the timestamp and truncated to 16 characters
// using the fmt package.
// The resulting trace ID is returned as a string.
// Example:
//
//	t := trace{}
//	traceId := t.GenerateTraceId()
//	fmt.Println(traceId)
//
// Output:
//
//	4ae6c92d16089e521626
func (t trace) GenerateTraceId() string {
	u := uuid.New().String()
	unixNano := time.Now().UnixNano()
	return fmt.Sprintf("%s%d", u[:8], unixNano)[:16]
}
