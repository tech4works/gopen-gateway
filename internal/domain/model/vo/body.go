/*
 * Copyright 2024 Gabriel Cataldo
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a Copy of the License at
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
	"encoding/json"
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"strconv"
)

type Body struct {
	contentType     ContentType
	contentEncoding ContentEncoding
	buffer          *bytes.Buffer
}

func NewBody(contentType, contentEncoding string, buffer *bytes.Buffer) *Body {
	if helper.IsNil(buffer) || helper.IsEmpty(buffer.Bytes()) {
		return nil
	}
	return &Body{
		contentType:     NewContentType(contentType),
		contentEncoding: NewContentEncoding(contentEncoding),
		buffer:          buffer,
	}
}

func NewBodyWithContentType(contentType ContentType, buffer *bytes.Buffer) *Body {
	if helper.IsNil(buffer) || helper.IsEmpty(buffer.Bytes()) {
		return nil
	}
	return &Body{
		contentType: contentType,
		buffer:      buffer,
	}
}

func NewBodyJson(buffer *bytes.Buffer) *Body {
	if helper.IsNil(buffer) || helper.IsEmpty(buffer.Bytes()) {
		return nil
	}
	return &Body{
		contentType: NewContentTypeJson(),
		buffer:      buffer,
	}
}

func (b *Body) ContentType() ContentType {
	return b.contentType
}

func (b *Body) ContentEncoding() ContentEncoding {
	return b.contentEncoding
}

func (b *Body) HasContentEncoding() bool {
	return b.contentEncoding.Valid()
}

func (b *Body) Buffer() *bytes.Buffer {
	return bytes.NewBuffer(b.RawBytes())
}

func (b *Body) Bytes() ([]byte, error) {
	if b.ContentEncoding().IsGzip() {
		return helper.DecompressWithGzip(b.RawBytes())
	} else if b.ContentEncoding().IsDeflate() {
		return helper.DecompressWithDeflate(b.RawBytes())
	}
	return b.RawBytes(), nil
}

func (b *Body) RawBytes() []byte {
	return b.buffer.Bytes()
}

func (b *Body) String() (string, error) {
	bs, err := b.Bytes()
	if helper.IsNotNil(err) {
		return string(b.RawBytes()), err
	}
	return string(bs), nil
}

func (b *Body) Raw() (string, error) {
	s, err := b.String()
	if helper.IsNotNil(err) {
		return "", err
	}
	if b.ContentType().IsText() {
		s = strconv.Quote(s)
	}
	return s, nil
}

func (b *Body) Resume() string {
	if b.contentType.IsJSON() || b.contentType.IsXML() || b.contentType.IsText() {
		s, _ := b.String()
		return helper.SimpleCompactString(s)
	}
	return fmt.Sprintf("contentType=%s contentLenght=%s contentEncoding=%s", b.contentType.String(),
		b.contentEncoding.String(), b.LenStr())
}

func (b *Body) Len() int {
	return len(b.RawBytes())
}

func (b *Body) LenStr() string {
	return helper.SimpleConvertToString(b.Len())
}

func (b *Body) Map() (any, error) {
	if b.ContentType().IsNotJSON() {
		return nil, nil
	}
	bs, err := b.Bytes()
	if helper.IsNotNil(err) {
		return nil, err
	}

	var dest any
	err = helper.ConvertToDest(bs, &dest)
	if helper.IsNotNil(err) {
		return nil, err
	}
	return dest, nil
}

func (b *Body) MarshalJSON() ([]byte, error) {
	base64, err := helper.ConvertToBase64(b.buffer.Bytes())
	if helper.IsNotNil(err) {
		return nil, err
	}
	return json.Marshal(map[string]any{
		"contentType":     b.contentType,
		"contentEncoding": b.contentEncoding,
		"buffer":          base64,
	})
}

func (b *Body) UnmarshalJSON(data []byte) error {
	var mapBody map[string]any

	err := json.Unmarshal(data, &mapBody)
	if helper.IsNotNil(err) {
		return err
	}

	if value, exists := mapBody["buffer"]; exists {
		valueStr, err := helper.ConvertToString(value)
		if helper.IsNotNil(err) {
			return err
		}
		bufferStr, err := helper.ConvertBase64ToString(valueStr)
		if helper.IsNotNil(err) {
			return err
		}
		b.buffer = bytes.NewBufferString(bufferStr)
	}

	if value, exists := mapBody["contentType"]; exists {
		contentType, ok := value.(ContentType)
		if ok {
			b.contentType = contentType
		}
	}
	if value, exists := mapBody["contentEncoding"]; exists {
		contentEncoding, ok := value.(ContentEncoding)
		if ok {
			b.contentEncoding = contentEncoding
		}
	}

	return nil
}
