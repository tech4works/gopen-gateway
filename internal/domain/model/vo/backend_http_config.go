package vo

type BackendHTTPConfig struct {
	hosts   []string
	path    string
	method  string
	request BackendHTTPRequestConfig
}

func NewBackendHTTPConfig(
	hosts []string,
	path,
	method string,
	request BackendHTTPRequestConfig,
) *BackendHTTPConfig {
	return &BackendHTTPConfig{
		hosts:   hosts,
		path:    path,
		method:  method,
		request: request,
	}
}

func (b *BackendHTTPConfig) Hosts() []string {
	return b.hosts
}

func (b *BackendHTTPConfig) Path() string {
	return b.path
}

func (b *BackendHTTPConfig) Method() string {
	return b.method
}

func (b *BackendHTTPConfig) Request() BackendHTTPRequestConfig {
	return b.request
}

func (b *BackendHTTPConfig) CountAllDataTransforms() (count int) {
	return b.Request().CountAllDataTransforms()
}
