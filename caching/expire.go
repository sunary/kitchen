package caching

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type expire struct {
	*cache.Cache
}

// NewCacheExpire ...
func NewCacheExpire(expireAfter time.Duration) Cache {
	return &expire{
		cache.New(expireAfter, 2*expireAfter),
	}
}

// Set ...
func (c expire) Set(k string, x interface{}) {
	c.Cache.Set(k, x, cache.DefaultExpiration)
}

// Add ...
func (c expire) Add(k string, x interface{}) error {
	return c.Cache.Add(k, x, cache.DefaultExpiration)
}

// SetForever ...
func (c *expire) SetForever(k string, x interface{}) error {
	return c.Cache.Add(k, x, cache.NoExpiration)
}

// Get ...
func (c expire) Get(k string) (interface{}, bool) {
	return c.Cache.Get(k)
}
