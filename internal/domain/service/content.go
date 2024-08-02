package service

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/tech4works/gopen-gateway/internal/domain"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"strconv"
)

type contentService struct {
	converter domain.Converter
}

type Content interface {
	ModifyBodyContentType(body *vo.Body, contentType enum.ContentType) (*vo.Body, error)
	ModifyBodyContentEncoding(body *vo.Body, contentEncoding enum.ContentEncoding) (*vo.Body, error)
}

func NewContent(converter domain.Converter) Content {
	return contentService{
		converter: converter,
	}
}

func (c contentService) ModifyBodyContentType(body *vo.Body, contentType enum.ContentType) (*vo.Body, error) {
	if !contentType.IsEnumValid() ||
		body.ContentType().IsUnknown() ||
		helper.Equals(body.ContentType().ToEnum(), contentType) {
		return body, nil
	}

	rawBytes, err := body.Bytes()
	if helper.IsNotNil(err) {
		return body, err
	}

	// todo: separar em metodos auxiliares
	var bodyBytes []byte
	var httpContentType vo.ContentType

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
	if helper.IsNotNil(err) {
		return body, err
	}

	buffer, err := helper.ConvertToBuffer(bodyBytes)
	if helper.IsNotNil(err) {
		return body, err
	}

	return vo.NewBodyWithContentType(httpContentType, buffer), nil
}

func (c contentService) ModifyBodyContentEncoding(body *vo.Body, contentEncoding enum.ContentEncoding) (*vo.Body,
	error) {
	if !contentEncoding.IsEnumValid() || helper.EqualsIgnoreCase(body.ContentEncoding(), contentEncoding) {
		return body, nil
	}

	rawBytes, err := body.Bytes()
	if helper.IsNotNil(err) {
		return body, err
	}

	// todo: separar em metodos auxiliares
	var bodyBytes []byte
	var httpContentEncoding vo.ContentEncoding

	switch contentEncoding {
	case enum.ContentEncodingGzip:
		httpContentEncoding = vo.NewContentEncodingGzip()
		bodyBytes, err = helper.CompressWithGzip(rawBytes)
	case enum.ContentEncodingDeflate:
		httpContentEncoding = vo.NewContentEncodingDeflate()
		bodyBytes, err = helper.CompressWithDeflate(rawBytes)
	default:
		bodyBytes = rawBytes
	}
	if helper.IsNotNil(err) {
		return body, err
	}

	buffer, err := helper.ConvertToBuffer(bodyBytes)
	if helper.IsNotNil(err) {
		return body, err
	}

	return vo.NewBody(body.ContentType().String(), httpContentEncoding.String(), buffer), nil
}
