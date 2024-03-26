package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"net/url"
)

type Query map[string][]string

func NewQuery(httpQuery url.Values) Query {
	return Query(httpQuery)
}

func (h Query) Add(key, value string) (r Query) {
	r = h.copy()
	r[key] = append(r[key], value)
	return r
}

func (h Query) Set(key, value string) (r Query) {
	r = h.copy()
	r[key] = []string{value}
	return r
}

func (h Query) Del(key string) (r Query) {
	r = h.copy()
	delete(r, key)
	return r
}

func (h Query) Get(key string) string {
	values := h[key]
	if helper.IsNotEmpty(values) {
		return values[len(values)-1]
	}
	return ""
}

func (h Query) FilterByForwarded(forwardedQueries []string) (r Query) {
	r = h.copy()
	for key := range h.copy() {
		if helper.NotContains(forwardedQueries, "*") && helper.NotContains(forwardedQueries, key) {
			r = h.Del(key)
		}
	}
	return r
}

func (h Query) Aggregate(anotherHeader Header) {

}

func (h Query) copy() (r Query) {
	for key, value := range h {
		r[key] = value
	}
	return r
}
