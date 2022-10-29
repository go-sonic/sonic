package cache

import (
	"time"

	goCache "github.com/patrickmn/go-cache"
)

type Cache interface {
	SetDefault(key string, value interface{})
	Set(key string, value interface{}, expiration time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
	BatchDelete(keys []string)
}

var _ Cache = &cacheImpl{}

type cacheImpl struct {
	goCache *goCache.Cache
}

func NewCache() Cache {
	return &cacheImpl{
		goCache: goCache.New(time.Hour, time.Hour),
	}
}

// SetDefault  to cache with defaultExpiration time
func (c *cacheImpl) SetDefault(key string, value interface{}) {
	c.goCache.SetDefault(key, value)
}

// Set to cache with expiration in params
func (c *cacheImpl) Set(key string, value interface{}, expiration time.Duration) {
	c.goCache.Set(key, value, expiration)
}

// Get key's value
func (c *cacheImpl) Get(key string) (interface{}, bool) {
	return c.goCache.Get(key)
}

func (c *cacheImpl) Delete(key string) {
	c.goCache.Delete(key)
}

func (c *cacheImpl) BatchDelete(keys []string) {
	for _, key := range keys {
		c.Delete(key)
	}
}
