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
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/decompressor"
)

type Body struct {
	contentType     ContentType
	contentEncoding ContentEncoding
	buffer          *bytes.Buffer
}

func NewBody(contentType, contentEncoding string, buffer *bytes.Buffer) *Body {
	if checker.IsNil(buffer) || checker.IsEmpty(buffer.Bytes()) {
		return nil
	}
	return &Body{
		contentType:     NewContentType(contentType),
		contentEncoding: NewContentEncoding(contentEncoding),
		buffer:          buffer,
	}
}

func NewBodyWithContentType(contentType ContentType, buffer *bytes.Buffer) *Body {
	if checker.IsNil(buffer) || checker.IsEmpty(buffer.Bytes()) {
		return nil
	}
	return &Body{
		contentType: contentType,
		buffer:      buffer,
	}
}

func NewBodyJson(buffer *bytes.Buffer) *Body {
	if checker.IsNil(buffer) || checker.IsEmpty(buffer.Bytes()) {
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
		return decompressor.ToBytesWithErr(decompressor.TypeGzip, b.RawBytes())
	} else if b.ContentEncoding().IsDeflate() {
		return decompressor.ToBytesWithErr(decompressor.TypeDeflate, b.RawBytes())
	}
	return b.RawBytes(), nil
}

func (b *Body) RawBytes() []byte {
	return b.buffer.Bytes()
}

func (b *Body) String() (string, error) {
	bs, err := b.Bytes()
	if checker.NonNil(err) {
		return string(b.RawBytes()), err
	}
	return string(bs), nil
}

func (b *Body) CompactString() (string, error) {
	bs, err := b.Bytes()
	if checker.NonNil(err) {
		bs = b.RawBytes()
	}
	return converter.ToCompactString(bs), err
}

func (b *Body) Raw() (string, error) {
	s, err := b.String()
	if checker.NonNil(err) {
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
		return converter.ToCompactString(s)
	}
	return fmt.Sprintf("type=%s encoding=%s contentLength=%s", b.contentType, b.contentEncoding, b.SizeInByteUnit())
}

func (b *Body) Size() int {
	return len(b.RawBytes())
}

func (b *Body) SizeInString() string {
	return converter.ToString(b.Size())
}

func (b *Body) SizeInByteUnit() string {
	bs := NewBytesByInt(b.Size())
	return bs.String()
}

func (b *Body) Map() (any, error) {
	if b.ContentType().IsNotJSON() {
		return nil, nil
	}

	str, err := b.String()
	if checker.NonNil(err) {
		return nil, err
	}

	if checker.IsSlice(str) {
		var dest []any
		err = converter.ToDestWithErr(str, &dest)
		if checker.NonNil(err) {
			return nil, err
		}
		return dest, nil
	}

	var dest map[string]any
	err = converter.ToDestWithErr(str, &dest)
	if checker.NonNil(err) {
		return nil, err
	}
	return dest, nil
}

func (b *Body) MarshalJSON() ([]byte, error) {
	base64, err := converter.ToBase64WithErr(b.buffer.Bytes())
	if checker.NonNil(err) {
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
	if checker.NonNil(err) {
		return err
	}

	if value, exists := mapBody["buffer"]; exists {
		valueStr, ok := value.(string)
		if !ok {
			return fmt.Errorf("buffer not a string")
		}

		bufferStr, err := converter.FromBase64ToStringWithErr(valueStr)
		if checker.NonNil(err) {
			return err
		}

		b.buffer = bytes.NewBufferString(bufferStr)
	}

	if value, exists := mapBody["contentType"]; exists {
		str, ok := value.(string)
		if ok {
			b.contentType = NewContentType(str)
		}
	}

	if value, exists := mapBody["contentEncoding"]; exists {
		str, ok := value.(string)
		if ok {
			b.contentEncoding = NewContentEncoding(str)
		}
	}

	return nil
}
