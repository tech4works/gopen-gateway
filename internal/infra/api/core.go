package api

import (
	"bytes"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/app/model/dto"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"io"
	"net/http"
)

type Request struct {
	framework *gin.Context
	gopen     vo.Gopen
	endpoint  vo.Endpoint
	writer    *dto.Writer
}

// Context returns the context of the Request. It delegates the call to the underlying framework's Request.Context() method.
func (r *Request) Context() context.Context {
	return r.framework.Request.Context()
}

// WithContext sets the context of the Request to the provided context.
// It updates the underlying framework's Request.Context() method to use the new context.
func (r *Request) WithContext(ctx context.Context) {
	r.framework.Request = r.framework.Request.WithContext(ctx)
}

// Http returns the underlying *http.Request object of the Request.
// It simply returns the framework's Request object.
func (r *Request) Http() *http.Request {
	return r.framework.Request
}

// Header returns the `vo.Header` of the `Request`. It creates a new `vo.Header` using the underlying `http.Header`
// from the `Request`.
func (r *Request) Header() vo.Header {
	return vo.NewHeader(r.Http().Header)
}

// HeaderValue returns the value of the specified header key. It delegates the call to the underlying Request's
// Header().Get method.
func (r *Request) HeaderValue(key string) string {
	return r.Header().Get(key)
}

// AddHeader adds a new header to the HTTP request.
// It takes a key and value as parameters and adds them to the request's headers.
// Example usage:
//
//	req.AddHeader("Content-Type", "application/json")
//	req.AddHeader("Authorization", "Bearer token123")
func (r *Request) AddHeader(key, value string) {
	r.Http().Header.Add(key, value)
}

// SetHeader sets the value of the specified header key for the Request object.
// It delegates the call to the underlying framework's Request.Header.Set() method.
// Example usage:
//
//	req.SetHeader("X-Forwarded-For", req.RemoteAddr())
//	req.SetHeader("X-TraceId", t.traceProvider.GenerateTraceId())
func (r *Request) SetHeader(key, value string) {
	r.Http().Header.Set(key, value)
}

// RemoteAddr returns the client's remote network address in the format "IP:port". It delegates the call to the
// underlying framework's ClientIP() method.
func (r *Request) RemoteAddr() string {
	return r.framework.ClientIP()
}

// Method returns the HTTP method of the Request.
// It retrieves the method from the underlying HTTP request.
func (r *Request) Method() string {
	return r.Http().Method
}

// Url returns the URL of the request. It delegates the call to the underlying framework's Request.URL.String() method.
func (r *Request) Url() string {
	return r.Http().URL.String()
}

// Body reads the request body and returns a vo.Body object representing the body content. It also updates the
// underlying request's Body with a new io.ReadCloser to ensure that the body
func (r *Request) Body() vo.Body {
	bytesBody, _ := io.ReadAll(r.Http().Body)
	r.Http().Body = io.NopCloser(bytes.NewBuffer(bytesBody))

	// no pior das hipóteses retornamos uma string do body
	return vo.NewBodyByContentType(r.framework.GetHeader("Content-Type"), bytesBody)
}

// BodyString returns the body of the Request as a string.
// It reads the bytes from the underlying framework's Request.Body and converts them to a string.
// The original Request.Body is replaced with a new io.ReadCloser that reads from a buffer containing the bytes.
func (r *Request) BodyString() string {
	bytesBody, _ := io.ReadAll(r.Http().Body)
	r.Http().Body = io.NopCloser(bytes.NewBuffer(bytesBody))

	return string(bytesBody)
}

// Params returns a copy of the parameters stored in the framework's request.
// It converts the parameters into a vo.Params map, where the key is the parameter's key and the value is the parameter's
// value. The method iterates over the framework's Params slice and adds each parameter to the result map.
// Returns an empty vo.Params map if there are no parameters.
func (r *Request) Params() vo.Params {
	result := vo.Params{}
	for _, param := range r.framework.Params {
		result[param.Key] = param.Value
	}
	return result
}

