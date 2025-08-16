package cache

import "sync"

// Cache para manejar respuestas paginadas grandes
type CacheEntry struct {
	Data      interface{}
	Timestamp int64
}

type Cache struct {
	Mu       sync.RWMutex
	Registry map[string]CacheEntry
}

func NewCache() *Cache {
	cacheRegistry := map[string]CacheEntry{}
	return &Cache{
		Registry: cacheRegistry,
	}
}
