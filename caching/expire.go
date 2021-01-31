package caching

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

type expire struct {
	*cache.Cache
	threadSafe bool
	lock       sync.RWMutex
}

// NewCacheExpire ...
func NewCacheExpire(expireAfter time.Duration) Cache {
	return &expire{
		Cache:      cache.New(expireAfter, 2*expireAfter),
		threadSafe: false,
	}
}

// Set ...
func (c expire) Set(k string, x interface{}) {
	if c.threadSafe {
		c.lock.Lock()
		defer c.lock.Unlock()
	}
	c.Cache.Set(k, x, cache.DefaultExpiration)
}

// Add ...
func (c expire) Add(k string, x interface{}) error {
	if c.threadSafe {
		c.lock.Lock()
		defer c.lock.Unlock()
	}
	return c.Cache.Add(k, x, cache.DefaultExpiration)
}

// SetForever ...
func (c *expire) SetForever(k string, x interface{}) error {
	if c.threadSafe {
		c.lock.Lock()
		defer c.lock.Unlock()
	}
	return c.Cache.Add(k, x, cache.NoExpiration)
}

// Get ...
func (c expire) Get(k string) (interface{}, bool) {
	if c.threadSafe {
		c.lock.RLock()
		defer c.lock.RUnlock()
	}
	return c.Cache.Get(k)
}
