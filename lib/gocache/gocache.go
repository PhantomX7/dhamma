package gocache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type Client interface {
	Set(key string, value string)
	Get(key string) string
}

type Cache struct {
	cache *cache.Cache
}

func (c *Cache) Set(key string, value string) {
	c.cache.Set(key, value, cache.NoExpiration)
}

func (c *Cache) Get(key string) string {
	value, _ := c.cache.Get(key)
	if value == nil {
		return ""
	}
	return value.(string)
}

func New() Client {
	return &Cache{
		cache: cache.New(cache.NoExpiration, 1*time.Hour),
	}
}
