package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"net/url"
)

type Query map[string][]string

func NewQuery(httpQuery url.Values) Query {
	return Query(httpQuery)
}

func (q Query) Add(key, value string) (r Query) {
	r = q.copy()
	r[key] = append(r[key], value)
	return r
}

func (q Query) Set(key, value string) (r Query) {
	r = q.copy()
	r[key] = []string{value}
	return r
}

func (q Query) Del(key string) (r Query) {
	r = q.copy()
	delete(r, key)
	return r
}

func (q Query) Get(key string) string {
	values := q[key]
	if helper.IsNotEmpty(values) {
		return values[len(values)-1]
	}
	return ""
}

func (q Query) FilterByForwarded(forwardedQueries []string) (r Query) {
	r = q.copy()
	for key := range q.copy() {
		if helper.NotContains(forwardedQueries, "*") && helper.NotContains(forwardedQueries, key) {
			r = q.Del(key)
		}
	}
	return r
}

func (q Query) copy() (r Query) {
	r = Query{}
	for key, value := range q {
		r[key] = value
	}
	return r
}
