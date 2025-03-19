package cachekv

import (
	"sync"

	"github.com/linakesi/lnksutils/fskv"
)

func New(backend fskv.Store) *cacheKv {
	return &cacheKv{
		backend: backend,
	}
}

type cacheKv struct {
	backend fskv.Store
	cache   sync.Map
}

func (c *cacheKv) Set(k string, v interface{}) error {
	data, err := c.backend.Codec.Marshal(v)
	if err != nil {
		return err
	}
	c.cache.Store(k, data)
	return c.backend.Set(k, v)
}

func (c *cacheKv) Get(k string, v interface{}) (found bool, err error) {
	value, ok := c.cache.Load(k)
	if ok {
		data := value.([]byte)
		err = c.backend.Codec.Unmarshal(data, v)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return c.backend.Get(k, v)
}

func (c *cacheKv) Delete(k string) error {
	c.cache.Delete(k)
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
