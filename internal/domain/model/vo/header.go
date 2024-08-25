/*
 * Copyright 2024 Tech4Works
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

package vo

import (
	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/gopen-gateway/internal/domain/mapper"
	"net/http"
	"strings"
)

type Header struct {
	values map[string][]string
}

func NewHeader(values map[string][]string) Header {
	cleanValues := map[string][]string{}
	for k, v := range values {
		if checker.IsNotEmpty(v) {
			cleanValues[k] = v
		}
	}
	return Header{values: cleanValues}
}

func NewHeaderByBody(body *Body) Header {
	if checker.IsNil(body) {
		return Header{}
	}

	values := map[string][]string{}
	values[mapper.ContentType] = []string{body.ContentType().String()}
	values[mapper.ContentLength] = []string{body.SizeInString()}
	if body.HasContentEncoding() {
		values[mapper.ContentEncoding] = []string{body.ContentEncoding().String()}
	}
	return Header{values: values}
}

func (h Header) Http() http.Header {
	return h.Copy()
}

func (h Header) String() string {
	return converter.ToCompactString(h.values)
}

func (h Header) GetAll(key string) []string {
	return h.values[key]
}

func (h Header) Get(key string) string {
	valuesByKey := h.values[key]
	if checker.IsNotEmpty(valuesByKey) {
		return strings.Join(valuesByKey, ", ")
	}
	return ""
}

func (h Header) GetFirst(key string) string {
	valuesByKey := h.values[key]
	if checker.IsNotEmpty(valuesByKey) {
		return valuesByKey[0]
	}
	return ""
}

func (h Header) Exists(key string) bool {
	_, ok := h.values[key]
	return ok
}

func (h Header) NotExists(key string) bool {
	return !h.Exists(key)
}

func (h Header) Copy() map[string][]string {
	copiedValues := map[string][]string{}
	for key, value := range h.values {
		copiedValues[key] = value
	}
	return copiedValues
}

func (h Header) Map() any {
	return h.Copy()
}

func (h Header) Size() int {
	size := 0
	for key, values := range h.values {
		size += len(key) + 2
		for _, value := range values {
			size += len(value)
			size += 2
		}
		size -= 2
		size += 2
	}
	size += 2
	return size
}

func (h Header) SizeStr() string {
	bs := NewBytesByInt(h.Size())
	return converter.ToCompactString(bs.String())
}

func (h Header) Keys() (ss []string) {
	for key := range h.values {
		ss = append(ss, key)
	}
	return ss
}

func (h Header) MarshalJSON() ([]byte, error) {
	return converter.ToBytesWithErr(h.Map())
}

func (h *Header) UnmarshalJSON(data []byte) error {
	var values map[string][]string

	err := converter.ToDestWithErr(data, &values)
	if checker.NonNil(err) {
		return err
	}

	h.values = values
	return nil
}
