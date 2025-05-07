package vo

import "github.com/tech4works/gopen-gateway/internal/domain/model/enum"

type Proxy struct {
	provider enum.ProxyProvider
	token    string
	domain   []string
}

func NewProxy(provider enum.ProxyProvider, token string, domains []string) *Proxy {
	return &Proxy{
		provider: provider,
		token:    token,
		domain:   domains,
	}
}

func (p Proxy) Provider() enum.ProxyProvider {
	return p.provider
}

func (p Proxy) Token() string {
	return p.token
}

func (p Proxy) Domains() []string {
	return p.domain
}
