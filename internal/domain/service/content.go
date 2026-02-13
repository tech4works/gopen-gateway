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
		return body, err
	}

	buffer, err := converter.ToBufferWithErr(bodyBytes)
	if checker.NonNil(err) {
		return body, err
	}

	return vo.NewBodyWithContentType(httpContentType, buffer), nil
}

func (c contentService) ModifyBodyContentEncoding(contentEncoding enum.ContentEncoding, body *vo.Body) (*vo.Body, error) {
	if !contentEncoding.IsEnumValid() || checker.EqualsIgnoreCase(body.ContentEncoding(), contentEncoding) {
		return body, nil
	}

	bodyBytes, httpContentEncoding, err := c.modifyBodyContentEncoding(body, contentEncoding)
	if checker.NonNil(err) {
		return body, err
	}

	buffer, err := converter.ToBufferWithErr(bodyBytes)
	if checker.NonNil(err) {
		return body, err
	}

	return vo.NewBody(body.ContentType().String(), httpContentEncoding.String(), buffer), nil
}

func (c contentService) modifyBodyContentType(body *vo.Body, contentType enum.ContentType) (
	bodyBytes []byte, httpContentType vo.ContentType, err error) {
	rawBytes, err := body.Bytes()
	if checker.NonNil(err) {
		return
	}

	switch contentType {
	case enum.ContentTypePlainText:
		httpContentType = vo.NewContentTypeTextPlain()
		bodyBytes = []byte(strconv.Quote(string(rawBytes)))
	case enum.ContentTypeJson:
		httpContentType = vo.NewContentTypeJson()
		if body.ContentType().IsXML() {
			bodyBytes, err = c.converter.ConvertXMLToJSON(rawBytes)
		} else {
			bodyBytes, err = c.converter.ConvertTextToJSON(rawBytes)
		}
	case enum.ContentTypeXml:
		httpContentType = vo.NewContentTypeXml()
		if body.ContentType().IsJSON() {
			bodyBytes, err = c.converter.ConvertJSONToXML(rawBytes)
		} else {
			bodyBytes, err = c.converter.ConvertTextToXML(rawBytes)
		}
	}

	return
}

func (c contentService) modifyBodyContentEncoding(body *vo.Body, contentEncoding enum.ContentEncoding) (
	bodyBytes []byte, httpContentEncoding vo.ContentEncoding, err error) {
	rawBytes, err := body.Bytes()
	if checker.NonNil(err) {
		return
	}

	switch contentEncoding {
	case enum.ContentEncodingGzip:
		httpContentEncoding = vo.NewContentEncodingGzip()
		bodyBytes, err = compressor.ToGzipWithErr(rawBytes)
	case enum.ContentEncodingDeflate:
		httpContentEncoding = vo.NewContentEncodingDeflate()
		bodyBytes, err = compressor.ToDeflateWithErr(rawBytes)
	default:
		bodyBytes = rawBytes
	}

	return
}
