package vo

type CacheInfo struct {
	hit          bool
	ttl          Duration
	remainingTTL Duration
}

func NewEmptyCacheInfo() CacheInfo {
	return CacheInfo{}
}

func NewCacheInfo(hit bool, ttl, remainingTTL Duration) CacheInfo {
	return CacheInfo{
		hit:          hit,
		ttl:          ttl,
		remainingTTL: remainingTTL,
	}
}

func (c CacheInfo) Hit() bool {
	return c.hit
}

func (c CacheInfo) RemainingTTL() Duration {
	return c.remainingTTL
}

func (c CacheInfo) HasRemainingTTL() bool {
	return c.remainingTTL > 0
}

func (c CacheInfo) Map() any {
	return map[string]any{
		"hit":           c.hit,
		"ttl":           c.ttl.Time().Milliseconds(),
		"remaining-ttl": c.remainingTTL.Time().Milliseconds(),
	}
}
