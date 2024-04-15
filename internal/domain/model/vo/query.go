package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"net/url"
	"sort"
	"strings"
)

// Query represents an HTTP query. It is a type alias for map[string][]string,
// where the key is the parameter name and the value is a slice of parameter values.
// The Query type provides various methods to manipulate and interact with the query
// parameters, such as adding, setting, deleting, filtering, and encoding.
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

// FilterByForwarded filters the Query map by the list of forwardedQueries.
// It removes any keys from the Query map that are not in forwardedQueries,
// except the wildcard character '*' which represents all keys.
// If forwardedQueries is empty or contains only the wildcard character '*',
// the original Query map is returned without any modifications.
//
// Returns:
//
//	A new copy of the Query map with keys filtered by forwardedQueries.
//
// Note:
//
//	The original Query map is not modified.
func (q Query) FilterByForwarded(forwardedQueries []string) (r Query) {
	r = q.copy()
	for key := range q.copy() {
		if helper.IsNotEmpty(forwardedQueries) && helper.NotContains(forwardedQueries, key) {
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

// Encode encodes the values into “URL encoded” form
// ("bar=baz&foo=qux") sorted by key.
func (q Query) Encode() string {
	// se for vazio retornamos a string vazia
	if helper.IsEmpty(q) {
		return ""
	}

	// instanciamos o valor string a ser usado para adicionar os valores
	var buf strings.Builder

	// obtemos as keys
	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
	}
	// fazemos o sort
	sort.Strings(keys)

	// iteramos as chaves ordenadas
	for _, k := range keys {
		// obtemos o valor da chave
		vs := q[k]
		// fazemos o sort dos valores
		sort.Strings(vs)

		// escapamos a chave
		keyEscaped := url.QueryEscape(k)

		// iteramos sobre os valores pela chave ja ordenados
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
		}
	}
	// retornamos o valor da query como string
	return buf.String()
}
