package caching

type Cache interface {
	Set(k string, x interface{})
	Get(k string) (interface{}, bool)
}
