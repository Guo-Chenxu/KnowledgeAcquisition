package logic

import (
	"KnowledgeAcquisition/model"
	"container/list"
	"sync"
)

type Cache struct {
	Mu        sync.Mutex
	Cache     map[string]*list.Element
	evictList *list.List
	capacity  int
}

type Entry struct {
	key   string
	value []model.SearchResult
}

func NewCache(capacity int) *Cache {
	return &Cache{
		Cache:     make(map[string]*list.Element),
		evictList: list.New(),
		capacity:  capacity,
	}
}

func (c *Cache) Get(key string) ([]model.SearchResult, bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if ent, ok := c.Cache[key]; ok {
		c.evictList.MoveToFront(ent)
		return ent.Value.(*Entry).value, true
	}

	return nil, false
}

func (c *Cache) Set(key string, val []model.SearchResult) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if ent, ok := c.Cache[key]; ok {
		c.evictList.MoveToFront(ent)
		ent.Value.(*Entry).value = val
		return
	}

	if c.evictList.Len() >= c.capacity {
		ent := c.evictList.Back()
		if ent != nil {
			c.removeElement(ent)
		}
	}

	ent := &Entry{key: key, value: val}
	element := c.evictList.PushFront(ent)
	c.Cache[key] = element
}

func (c *Cache) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	kv := e.Value.(*Entry)
	delete(c.Cache, kv.key)
}
