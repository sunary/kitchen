package caching

import (
	"sync"

	grlu "github.com/hashicorp/golang-lru"
)

type lru struct {
	grlu.Cache
	threadSafe bool
	lock       sync.RWMutex
}

// NewCacheRLU ...
func NewCacheRLU(size int) (Cache, error) {
	c, err := grlu.New(size)
	if err != nil {
		return nil, err
	}

	return &lru{
		Cache:      *c,
		threadSafe: false,
	}, nil
}

// Set ...
func (c lru) Set(k string, x interface{}) {
	if c.threadSafe {
		c.lock.Lock()
		defer c.lock.Unlock()
	}
	c.Cache.Add(k, x)
}

// Get ...
func (c lru) Get(k string) (interface{}, bool) {
	if c.threadSafe {
		c.lock.RLock()
		defer c.lock.RUnlock()
	}
	return c.Cache.Get(k)
}
