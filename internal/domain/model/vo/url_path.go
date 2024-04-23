package vo

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"strings"
)

// UrlPath represents a string type that represents a URL path.
type UrlPath string

// ContainsParam checks if the UrlPath contains the specified key as a path parameter.
// It returns true if the key is found, and false otherwise.
func (u UrlPath) ContainsParam(key string) bool {
	return helper.Contains(u, u.patternParamKey(key))
}

// FillParamValue replaces all occurrences of the path parameter key with the specified value.
// It returns a new UrlPath with the modified path.
func (u UrlPath) FillParamValue(key string, value string) UrlPath {
	newPath := strings.ReplaceAll(u.String(), u.patternParamKey(key), value)
	return UrlPath(newPath)
}

// String returns the string representation of the UrlPath instance.
func (u UrlPath) String() string {
	return string(u)
}

// SetParamKey adds the specified key as a path parameter to the UrlPath instance.
// If the key is not already present in the UrlPath, it appends the key to the end,
// prefixed with a "/" character. It returns the modified UrlPath.
func (u UrlPath) SetParamKey(key string) UrlPath {
	patternParamKey := u.patternParamKey(key)
	if helper.NotContains(u, patternParamKey) {
		s := fmt.Sprintf("%s/%s", u, patternParamKey)
		return UrlPath(s)
	}
	return u
}

// RenameParamKey renames the path parameter key from oldKey to newKey in the UrlPath.
// It replaces all occurrences of the oldKey with the newKey and returns a modified UrlPath instance.
func (u UrlPath) RenameParamKey(oldKey, newKey string) UrlPath {
	s := strings.ReplaceAll(u.String(), u.patternParamKey(oldKey), u.patternParamKey(newKey))
	return UrlPath(s)
}

// DeleteParamKey removes the specified key as a path parameter from the UrlPath instance.
// It replaces all occurrences of the key prefixed with a "/" character and returns
// a modified UrlPath with the key removed.
func (u UrlPath) DeleteParamKey(key string) UrlPath {
	patternParamKey := u.patternParamKey(key)
	patternParamKeyUrl := fmt.Sprintf("/%s", patternParamKey)
	s := strings.ReplaceAll(u.String(), patternParamKeyUrl, "")
	return UrlPath(s)
}

// patternParamKey generates a string representation of the specified key as a path parameter.
func (u UrlPath) patternParamKey(key string) string {
	return fmt.Sprintf(":%s", key)
}
