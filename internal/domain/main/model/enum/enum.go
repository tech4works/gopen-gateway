/*
 * Copyright 2024 Gabriel Cataldo
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package enum

import "github.com/GabrielHCataldo/go-helper/helper"

// CacheControl represents the header value of cache control.
type CacheControl string

// ContentType represents the format of the content.
type ContentType string

type ContentEncoding string

const (
	ContentEncodingGzip ContentEncoding = "gzip"
)
const (
	CacheControlNoCache CacheControl = "no-cache"
	CacheControlNoStore CacheControl = "no-store"
)

const (
	ContentTypeJson ContentType = "JSON"
	ContentTypeXml  ContentType = "XML"
	ContentTypeYml  ContentType = "YML"
	ContentTypeText ContentType = "TEXT"
)

// ContentTypeFromString converts a string representation of a content type
// to its corresponding ContentType value. It checks if the given string
// contains the string representation of ContentTypeJson or ContentTypeText case-insensitively.
// If the string contains ContentTypeJson, it returns ContentTypeJson.
// If the string contains ContentTypeText, it returns ContentTypeText.
// Otherwise, it returns an empty string.
// This function is used to convert a string content type to the ContentType enumeration value.
func ContentTypeFromString(s string) ContentType {
	// todo: aqui podemos ter XML, YAML, form-data
	if helper.ContainsIgnoreCase(s, ContentTypeJson.String()) {
		return ContentTypeJson
	} else if helper.ContainsIgnoreCase(s, ContentTypeText.String()) {
		return ContentTypeText
	}
	return ""
}

func ContentEncodingFromString(s string) ContentEncoding {
	if helper.ContainsIgnoreCase(s, ContentEncodingGzip) {
		return ContentEncodingGzip
	}
	return ""
}

// IsEnumValid checks if the CacheControl is a valid enumeration value.
// It returns true if the CacheControl is either CacheControlNoCache or CacheControlNoStore,
// otherwise it returns false.
func (c CacheControl) IsEnumValid() bool {
	switch c {
	case CacheControlNoCache, CacheControlNoStore:
		return true
	}
	return false
}

// IsEnumValid checks if the ContentType is a valid enumeration value.
// It returns true if the ContentType is either ContentTypeText, ContentTypeJson,
// ContentTypeXml, or ContentTypeYml, otherwise it returns false.
func (c ContentType) IsEnumValid() bool {
	switch c {
	case ContentTypeText, ContentTypeJson, ContentTypeXml, ContentTypeYml:
		return true
	}
	return false
}

func (c ContentEncoding) IsEnumValid() bool {
	switch c {
	case ContentEncodingGzip:
		return true
	}
	return false
}

// String returns the string representation of the ContentType value.
// It returns "application/json" if c is ContentTypeJson, "application/xml" if c is ContentTypeXml,
// "application/x-yaml" if c is ContentTypeYml, and "text/plain" for any other value of c.
// This method is used to convert the ContentType value to its corresponding MIME type string representation.
func (c ContentType) String() string {
	switch c {
	case ContentTypeJson:
		return "application/json"
	case ContentTypeXml:
		return "application/xml"
	case ContentTypeYml:
		return "application/x-yaml"
	default:
		return "text/plain"
	}
}
