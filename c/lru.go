package c

import (
	grlu "github.com/hashicorp/golang-lru"
)

type lru struct {
	grlu.Cache
}

func NewCacheRLU(size int) (Cache, error) {
	c, err := grlu.New(size)
	if err != nil {
		return nil, err
	}

	return &lru{
		*c,
	}, nil
}

func (c lru) Set(k string, x interface{}) {
	c.Cache.Add(k, x)
}

func (c lru) Get(k string) (interface{}, bool) {
	return c.Cache.Get(k)
}