// Query returns a vo.Query object representing the query parameters of the request's URL.
// It creates a new vo.Query object using the URL query parameters obtained from the HTTP request.
func (r *Request) Query() vo.Query {
	return vo.NewQuery(r.Http().URL.Query())
}

// Next calls the underlying framework's Next method to proceed to the next handler in the request chain.
func (r *Request) Next() {
	r.framework.Next()
}

// Gopen returns the Gopen object associated with the Request. It retrieves the Gopen value from the Request object.
func (r *Request) Gopen() vo.Gopen {
	return r.gopen
}

// Endpoint returns the endpoint associated with the request.
// It retrieves the endpoint value from the `endpoint` field of the Request struct.
func (r *Request) Endpoint() vo.Endpoint {
	return r.endpoint
}

// Request returns a new `vo.Request` object based on the current `Request` instance.
// It creates a new `vo.Request` object using the URL, method, header, params, query, and body of the Request struct
func (r *Request) Request() vo.Request {
	return vo.NewRequest(r.Url(), r.Method(), r.Header(), r.Params(), r.Query(), r.Body())
}

// Writer returns the Writer object associated with the Request. It allows writing response data to the client.
// The Writer object is obtained from the underlying `dto.Writer` field of the Request.
func (r *Request) Writer() dto.Writer {
	return *r.writer
}

// Write writes the response to the client.
// It first checks if the request has already been aborted, in which case it does nothing.
// Then, it writes the response headers.
// It retrieves the status code and body from the responseVO.
// If the body is not empty, it writes the body along with the status code.
// Otherwise, it only writes the status code.
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
		r.writeBody(statusCode, body.ToWrite())
	} else {
		r.writeStatusCode(statusCode)
	}

	// abortamos
	r.framework.Abort()
}

// WriteCacheResponse writes the cache response to the client's response.
// It creates a new response using the cache response and writes it.
func (r *Request) WriteCacheResponse(cacheResponse vo.CacheResponse) {
	// preparamos a resposta
	responseVO := vo.NewResponseByCache(r.endpoint, cacheResponse)
	// escrevemos a resposta
	r.Write(responseVO)
}

// WriteError writes an error response to the client.
// It creates a new Response object with the provided code and error, and delegates the writing of the response to the
// Write method.
func (r *Request) WriteError(code int, err error) {
	// escrevemos a resposta
	r.Write(vo.NewResponseByErr(r.endpoint, code, err))
}

// writeHeader sets the headers in the HTTP response as received in the `header` argument, excluding certain headers.
// Headers to be ignored: "Content-Length", "Content-Type", "Date".
// The method delegates the actual header setting to the underlying framework's `Header()` method.
// Example usage can be found in the `Write()` method.
func (r *Request) writeHeader(header vo.Header) {
	for key := range header {
		if helper.EqualsIgnoreCase(key, "Content-Length") || helper.EqualsIgnoreCase(key, "Content-Type") ||
			helper.EqualsIgnoreCase(key, "Date") {
			continue
		}
		r.framework.Header(key, header.Get(key))
	}
}

// writeBody writes the response body based on the configured response encoding of the endpoint.
// If the framework is aborted, the method returns early without writing the body.
// If the response encoding is set to ResponseEncodeText, the body is written as a string using the given code.
// If the response encoding is set to ResponseEncodeJson, the body is written as JSON using the given code.
// If the response encoding is set to ResponseEncodeXml, the body is written as XML using the given code.
// If the response encoding is set to ResponseEncodeYaml, the body is written as YAML using the given code.
// If none of the above cases match and the body is of JSON type, it is written as JSON using the given code.
// If none of the above cases match and the body is not of JSON type, it is written as a string using the given code.
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

// writeStatusCode writes the HTTP status code to the response.
// If the request is already aborted, it does nothing.
// It sets the status code in the underlying framework using the given code.
// Parameter:
//   - code: the HTTP status code to be set in the response.
func (r *Request) writeStatusCode(code int) {
	if r.framework.IsAborted() {
		return
	}
	r.framework.Status(code)
}
