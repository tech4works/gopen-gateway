package vo

type ProxyConfig struct {
	token  string
	domain []string
}

func NewProxyConfig(token string, domains []string) *ProxyConfig {
	return &ProxyConfig{
		token:  token,
		domain: domains,
	}
}

func (p ProxyConfig) Token() string {
	return p.token
}

func (p ProxyConfig) Domains() []string {
	return p.domain
}
