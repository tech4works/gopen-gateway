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
	"github.com/tech4works/errors"
)

type Payload struct {
	contentType     ContentType
	contentEncoding ContentEncoding
	buffer          *bytes.Buffer
}

func NewPayload(contentType, contentEncoding string, buffer *bytes.Buffer) *Payload {
	if checker.IsNil(buffer) || checker.IsEmpty(buffer.Bytes()) {
		return nil
	}
	return &Payload{
		contentType:     NewContentType(contentType),
		contentEncoding: NewContentEncoding(contentEncoding),
		buffer:          buffer,
	}
}

func NewPayloadWithContentType(contentType ContentType, buffer *bytes.Buffer) *Payload {
	return NewPayload(contentType.String(), "", buffer)
}

func NewPayloadJSON(buffer *bytes.Buffer) *Payload {
	return NewPayloadWithContentType(NewContentTypeJSON(), buffer)
}

func (p *Payload) ContentType() ContentType {
	return p.contentType
}

func (p *Payload) HasContentEncoding() bool {
	return checker.IsNotEmpty(p.contentEncoding)
}

func (p *Payload) ContentEncoding() ContentEncoding {
	return p.contentEncoding
}

func (p *Payload) IsValid() bool {
	return checker.NonNil(p.Validate())
}

func (p *Payload) IsNotValid() bool {
	return !p.IsValid()
}

func (p *Payload) Validate() error {
	bs, err := p.Bytes()
	if checker.NonNil(err) {
		return err
	}

	if (p.ContentType().IsJSON() && checker.IsNotJSON(bs)) || (p.ContentType().IsXML() && checker.IsNotXML(bs)) {
		return errors.Newf("payload failed: op=validate kind=%s", p.contentType)
	}

	return nil
}

func (p *Payload) Buffer() *bytes.Buffer {
	return bytes.NewBuffer(p.RawBytes())
}

func (p *Payload) Bytes() ([]byte, error) {
	if p.ContentEncoding().IsGzip() {
		return decompressor.ToBytesWithErr(decompressor.TypeGzip, p.RawBytes())
	} else if p.ContentEncoding().IsDeflate() {
		return decompressor.ToBytesWithErr(decompressor.TypeDeflate, p.RawBytes())
	} else {
		return p.RawBytes(), nil
	}
}

func (p *Payload) RawBytes() []byte {
	return p.buffer.Bytes()
}

func (p *Payload) String() (string, error) {
	bs, err := p.Bytes()
	if checker.NonNil(err) {
		return "", err
	} else {
		return string(bs), nil
	}
}

func (p *Payload) CompactString() (string, error) {
	bs, err := p.Bytes()
	if checker.NonNil(err) {
		return "", err
	} else {
		return converter.ToCompactStringWithErr(bs)
	}
}

func (p *Payload) Raw() (string, error) {
	s, err := p.String()
	if checker.NonNil(err) {
		return "", err
	}
	if p.ContentType().IsPlainText() {
		s = strconv.Quote(s)
	}
	return s, nil
}

func (p *Payload) Size() int {
	return len(p.RawBytes())
}

func (p *Payload) SizeInString() string {
	return converter.ToString(p.Size())
}

func (p *Payload) SizeInByteUnit() string {
	bs := NewBytesByInt(p.Size())
	return bs.String()
}

func (p *Payload) Map() (any, error) {
	if !p.ContentType().IsJSON() {
		return nil, nil
	}

	str, err := p.String()
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

func (p *Payload) MarshalJSON() ([]byte, error) {
	base64, err := converter.ToBase64WithErr(p.buffer.Bytes())
	if checker.NonNil(err) {
		return nil, err
	}

	return json.Marshal(map[string]any{
		"content-type":     p.contentType,
		"content-encoding": p.contentEncoding,
		"buffer":           base64,
	})
}

func (p *Payload) UnmarshalJSON(data []byte) error {
	var mapPayload map[string]any

	err := json.Unmarshal(data, &mapPayload)
	if checker.NonNil(err) {
		return err
	}

	if value, exists := mapPayload["buffer"]; exists {
		valueStr, ok := value.(string)
		if !ok {
			return fmt.Errorf("buffer not a string")
		}

		bufferStr, err := converter.FromBase64ToStringWithErr(valueStr)
		if checker.NonNil(err) {
			return err
		}

		p.buffer = bytes.NewBufferString(bufferStr)
	}

	if value, exists := mapPayload["content-type"]; exists {
		str, ok := value.(string)
		if ok {
			p.contentType = NewContentType(str)
		}
	}

	if value, exists := mapPayload["content-encoding"]; exists {
		str, ok := value.(string)
		if ok {
			p.contentEncoding = NewContentEncoding(str)
		}
	}

	return nil
}
