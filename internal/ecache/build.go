package ecache

import "github.com/orca-zhang/ecache"

type Cache struct {
	cache *ecache.Cache
}

func NewCache() *Cache {
	lru2 := ecache.NewLRUCache(16, 200, 0).LRU2(1024)
	return &Cache{
		cache: lru2,
	}
}
