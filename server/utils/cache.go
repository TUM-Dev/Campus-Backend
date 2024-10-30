package utils

import (
	"sync"
	"time"
)

type CacheKey string

const (
	// CacheKeyAllNewsSources is the key for all news sources from the database
	CacheKeyAllNewsSources CacheKey = "all_news_sources"
	// CacheKeyNews is the key for news
	CacheKeyNews CacheKey = "news"
)

// Cache is a thread-safe cache
type Cache struct {
	cache      map[string]any
	deleteJobs map[string]time.Time
	lock       sync.RWMutex
}

// NewCache creates a new cache
func NewCache() *Cache {
	cache := &Cache{
		cache:      make(map[string]any),
		deleteJobs: make(map[string]time.Time),
		lock:       sync.RWMutex{},
	}
	go cache.deleteLoop()
	return cache
}

// Set adds an entry to the cache
func (c *Cache) Set(key CacheKey, params string, value any, expire time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache[c.combine(key, params)] = value
	c.deleteJobs[c.combine(key, params)] = time.Now().Add(expire)
}

// Get returns an entry from the cache
func (c *Cache) Get(key CacheKey, params string) any {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cache[c.combine(key, params)]
}

// Exists checks if an entry exists in the cache
func (c *Cache) Exists(key CacheKey, params string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok := c.cache[c.combine(key, params)]
	return ok
}

// Delete removes an entry from the cache
func (c *Cache) delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.cache, key)
	delete(c.deleteJobs, key)
}

// deleteLoop deletes all entries that have expired from the cache
func (c *Cache) deleteLoop() {
	for {
		for s, t := range c.deleteJobs {
			if time.Now().After(t) {
				c.delete(s)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func (c *Cache) combine(key CacheKey, params string) string {
	return string(key) + params
}
