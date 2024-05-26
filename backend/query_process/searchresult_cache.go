package query_process

import (
	. "MariaInfoRetrieval/maria_types"
	"container/list"
	"sync"
)

// Simple in-memory cache
type Cache struct {
	Mu        sync.Mutex
	Cache     map[string]*list.Element
	evictList *list.List
	capacity  int
}

type Entry struct {
	key   string
	value []SearchResult
}

// Create a new cache
func NewCache(capacity int) *Cache {
	return &Cache{
		Cache:     make(map[string]*list.Element),
		evictList: list.New(),
		capacity:  capacity,
	}
}

// Retrieve an item from cache
func (c *Cache) Get(key string) ([]SearchResult, bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if ent, ok := c.Cache[key]; ok {
		c.evictList.MoveToFront(ent)
		return ent.Value.(*Entry).value, true
	}

	return nil, false
}

func (c *Cache) Set(key string, val []SearchResult) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	// If key hit in cache, move it to the front of LRU Cache
	if ent, ok := c.Cache[key]; ok {
		c.evictList.MoveToFront(ent)
		ent.Value.(*Entry).value = val
		return
	}

	// Remove Least recently used item from LRU Cache
	if c.evictList.Len() >= c.capacity {
		ent := c.evictList.Back()
		if ent != nil {
			c.removeElement(ent)
		}
	}

	// Create new entry and add to LRU Cache
	ent := &Entry{key: key, value: val}
	element := c.evictList.PushFront(ent)
	c.Cache[key] = element
}

func (c *Cache) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	kv := e.Value.(*Entry)
	delete(c.Cache, kv.key)
}
