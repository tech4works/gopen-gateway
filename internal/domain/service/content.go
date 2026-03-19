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

package service

import (
	"strconv"

	"github.com/tech4works/checker"
	"github.com/tech4works/compressor"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type content struct {
	converter domain.Converter
}

type Content interface {
	ModifyPayloadContentType(config vo.ContentType, payload *vo.Payload) (*vo.Payload, error)
	ModifyPayloadContentEncoding(config vo.ContentEncoding, payload *vo.Payload) (*vo.Payload, error)
}

func NewContent(converter domain.Converter) Content {
	return content{
		converter: converter,
	}
}

func (c content) ModifyPayloadContentType(config vo.ContentType, payload *vo.Payload) (*vo.Payload, error) {
	if checker.IsNil(payload) || config.IsUnsupported() || config.Equals(payload.ContentType()) {
		return payload, nil
	}

	bs, err := c.modifyPayloadContentType(config, payload)
	if checker.NonNil(err) {
		return payload, errors.Inheritf(err, "content-type failed: op=modify-payload from=%s to=%s",
			payload.ContentType(), config)
	}

	buffer, err := converter.ToBufferWithErr(bs)
	if checker.NonNil(err) {
		return payload, errors.Inheritf(err, "content-type failed: op=buffer to=%s", config)
	}

	return vo.NewPayloadWithContentType(config, buffer), nil
}

func (c content) ModifyPayloadContentEncoding(config vo.ContentEncoding, payload *vo.Payload) (
	*vo.Payload, error) {
	if checker.IsNil(payload) || config.IsUnsupported() || config.Equals(payload.ContentEncoding()) {
		return payload, nil
	}

	bs, err := c.modifyPayloadContentEncoding(config, payload)
	if checker.NonNil(err) {
		return payload, errors.Inheritf(err, "content-encoding failed: op=modify-payload from=%s to=%s",
			payload.ContentEncoding(), config)
	}

	buffer, err := converter.ToBufferWithErr(bs)
	if checker.NonNil(err) {
		return payload, errors.Inheritf(err, "content-encoding failed: op=buffer to=%s", config.String())
	}

	return vo.NewPayload(payload.ContentType().String(), config.String(), buffer), nil
}

func (c content) modifyPayloadContentType(config vo.ContentType, payload *vo.Payload) ([]byte, error) {
	rawBytes, err := payload.Bytes()
	if checker.NonNil(err) {
		return nil, errors.Inherit(err, "content-type failed: op=bytes")
	} else if config.IsPlainText() {
		return c.asPlainText(rawBytes)
	} else if config.IsJSON() {
		return c.asJSON(payload, rawBytes)
	} else if config.IsXML() {
		return c.asXML(payload, rawBytes)
	} else {
		return nil, errors.New("content-type failed: op=unsupported")
	}
}

func (c content) modifyPayloadContentEncoding(config vo.ContentEncoding, payload *vo.Payload) ([]byte,
	error) {
	rawBytes, err := payload.Bytes()
	if checker.NonNil(err) {
		return nil, errors.Inherit(err, "content-encoding failed: op=bytes")
	} else if config.IsGzip() {
		return c.asGzip(rawBytes)
	} else if config.IsDeflate() {
		return c.asDeflate(rawBytes)
	} else {
		return nil, errors.New("content-encoding failed: op=unsupported")
	}
}

func (c content) asPlainText(rawBytes []byte) ([]byte, error) {
	return []byte(strconv.Quote(string(rawBytes))), nil
}

func (c content) asJSON(payload *vo.Payload, rawBytes []byte) ([]byte, error) {
	if payload.ContentType().IsXML() {
		converted, err := c.converter.ConvertXMLToJSON(rawBytes)
		if checker.NonNil(err) {
			return nil, errors.Inherit(err, "content-type failed: op=convert xml->json")
		}
		return converted, nil
	}

	converted, err := c.converter.ConvertTextToJSON(rawBytes)
	if checker.NonNil(err) {
		return nil, errors.Inherit(err, "content-type failed: op=convert text->json")
	}

	return converted, nil
}

func (c content) asXML(payload *vo.Payload, rawBytes []byte) ([]byte, error) {
	if payload.ContentType().IsJSON() {
		converted, err := c.converter.ConvertJSONToXML(rawBytes)
		if checker.NonNil(err) {
			return nil, errors.Inherit(err, "content-type failed: op=convert json->xml")
		}
		return converted, nil
	}

	converted, err := c.converter.ConvertTextToXML(rawBytes)
	if checker.NonNil(err) {
		return nil, errors.Inherit(err, "content-type failed: op=convert text->xml")
	}

	return converted, nil
}

func (c content) asGzip(rawBytes []byte) ([]byte, error) {
	compressed, err := compressor.ToGzipWithErr(rawBytes)
	if checker.NonNil(err) {
		return nil, errors.Inherit(err, "content-encoding failed: op=gzip")
	}
	return compressed, nil
}

func (c content) asDeflate(rawBytes []byte) ([]byte, error) {
	compressed, err := compressor.ToDeflateWithErr(rawBytes)
	if checker.NonNil(err) {
		return nil, errors.Inherit(err, "content-encoding failed: op=deflate")
	}
	return compressed, nil
}
