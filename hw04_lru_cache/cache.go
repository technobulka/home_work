package hw04_lru_cache //nolint:golint,stylecheck
import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*cacheItem
	mutex    *sync.Mutex
}

type cacheItem struct {
	key   Key
	ptr   *listItem
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item, ok := c.items[key]; ok {
		c.items[key].value = value
		c.queue.MoveToFront(item.ptr)
		return true
	}

	ptr := c.queue.PushFront(key)
	c.items[key] = &cacheItem{key, ptr, value}

	if c.queue.Len() > c.capacity {
		delete(c.items, c.queue.Back().value.(Key))
		c.queue.Remove(c.queue.Back())
	}

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item.ptr)
		return item.value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = new(list)
	c.items = make(map[Key]*cacheItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity,
		new(list),
		make(map[Key]*cacheItem, capacity),
		new(sync.Mutex),
	}
}
