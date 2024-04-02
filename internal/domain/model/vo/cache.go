package vo

import (
	"fmt"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"net/http"
	"strings"
	"time"
)

type Cache struct {
	duration          time.Duration
	strategyHeaders   []string
	allowCacheControl *bool
}

// Duration returns the value of the duration field in the Cache struct.
func (c Cache) Duration() time.Duration {
	return c.duration
}

// Enabled returns true if the cache is enabled, false otherwise.
// The cache is considered enabled if the duration field in the Cache struct is greater than 0.
func (c Cache) Enabled() bool {
	return helper.IsGreaterThan(c.duration, 0)
}

// Disabled returns the opposite of the Enabled method. It indicates if the cache is disabled or not.
func (c Cache) Disabled() bool {
	return !c.Enabled()
}

// CanRead checks if the cache is active and if the Cache-Control header value is not "no-cache" and the HTTP method is "GET".
// If the cache is disabled, it returns false.
// It returns true if the cache is active and the conditions are met, otherwise false.
func (c Cache) CanRead(httpMethod string, header Header) bool {
	// verificamos se ta ativo
	if c.Disabled() {
		return false
	}

	// obtemos o cache control enum do ctx de requisição
	cacheControl := c.CacheControlEnum(header)

	// verificamos se no Cache-Control enviado veio como "no-cache" e se o método da requisição é GET
	return helper.IsNotEqualTo(enum.CacheControlNoCache, cacheControl) && helper.Equals(httpMethod, http.MethodGet)
}

// CanWrite checks if it is possible to write to the cache based on the HTTP method and header.
// If the cache is disabled, it returns false.
// It retrieves the cache control enum from the request header.
// It returns false if the Cache-Control header contains "no-store" and the HTTP method is GET, otherwise it returns true.
func (c Cache) CanWrite(httpMethod string, header Header) bool {
	// verificamos se ta ativo
	if c.Disabled() {
		return false
	}

	// obtemos o cache control enum do ctx de requisição
	cacheControl := c.CacheControlEnum(header)

	// verificamos se no Cache-Control enviado veio como "no-store" e se o método da requisição é GET
	return helper.IsNotEqualTo(enum.CacheControlNoStore, cacheControl) && helper.Equals(httpMethod, http.MethodGet)
}

// CacheControlEnum takes a Header and returns the CacheControl enum value.
// If caching is allowed and the Cache-Control header is present, it tries to parse the value and convert it
func (c Cache) CacheControlEnum(header Header) (cacheControl enum.CacheControl) {
	// caso esteja permitido o cache control obtemos do header
	if helper.IsNotNil(c.allowCacheControl) && *c.allowCacheControl {
		cacheControl = enum.CacheControl(header.Get("Cache-Control"))
	}
	return cacheControl
}

// StrategyKey generates a key for caching based on the HTTP method, URL, and header values.
// The key is initially constructed with the request method and URL.
// Then, the method iterates through the strategyHeaders field of the Cache object
// to collect the corresponding values from the header parameter.
// If the values are found, they are separated with a colon delimiter.
// If the strategyKey is not empty, it is appended to the key string.
// The final key is returned.
func (c Cache) StrategyKey(httpMethod string, httpUrl string, header Header) string {
	// construímos a chave inicialmente com os valores de requisição
	key := fmt.Sprintf("%s:%s", httpMethod, httpUrl)

	var strategyValues []string
	// iteramos as chaves para obter os valores
	for _, strategyKey := range c.strategyHeaders {
		valueByStrategyKey := header.Get(strategyKey)
		if helper.IsNotEmpty(valueByStrategyKey) {
			strategyValues = append(strategyValues, valueByStrategyKey)
		}
	}
	// caso tenha encontrado valores, separamos os mesmos
	strategyKey := strings.Join(strategyValues, ":")

	// caso o valor não esteja vazio retornamos o key padrão com a estratégia imposto no objeto de valor
	if helper.IsNotEmpty(strategyKey) {
		key = fmt.Sprintf("%s:%s", key, strategyKey)
	}
	return key
}
