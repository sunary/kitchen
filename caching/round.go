package caching

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

type round struct {
	*cache.Cache
	threadSafe bool
	lock       sync.RWMutex
	seconds    int64
}

// NewCacheRound ...
func NewCacheRound(roundDuration time.Duration) Cache {
	return &round{
		Cache:      cache.New(roundDuration, 2*roundDuration),
		threadSafe: false,
		seconds:    int64(roundDuration.Seconds()),
	}
}

// Set ...
func (c round) Set(k string, x interface{}) {
	if c.threadSafe {
		c.lock.Lock()
		defer c.lock.Unlock()
	}
	c.Cache.Set(c.newKey(k), x, cache.DefaultExpiration)
}

// Add ...
func (c round) Add(k string, x interface{}) error {
	if c.threadSafe {
		c.lock.Lock()
		defer c.lock.Unlock()
	}
	return c.Cache.Add(c.newKey(k), x, cache.DefaultExpiration)
}

// Get ...
func (c round) Get(k string) (interface{}, bool) {
	if c.threadSafe {
		c.lock.RLock()
		defer c.lock.RUnlock()
	}
	return c.Cache.Get(c.newKey(k))
}

func (c round) newKey(k string) string {
	var sb strings.Builder
	sb.WriteString(k)
	sb.WriteString("-")
	sb.Write(strconv.AppendInt(nil, time.Now().Unix()/c.seconds, 10))
	return sb.String()
}
