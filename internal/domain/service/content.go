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
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
)

type contentService struct {
	converter domain.Converter
}

type Content interface {
	ModifyBodyContentType(contentType enum.ContentType, body *vo.Body) (*vo.Body, error)
	ModifyBodyContentEncoding(contentEncoding enum.ContentEncoding, body *vo.Body) (*vo.Body, error)
}

func NewContent(converter domain.Converter) Content {
	return contentService{
		converter: converter,
	}
}

func (c contentService) ModifyBodyContentType(contentType enum.ContentType, body *vo.Body) (*vo.Body, error) {
	if !contentType.IsEnumValid() ||
		body.ContentType().IsUnknown() ||
		checker.Equals(body.ContentType().ToEnum(), contentType) {
		return body, nil
	}

	bodyBytes, httpContentType, err := c.modifyBodyContentType(body, contentType)
	if checker.NonNil(err) {
		return body, errors.Inheritf(err, "content-type failed: op=modify-body from=%s to=%s", body.ContentType(),
			contentType)
	}

	buffer, err := converter.ToBufferWithErr(bodyBytes)
	if checker.NonNil(err) {
		return body, errors.Inheritf(err, "content-type failed: op=buffer to=%s", httpContentType)
	}

	return vo.NewBodyWithContentType(httpContentType, buffer), nil
}

func (c contentService) ModifyBodyContentEncoding(contentEncoding enum.ContentEncoding, body *vo.Body) (*vo.Body, error) {
	if !contentEncoding.IsEnumValid() ||
		checker.EqualsIgnoreCase(body.ContentEncoding(), contentEncoding) {
		return body, nil
	}

	bodyBytes, httpContentEncoding, err := c.modifyBodyContentEncoding(body, contentEncoding)
	if checker.NonNil(err) {
		return body, errors.Inheritf(err,
			"content-encoding failed: op=modify-body from=%s to=%s", body.ContentEncoding(), contentEncoding)
	}

	buffer, err := converter.ToBufferWithErr(bodyBytes)
	if checker.NonNil(err) {
		return body, errors.Inheritf(err, "content-encoding failed: op=buffer to=%s", httpContentEncoding.String())
	}

	return vo.NewBody(
		body.ContentType().String(),
		httpContentEncoding.String(),
		buffer,
	), nil
}

func (c contentService) modifyBodyContentType(
	body *vo.Body,
	contentType enum.ContentType,
) ([]byte, vo.ContentType, error) {
	rawBytes, err := body.Bytes()
	if checker.NonNil(err) {
		return nil, "", errors.Inherit(err, "content-type failed: op=bytes")
	}

	switch contentType {
	case enum.ContentTypePlainText:
		return c.asPlainText(rawBytes)
	case enum.ContentTypeJson:
		return c.asJSON(body, rawBytes)
	case enum.ContentTypeXml:
		return c.asXML(body, rawBytes)
	default:
		return rawBytes, body.ContentType(), nil
	}
}

func (c contentService) modifyBodyContentEncoding(
	body *vo.Body,
	contentEncoding enum.ContentEncoding,
) ([]byte, vo.ContentEncoding, error) {
	rawBytes, err := body.Bytes()
	if checker.NonNil(err) {
		return nil, "", errors.Inherit(err, "content-encoding failed: op=bytes")
	}

	switch contentEncoding {
	case enum.ContentEncodingGzip:
		return c.asGzip(rawBytes)
	case enum.ContentEncodingDeflate:
		return c.asDeflate(rawBytes)
	default:
		return rawBytes, body.ContentEncoding(), nil
	}
}

func (c contentService) asPlainText(rawBytes []byte) ([]byte, vo.ContentType, error) {
	return []byte(strconv.Quote(string(rawBytes))), vo.NewContentTypeTextPlain(), nil
}

func (c contentService) asJSON(body *vo.Body, rawBytes []byte) ([]byte, vo.ContentType, error) {
	httpContentType := vo.NewContentTypeJson()

	if body.ContentType().IsXML() {
		converted, err := c.converter.ConvertXMLToJSON(rawBytes)
		if checker.NonNil(err) {
			return nil, "", errors.Inherit(err, "content-type failed: op=convert xml->json")
		}
		return converted, httpContentType, nil
	}

	converted, err := c.converter.ConvertTextToJSON(rawBytes)
	if checker.NonNil(err) {
		return nil, "", errors.Inherit(err, "content-type failed: op=convert text->json")
	}

	return converted, httpContentType, nil
}

func (c contentService) asXML(body *vo.Body, rawBytes []byte) ([]byte, vo.ContentType, error) {
	httpContentType := vo.NewContentTypeXml()

	if body.ContentType().IsJSON() {
		converted, err := c.converter.ConvertJSONToXML(rawBytes)
		if checker.NonNil(err) {
			return nil, "", errors.Inherit(err, "content-type failed: op=convert json->xml")
		}
		return converted, httpContentType, nil
	}

	converted, err := c.converter.ConvertTextToXML(rawBytes)
	if checker.NonNil(err) {
		return nil, "", errors.Inherit(err, "content-type failed: op=convert text->xml")
	}

	return converted, httpContentType, nil
}

func (c contentService) asGzip(rawBytes []byte) ([]byte, vo.ContentEncoding, error) {
	compressed, err := compressor.ToGzipWithErr(rawBytes)
	if checker.NonNil(err) {
		return nil, "", errors.Inherit(err, "content-encoding failed: op=gzip")
	}
	return compressed, vo.NewContentEncodingGzip(), nil
}

func (c contentService) asDeflate(rawBytes []byte) ([]byte, vo.ContentEncoding, error) {
	compressed, err := compressor.ToDeflateWithErr(rawBytes)
	if checker.NonNil(err) {
		return nil, "", errors.Inherit(err, "content-encoding failed: op=deflate")
	}
	return compressed, vo.NewContentEncodingDeflate(), nil
}
