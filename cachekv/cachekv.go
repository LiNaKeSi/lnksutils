package cachekv

import (
	"reflect"
	"sync"

	"github.com/philippgille/gokv"
)

func New(backend gokv.Store) *cacheKv {
	return &cacheKv{
		backend: backend,
	}
}

type cacheKv struct {
	backend gokv.Store
	cache   sync.Map
}

func (c *cacheKv) Set(k string, v interface{}) error {
	c.cache.Store(k, v)
	return c.backend.Set(k, v)
}

func (c *cacheKv) Get(k string, v interface{}) (found bool, err error) {
	value, ok := c.cache.Load(k)
	if ok {
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf(value))
		return true, nil
	}
	return c.backend.Get(k, v)
}

func (c *cacheKv) Delete(k string) error {
	return c.backend.Delete(k)
}

func (c *cacheKv) Close() error {
	return c.backend.Close()
}

func (c *cacheKv) ClearCache() {
	c.cache.Range(func(k, v interface{}) bool {
		c.cache.Delete(k)
		return true
	})
}
