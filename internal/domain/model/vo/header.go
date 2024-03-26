package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"net/http"
)

type Header map[string][]string

func NewHeader(httpHeader http.Header) Header {
	return Header(httpHeader)
}

func newResponseHeader(complete, success bool) Header {
	return Header{
		consts.XGatewayCache:    {"false"},
		consts.XGatewayComplete: {helper.SimpleConvertToString(complete)},
		consts.XGatewaySuccess:  {helper.SimpleConvertToString(success)},
	}
}

func newResponseHeaderFailed() Header {
	return Header{
		consts.XGatewayCache:    {"false"},
		consts.XGatewayComplete: {"false"},
		consts.XGatewaySuccess:  {"false"},
	}
}

func (h Header) Http() http.Header {
	return http.Header(h)
}

func (h Header) AddAll(key string, values []string) (r Header) {
	r = h.copy()
	r[key] = append(r[key], values...)
	return r
}

func (h Header) Add(key, value string) (r Header) {
	r = h.copy()
	r[key] = append(r[key], value)
	return r
}

func (h Header) Set(key, value string) (r Header) {
	r = h.copy()
	r[key] = []string{value}
	return r
}

func (h Header) Del(key string) (r Header) {
	r = h.copy()
	delete(r, key)
	return r
}

func (h Header) Get(key string) string {
	values := h[key]
	if helper.IsNotEmpty(values) {
		return values[len(values)-1]
	}
	return ""
}

func (h Header) FilterByForwarded(forwardedHeaders []string) (r Header) {
	r = h.copy()
	for key := range h.copy() {
		if helper.NotContains(forwardedHeaders, "*") && helper.NotContains(forwardedHeaders, key) &&
			helper.IsNotEqualTo(key, consts.XForwardedFor) && helper.IsNotEqualTo(key, consts.XTraceId) {
			r = h.Del(key)
		}
	}
	return r
}

func (h Header) Aggregate(anotherHeader Header) (r Header) {
	r = h.copy()
	for key, values := range anotherHeader {
		r = r.AddAll(key, values)
	}
	return r
}

func (h Header) copy() (r Header) {
	for key, value := range h {
		r[key] = value
	}
	return r
}
