package redis_implementation

import (
	"time"
)

func (c *LRU) Set(key string, value interface{}, expiration int64) bool {
	if element, exists := c.items[key]; exists == true {
		c.queue.MoveToFront(element)
		element.Value.(*Item).Value = value
		return true
	}

	if c.queue.Len() == c.capacity {
		c.purge()
	}

	if expiration != 0 {
		expiration += time.Now().Unix()
	}

	item := NewItem(key, value, expiration)

	go func() {
		for now := range time.Tick(time.Second) {
			c.mutex.Lock()
			for k, _ := range c.items {
				if now.Unix() > item.TTL {
					delete(c.items, k)
				}
			}
		}
	}()

	element := c.queue.PushFront(item)
	c.items[item.Key] = element

	return true
}

func (c *LRU) Get(key string) interface{} {
	element, exists := c.items[key]
	if exists == false {
		return nil
	}
	c.queue.MoveToFront(element)

	return element.Value.(*Item).Value
}