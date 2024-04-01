package api

import (
	"bytes"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"time"
)

type Request struct {
	framework *gin.Context
	gopen     vo.GOpen
	endpoint  vo.Endpoint
	writer    *dto.Writer
}

type errorResponse struct {
	File      string    `json:"file,omitempty"`
	Line      int       `json:"line,omitempty"`
	Endpoint  string    `json:"endpoint,omitempty"`
	Message   string    `json:"message,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

func (r *Request) Context() context.Context {
	return r.framework.Request.Context()
}

func (r *Request) WithContext(ctx context.Context) {
	r.framework.Request = r.framework.Request.WithContext(ctx)
}

func (r *Request) Http() *http.Request {
	return r.framework.Request
}

func (r *Request) Header() vo.Header {
	return vo.NewHeader(r.Http().Header)
}

func (r *Request) HeaderValue(key string) string {
	return r.Header().Get(key)
}

func (r *Request) AddHeader(key string, value string) {
	r.Http().Header.Add(key, value)
}

func (r *Request) SetHeader(key string, value string) {
	r.Http().Header.Set(key, value)
}

func (r *Request) RemoteAddr() string {
	return r.framework.ClientIP()
}

func (r *Request) Method() string {
	return r.Http().Method
}

func (r *Request) Url() string {
	return r.Http().URL.String()
}

func (r *Request) Body() vo.Body {
	bytesBody, _ := io.ReadAll(r.Http().Body)
	r.Http().Body = io.NopCloser(bytes.NewBuffer(bytesBody))

	// no pior das hipóteses retornamos uma string do body
	return vo.NewBodyByContentType(r.framework.GetHeader("Content-Type"), bytesBody)
}

func (r *Request) BodyString() string {
	bytesBody, _ := io.ReadAll(r.Http().Body)
	r.Http().Body = io.NopCloser(bytes.NewBuffer(bytesBody))

	return string(bytesBody)
}

func (r *Request) Params() vo.Params {
	result := vo.Params{}
	for _, param := range r.framework.Params {
		result[param.Key] = param.Value
	}
	return result
}

func (r *Request) Query() vo.Query {
	return vo.NewQuery(r.Http().URL.Query())
}

func (r *Request) Next() {
	r.framework.Next()
}

func (r *Request) GOpen() vo.GOpen {
	return r.gopen
}

func (r *Request) Endpoint() vo.Endpoint {
	return r.endpoint
}

func (r *Request) Request() vo.Request {
	return vo.NewRequest(r.Url(), r.Method(), r.Header(), r.Params(), r.Query(), r.Body())
}

func (r *Request) Writer() dto.Writer {
	return *r.writer
}

func (r *Request) Write(responseVO vo.Response) {
	// se ja tiver abortado não fazemos nada
	if r.framework.IsAborted() {
		return
	}

	// escrevemos os headers de resposta
	r.writeHeader(responseVO.Header())

	statusCode := responseVO.StatusCode()
	body := responseVO.Body()

	// verificamos se tem valor o body
	if body.IsNotEmpty() {
		r.writeBody(statusCode, responseVO.Body())
	} else {
		r.writeStatusCode(statusCode)
	}
}

func (r *Request) WriteCacheResponse(cacheResponse vo.CacheResponse) {
	// escrevemos a resposta
	r.Write(vo.NewResponseByCache(r.endpoint, cacheResponse))
}

func (r *Request) WriteError(code int, err error) {
	// escrevemos a resposta
	r.Write(vo.NewResponseByErr(r.endpoint, code, err))
}

func (r *Request) writeHeader(header vo.Header) {
	for key := range header {
		if helper.EqualsIgnoreCase(key, "Content-Length") || helper.EqualsIgnoreCase(key, "Content-Type") ||
			helper.EqualsIgnoreCase(key, "Date") {
			continue
		}
		r.framework.Header(key, header.Get(key))
	}
}

func (r *Request) writeBody(code int, body any) {
	if r.framework.IsAborted() {
		return
	}
	// respondemos o body a partir do encode configurado
	switch r.endpoint.ResponseEncode() {
	case enum.ResponseEncodeText:
		r.framework.String(code, "%s", body)
		break
	case enum.ResponseEncodeJson:
		r.framework.JSON(code, body)
		break
	case enum.ResponseEncodeXml:
		r.framework.XML(code, body)
		break
	case enum.ResponseEncodeYaml:
		r.framework.YAML(code, body)
		break
	default:
		if helper.IsJsonType(body) {
			r.framework.JSON(code, body)
		} else {
			r.framework.String(code, "%s", body)
		}
		break
	}
}

func (r *Request) writeStatusCode(code int) {
	if r.framework.IsAborted() {
		return
	}
	r.framework.Status(code)
}

func buildErrorResponse(endpoint string, err error) *errorResponse {
	detailsErr := errors.Details(err)
	if helper.IsNil(detailsErr) {
		return nil
	}
	return &errorResponse{
		File:      detailsErr.GetFile(),
		Line:      detailsErr.GetLine(),
		Endpoint:  endpoint,
		Message:   detailsErr.GetMessage(),
		Timestamp: time.Now(),
	}
}
