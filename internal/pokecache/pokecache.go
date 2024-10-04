package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	cacheMap map[string]cacheEntry
	mutex    sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	newCache := &Cache{
		cacheMap: make(map[string]cacheEntry),
		interval: interval,
	}
	go newCache.reapLoop()
	return newCache
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cacheMap[key] = cacheEntry{
		time.Now(),
		val,
	}
}

func (c *Cache) Get(key string) (val []byte, found bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, found := c.cacheMap[key]

	return entry.val, found
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mutex.Lock()

			now := time.Now()
			for key, entry := range c.cacheMap {
				if now.Sub(entry.createdAt) >= c.interval {
					delete(c.cacheMap, key)
				}
			}

			c.mutex.Unlock()
		}
	}
}

func TestFunc() {
	fmt.Println("Imported pokecache")
}
