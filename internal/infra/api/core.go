package api

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/vo"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"net/http"
	"sync"
)

type Context struct {
	mutex     *sync.RWMutex
	framework *gin.Context
	gopen     vo.Gopen
	endpoint  vo.Endpoint
	request   vo.Request
	response  vo.Response
}

// Context returns the context of the Context. It delegates the call to the underlying framework's Context.Context() method.
func (r *Context) Context() context.Context {
	return r.framework.Request.Context()
}

// Gopen returns the Gopen object associated with the Context. It retrieves the Gopen value from the Context object.
func (r *Context) Gopen() vo.Gopen {
	return r.gopen
}

// Endpoint returns the endpoint associated with the request.
// It retrieves the endpoint value from the `endpoint` field of the Context struct.
func (r *Context) Endpoint() vo.Endpoint {
	return r.endpoint
}

// Request returns the request object of the Context.
// It returns the `request` field of the Context struct.
func (r *Context) Request() vo.Request {
	return r.request
}

// Response returns the response of the Context. It returns the response object stored
// in the Context struct.
func (r *Context) Response() vo.Response {
	return r.response
}

// Http returns the underlying HTTP request object of the Context.
// It delegates the call to the underlying framework's Request property.
func (r *Context) Http() *http.Request {
	return r.framework.Request
}

// SetRequestContext sets the context of the Context to the provided context.
// It updates the underlying framework's Context.Context() method to use the new context.
func (r *Context) SetRequestContext(ctx context.Context) {
	r.framework.Request = r.framework.Request.WithContext(ctx)
}

// Header returns the `vo.Header` of the `Request`. It creates a new `vo.Header` using the underlying `http.Header`
// from the `Request`.
func (r *Context) Header() vo.Header {
	return r.request.Header()
}

// HeaderValue returns the value of the specified header key. It delegates the call to the underlying Context's
// Header().Get method.
func (r *Context) HeaderValue(key string) string {
	return r.Header().Get(key)
}

// AddHeader adds a new header to the HTTP request.
// It takes a key and value as parameters and adds them to the request's headers.
// Example usage:
//
//	req.AddHeader("Content-Type", "application/json")
//	req.AddHeader("Authorization", "Bearer token123")
//
// The method first creates a new header using the provided key and value.
// It then adds the header to the request using the Header() method of the context.
// Finally, it sets the updated request header using the SetHeader() method of the request.
func (r *Context) AddHeader(key, value string) {
	header := r.Header().Add(key, value)
	r.request = r.Request().SetHeader(header)
}

// SetHeader sets the value of the specified header key for the Context object.
// It delegates the call to the underlying framework's Request.Header.Set() method.
// Example usage: req.SetHeader("X-Forwarded-For", req.RemoteAddr()) and req.SetHeader("X-TraceId", t.traceProvider.GenerateTraceId())
// The SetHeader method takes a key and value as parameters, set the key value pair in the Context object's header.
// It uses the underlying framework's Request.Header.Set() method to update the header value.
// It returns nothing.
func (r *Context) SetHeader(key, value string) {
	header := r.Header().Set(key, value)
	r.request = r.Request().SetHeader(header)
}

// RemoteAddr returns the client's remote network address in the format "IP:port". It delegates the call to the
// underlying framework's ClientIP() method.
func (r *Context) RemoteAddr() string {
	return r.framework.ClientIP()
}

// Method returns the HTTP method of the Context.
// It delegates the call to the underlying framework's Request.Method() method.
func (r *Context) Method() string {
	return r.Request().Method()
}

// Url returns the URL of the request.
// It delegates the call to the underlying framework's Request.Url() method.
func (r *Context) Url() string {
	return r.Request().Url()
}

// Uri returns the URI of the Context. It delegates the call to the
// underlying Request's Uri() method.
func (r *Context) Uri() string {
	return r.Request().Uri()
}

// Body returns the body of the Context. It delegates the call to the
// underlying framework's Request.Body() method.
func (r *Context) Body() vo.Body {
	return r.Request().Body()
}

// BodyString returns the string representation of the body. It delegates the call to the
// underlying Body() method and then calls String() on the returned value.
func (r *Context) BodyString() string {
	body := r.Body()
	return body.String()
}

// Params returns the params of the Context.
// It delegates the call to the underlying Request's Params() method.
func (r *Context) Params() vo.Params {
	return r.Request().Params()
}

// Query returns the query object associated with the current context. It retrieves the query object
// by delegating the call to the underlying framework's Request().Query() method.
func (r *Context) Query() vo.Query {
	return r.Request().Query()
}

// Next calls the underlying framework's Next method to proceed to the next handler in the request chain.
func (r *Context) Next() {
	r.framework.Next()
}

// Write writes the response to the client.
// It first checks if the request has already been aborted, in which case it does nothing.
// Then, it writes the response headers.
// It retrieves the status code and body from the responseVO.
// If the body is not empty, it writes the body along with the status code.
// Otherwise, it only writes the status code.
func (r *Context) Write(responseVO vo.Response) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// se ja tiver abortado não fazemos nada
	if r.framework.IsAborted() {
		return
	}

	// escrevemos os headers de resposta
	r.writeHeader(responseVO.Header())

	// instanciamos os valores a serem utilizados
	statusCode := responseVO.StatusCode()
	body := responseVO.Body()

	// verificamos se tem valor o body
	if body.IsNotEmpty() {
		r.writeBody(statusCode, body.ToWrite())
	} else {
		r.writeStatusCode(statusCode)
	}

	// abortamos a requisição
	r.framework.Abort()

	// setamos a resposta VO escrita
	r.response = responseVO
}

// WriteCacheResponse writes the cache response to the client's response.
// It creates a new response using the cache response and writes it.
func (r *Context) WriteCacheResponse(cacheResponse vo.CacheResponse) {
	// preparamos a resposta
	responseVO := vo.NewResponseByCache(r.endpoint, cacheResponse)
	// escrevemos a resposta
	r.Write(responseVO)
}

// WriteError writes an error response to the client.
// It creates a new Response object with the provided code and error, and delegates the writing of the response to the
// Write method.
func (r *Context) WriteError(code int, err error) {
	// escrevemos a resposta
	r.Write(vo.NewResponseByErr(r.endpoint, code, err))
}

// writeHeader sets the headers in the HTTP response as received in the `header` argument, excluding certain headers.
// Headers to be ignored: "Content-Length", "Content-Type", "Date".
// The method delegates the actual header setting to the underlying framework's `Header()` method.
// Example usage can be found in the `Write()` method.
func (r *Context) writeHeader(header vo.Header) {
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
func (r *Context) writeBody(code int, body any) {
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
func (r *Context) writeStatusCode(code int) {
	if r.framework.IsAborted() {
		return
	}
	r.framework.Status(code)
}
