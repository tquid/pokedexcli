package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	reapInterval time.Duration
	entries      map[string]cacheEntry
	mu           sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	c := Cache{
		reapInterval: interval,
		entries:      make(map[string]cacheEntry),
	}
	c.reapLoop()
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, exists := c.entries[key]
	if !exists {
		return []byte{}, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.reapInterval)
	go func() {
		for {
			<-ticker.C
			c.mu.Lock()
			for key, entry := range c.entries {
				if time.Since(entry.createdAt) > c.reapInterval {
					delete(c.entries, key)
				}
			}
			c.mu.Unlock()
		}
	}()
}
