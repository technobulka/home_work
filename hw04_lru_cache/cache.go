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
	items    map[Key]cacheItem
}

type cacheItem struct {
	key   *listItem
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if item, ok := c.items[key]; ok {
		c.items[key] = cacheItem{item.key, value}
		c.queue.MoveToFront(item.key)
		return true
	}

	if c.queue == nil {
		c.Clear()
	}

	var newItem = cacheItem{
		value: value,
	}

	newItem.key = c.queue.PushFront(newItem)
	c.items[key] = newItem

	if c.queue.Len() > c.capacity {
		for mapKey, item := range c.items {
			if item.key == c.queue.Back() {
				delete(c.items, mapKey)
			}
		}

		c.queue.Remove(c.queue.Back())
	}

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item.key)
		return item.value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = &list{}
	c.items = make(map[Key]cacheItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{capacity: capacity}
}
