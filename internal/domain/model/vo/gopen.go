package vo

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/gopen-gateway/internal/domain/model/enum"
	"github.com/chenyahui/gin-cache/persist"
	"net/http"
	"time"
)

type Gopen struct {
	middlewares map[string]Backend
}

type Cache struct {
	Duration          time.Duration
	StrategyKeys      []string
	AllowCacheControl bool
	MemoryStore       persist.CacheStore
}

type SecurityCors struct {
	AllowCountries []string
	AllowOrigins   []string
	AllowMethods   []string
	AllowHeaders   []string
}

type Endpoint struct {
	aggregateResponses bool
	abortIfStatusCodes []int
	beforeware         []string
	afterware          []string
	backends           []Backend
}

type Modifier struct {
	Context enum.ModifierContext
	Scope   enum.ModifierScope
	Action  enum.ModifierAction
	Key     string
	Value   string
}

func (e Endpoint) Beforeware() []string {
	return e.beforeware
}

func (e Endpoint) CountBackends() int {
	return len(e.beforeware) + len(e.backends) + len(e.afterware)
}

func (e Endpoint) Completed(responseHistorySize int) bool {
	return helper.Equals(responseHistorySize, e.CountBackends())
}

func (e Endpoint) AbortSequencial(responseVO Response) bool {
	if helper.IsEmpty(e.abortIfStatusCodes) {
		return helper.IsGreaterThan(responseVO.statusCode, http.StatusBadRequest)
	}
	return helper.Contains(e.abortIfStatusCodes, responseVO.statusCode)
}

func (m Modifier) Valid() bool {
	return helper.IsNotEmpty(m) && helper.IsNotEmpty(m.Value)
}
