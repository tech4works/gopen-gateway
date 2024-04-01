package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"net/url"
)

type Query map[string][]string

// NewQuery takes an HTTP query and creates a new instance of Query.
func NewQuery(httpQuery url.Values) Query {
	return Query(httpQuery)
}

// Add appends the given value to the slice associated with the given key in a new copy of the Query map.
// If the key does not exist in the original Query map, a new slice containing only the given value will be created.
// The new copy of the Query map is then returned.
func (q Query) Add(key, value string) (r Query) {
	r = q.copy()
	r[key] = append(r[key], value)
	return r
}

// Set sets the value of the given key in a new copy of the Query map.
// The value is a string that will be stored in a new single-element string slice.
// The new copy of the Query map is then returned.
func (q Query) Set(key, value string) (r Query) {
	r = q.copy()
	r[key] = []string{value}
	return r
}

// Del deletes the slice associated with the given key in a new copy of the Query map.
// If the key does not exist in the original Query map, the new copy remains unchanged.
// The new copy of the Query map is then returned.
func (q Query) Del(key string) (r Query) {
	r = q.copy()
	delete(r, key)
	return r
}

// Get retrieves the last value associated with the given key from the Query map.
// If the key does not exist or the associated value slice is empty, it returns an empty string.
func (q Query) Get(key string) string {
	values := q[key]
	if helper.IsNotEmpty(values) {
		return values[len(values)-1]
	}
	return ""
}

// FilterByForwarded filters the Query map by removing keys that are not present in the forwardedQueries slice.
// If the forwardedQueries slice does not contain "*" and does not contain the key, the key is removed from the copied Query map.
// The filtered Query map is then returned.
func (q Query) FilterByForwarded(forwardedQueries []string) (r Query) {
	r = q.copy()
	for key := range q.copy() {
		if helper.NotContains(forwardedQueries, "*") && helper.NotContains(forwardedQueries, key) {
			r = q.Del(key)
		}
	}
	return r
}

// copy creates a shallow copy of the Query map.
// It iterates over each key-value pair in the original Query map and assigns them to the new copy.
// The copied Query map is then returned.
func (q Query) copy() (r Query) {
	r = Query{}
	for key, value := range q {
		r[key] = value
	}
	return r
}
