package vo

import "github.com/tech4works/checker"

type BackendHTTPRequestConfig struct {
	header  *MetadataConfig
	urlPath *URLPathConfig
	query   *QueryConfig
	body    *PayloadConfig
}

func NewBackendHTTPRequestConfig(
	header *MetadataConfig,
	param *URLPathConfig,
	query *QueryConfig,
	body *PayloadConfig,
) BackendHTTPRequestConfig {
	return BackendHTTPRequestConfig{
		header:  header,
		urlPath: param,
		query:   query,
		body:    body,
	}
}

func (b BackendHTTPRequestConfig) HasURLPath() bool {
	return checker.NonNil(b.urlPath)
}

func (b BackendHTTPRequestConfig) URLPath() *URLPathConfig {
	return b.urlPath
}

func (b BackendHTTPRequestConfig) HasQuery() bool {
	return checker.NonNil(b.query)
}

func (b BackendHTTPRequestConfig) Query() *QueryConfig {
	return b.query
}

func (b BackendHTTPRequestConfig) HasHeader() bool {
	return checker.NonNil(b.header)
}

func (b BackendHTTPRequestConfig) Header() *MetadataConfig {
	return b.header
}

func (b BackendHTTPRequestConfig) HasBody() bool {
	return checker.NonNil(b.body)
}

func (b BackendHTTPRequestConfig) Body() *PayloadConfig {
	return b.body
}

func (b BackendHTTPRequestConfig) HasDataTransforms() bool {
	return checker.IsGreaterThan(b.CountAllDataTransforms(), 0)
}

func (b BackendHTTPRequestConfig) CountAllDataTransforms() (count int) {
	if b.HasURLPath() {
		count += b.URLPath().CountDataTransforms()
	}
	if b.HasQuery() {
		count += b.Query().CountDataTransforms()
	}
	if b.HasHeader() {
		count += b.Header().CountDataTransforms()

	}
	if b.HasBody() {
		count += b.Body().CountDataTransforms()
	}
	return count
}
