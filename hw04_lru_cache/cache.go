package hw04_lru_cache //nolint:golint,stylecheck

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
}

type cacheItem struct {
	key   Key
	ptr   *listItem
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if item, ok := c.items[key]; ok {
		c.items[key].value = value
		c.queue.MoveToFront(item.ptr)
		return true
	}

	if c.queue == nil {
		c.init()
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
	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item.ptr)
		return item.value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.init()
}

func (c *lruCache) init() {
	c.queue = &list{}
	c.items = make(map[Key]*cacheItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{capacity: capacity}
}
