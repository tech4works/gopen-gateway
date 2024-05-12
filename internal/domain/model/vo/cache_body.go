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

package vo

import (
	"bytes"
	"github.com/GabrielHCataldo/go-helper/helper"
	"strconv"
)

// CacheBody represents the caching value of an HTTP httpResponse body.
type CacheBody struct {
	// ContentType represents the format of the content.
	ContentType ContentType `json:"content-type,omitempty"`
	// Value represents the caching content of an HTTP httpResponse body.
	// It is a pointer to the CacheBodyValue type, which is an alias for bytes.Buffer.
	// The value is nullable and is omitted in JSON if it is empty.
	Value *CacheBodyValue `json:"value,omitempty"`
}

// CacheBodyValue is an alias for bytes.Buffer type used to represent the caching value
// of an HTTP httpResponse body. It contains methods to convert the value to different
// representations, such as string and JSON.
type CacheBodyValue bytes.Buffer

// newCacheBody creates a new instance of CacheBody based on the provided body.
// If the body is nil, it returns nil.
// Otherwise, it sets the ContentType field of CacheBody based on the ContentType method of body.
// It sets the Value field of CacheBody by calling newCacheBodyValue with the Value method of body as an argument.
// It returns a pointer to the constructed CacheBody instance.
func newCacheBody(body *Body) *CacheBody {
	if helper.IsNil(body) {
		return nil
	}
	return &CacheBody{
		ContentType: body.ContentType(),
		Value:       newCacheBodyValue(body.Buffer()),
	}
}

// newCacheBodyValue creates a new instance of CacheBodyValue based on the provided buffer.
// If the buffer is nil or empty, it returns nil.
// Otherwise, it returns a pointer to the CacheBodyValue instance, casting the buffer to the CacheBodyValue type.
// CacheBodyValue is an alias for bytes.Buffer type used to represent the caching value of an HTTP httpResponse body.
// It contains methods to convert the value to different representations, such as string and JSON.
//
// The CacheBodyValue type has the following methods:
// - String(): returns the string representation of the CacheBodyValue instance.
// - Bytes(): returns the byte slice representation of the CacheBodyValue instance.
// - MarshalJSON(): returns the JSON encoding of the CacheBodyValue instance.
// - UnmarshalJSON(data []byte): decodes the JSON data into a string and writes it to the CacheBodyValue instance.
//
// Example usage:
//
//	cacheBody := newCacheBody(body)
//	value := newCacheBodyValue(body.Buffer())
func newCacheBodyValue(buffer *bytes.Buffer) *CacheBodyValue {
	if helper.IsNil(buffer) || helper.IsEmpty(buffer.Bytes()) {
		return nil
	}
	return (*CacheBodyValue)(buffer)
}

// String returns the string representation of the CacheBodyValue instance.
// It calls the String method of the underlying bytes.Buffer type to get the string representation.
func (c *CacheBodyValue) String() string {
	return (*bytes.Buffer)(c).String()
}

// Bytes returns the byte slice representation of the CacheBodyValue instance.
// It calls the Bytes method of the underlying bytes.Buffer type to get the byte slice representation.
func (c *CacheBodyValue) Bytes() []byte {
	return (*bytes.Buffer)(c).Bytes()
}

func (c *CacheBodyValue) MarshalJSON() ([]byte, error) {
	b64, err := helper.CompressWithGzipToBase64(c.Bytes())
	if helper.IsNotNil(err) {
		return nil, err
	}
	b64 = strconv.Quote(b64)
	return helper.SimpleConvertToBytes(b64), nil
}

func (c *CacheBodyValue) UnmarshalJSON(data []byte) error {
	if helper.IsEmpty(data) {
		return nil
	}

	unquote, _ := strconv.Unquote(string(data))
	bs, err := helper.DecompressFromBase64WithGzip(unquote)
	if helper.IsNotNil(err) {
		return err
	}
	*c = *(*CacheBodyValue)(bytes.NewBuffer(bs))

	return nil
}
