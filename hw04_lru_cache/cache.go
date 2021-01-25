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
	key   Key
	ptr   *listItem
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if item, ok := c.items[key]; ok {
		if item.ptr != c.queue.Front() {
			c.queue.MoveToFront(item.ptr)
		}

		if item.value != value {
			c.items[key] = cacheItem{
				key,
				c.queue.Front(),
				value,
			}
		}

		return true
	}

	if c.queue == nil {
		c.Clear()
	}

	var newItem = cacheItem{key, nil, value}
	newItem.ptr = c.queue.PushFront(newItem)
	c.items[key] = newItem

	if c.queue.Len() > c.capacity {
		var last = c.queue.Back().value.(cacheItem)
		delete(c.items, last.key)
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
	c.queue = &list{}
	c.items = make(map[Key]cacheItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{capacity: capacity}
}
