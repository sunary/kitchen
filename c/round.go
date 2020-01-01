package c

import (
	"strconv"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

type round struct {
	*cache.Cache
	seconds int64
}

func NewCacheRound(roundDuration time.Duration) Cache {
	return &round{
		cache.New(roundDuration, 2*roundDuration),
		int64(roundDuration.Seconds()),
	}
}

func (c round) Set(k string, x interface{}) {
	c.Cache.Set(c.newKey(k), x, cache.DefaultExpiration)
}

func (c round) Add(k string, x interface{}) error {
	return c.Cache.Add(c.newKey(k), x, cache.DefaultExpiration)
}

func (c round) Get(k string) (interface{}, bool) {
	return c.Cache.Get(c.newKey(k))
}
func (c round) newKey(k string) string {
	var sb strings.Builder
	sb.WriteString(k)
	sb.WriteString("-")
	sb.Write(strconv.AppendInt(nil, time.Now().Unix()/c.seconds*c.seconds, 10))
	return sb.String()
}
