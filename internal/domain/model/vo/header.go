package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/consts"
	"net/http"
	"strings"
)

type Header map[string][]string

// NewHeader creates a new Header object from an existing http.Header object.
func NewHeader(httpHeader http.Header) Header {
	return Header(httpHeader)
}

// newHeaderFailed creates a new Header object with failed status values for consts.XGOpenCache, consts.XGOpenComplete,
// and consts.XGOpenSuccess.
func newHeaderFailed() Header {
	return Header{
		consts.XGOpenCache:    {"false"},
		consts.XGOpenComplete: {"false"},
		consts.XGOpenSuccess:  {"false"},
	}
}

// newResponseHeader creates a new Header object with specific values for consts.XGOpenCache, consts.XGOpenComplete, consts.XGOpenSuccess.
// The complete parameter is used to set the value of consts.XGOpenComplete header.
// The success parameter is used to set the value of consts.XGOpenSuccess header.
// The returned Header object contains the updated values for consts.XGOpenCache, consts.XGOpenComplete, consts.XGOpenSuccess headers.
func newResponseHeader(complete, success bool) Header {
	return Header{
		consts.XGOpenCache:    {"false"},
		consts.XGOpenComplete: {helper.SimpleConvertToString(complete)},
		consts.XGOpenSuccess:  {helper.SimpleConvertToString(success)},
	}
}

// Http converts the Header object to an http.Header object.
// It returns the converted http.Header object.
func (h Header) Http() http.Header {
	return http.Header(h)
}

// AddAll accepts a key (in string format) and an array of values (in string format).
// It adds these values to the existing header under the provided key.
// The function makes a copy of the existing header before performing the operation
// to avoid mutating the original header.
func (h Header) AddAll(key string, values []string) (r Header) {
	r = h.copy()
	r[key] = append(r[key], values...)
	return r
}

// Add is a method on the Header struct.
// It accepts a key and a value, both strings, as parameters.
// It copies the existing Header, appends the provided value to the slice
// of values associated with the provided key in the copied header,
// and then returns the updated header.
func (h Header) Add(key, value string) (r Header) {
	r = h.copy()
	r[key] = append(r[key], value)
	return r
}

// Set is a method on the Header type. It takes a key and a value, both of type string, and returns a Header.
// The method makes a copy of the original Header, sets the value of the given key in the copied Header to a new
// string slice containing the provided value, and then returns the modified Header copy.
func (h Header) Set(key, value string) (r Header) {
	r = h.copy()
	r[key] = []string{value}
	return r
}

// Del removes the value associated with the given key from the Header h.
// It returns a new Header object with the key removed.
// If the key does not exist in the Header, the returned Header is identical to the original.
func (h Header) Del(key string) (r Header) {
	r = h.copy()
	delete(r, key)
	return r
}

// Get retrieves the value for a specific key from a Header. If a value exists,
// it concatenates its items with a comma separator and returns them as a string.
// If no value exists for the given key or the value is empty, it returns an empty string.
func (h Header) Get(key string) string {
	values := h[key]
	if helper.IsNotEmpty(values) {
		return strings.Join(values, ", ")
	}
	return ""
}

// FilterByForwarded filters the headers of the Header object based on the forwardedHeaders array.
// The method removes the keys from the Header object if they do not match the conditions:
// - The forwardedHeaders array does not contain the "*" wildcard.
// - The forwardedHeaders array does not contain the current key.
// - The current key is not equal to the consts.XForwardedFor and consts.XTraceId constants.
//
// It returns a new Header object with the filtered headers.
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

// Aggregate combines the headers of two Header objects.
// It takes anotherHeader as a parameter and returns a new Header object that contains the combined headers.
// The method iterates through each key-value pair in anotherHeader and adds the values to the corresponding key in the new Header object.
// It uses the AddAll method to append the values to the existing ones, creating a new slice.
// The resulting Header object is returned at the end of the method.
func (h Header) Aggregate(anotherHeader Header) (r Header) {
	r = h.copy()
	for key, values := range anotherHeader {
		r = r.AddAll(key, values)
	}
	return r
}

// copy creates a deep copy of the Header object.
// It returns a new Header object that is a copy of the original Header object.
func (h Header) copy() (r Header) {
	r = Header{}
	for key, value := range h {
		r[key] = value
	}
	return r
}
