package memcache

import (
	"container/list"
)

type HashItem struct {
	Key   string
	Value map[interface{}] interface{}
	TTL   int64
}

type Item struct {
	Key   string
	Value interface{}
	TTL   int64
}

type LRU struct {
	capacity int
	items    map[string]*list.Element
	queue    *list.List
	password string
}

func NewLru(config *Config) *LRU {
	return &LRU{
		capacity: config.Capacity,
		items:    make(map[string]*list.Element),
		queue:    list.New(),
		password: config.Password,
	}
}

func NewItem(key string, value interface{}, ttl int64) *Item {
	return &Item{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}
}

func NewHashItem(key string, value map[interface{}] interface{}, ttl int64) *HashItem {
	return &HashItem{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}
}

func (c *LRU) purge() {
	if element := c.queue.Back(); element != nil {
		item := c.queue.Remove(element).(*Item)
		delete(c.items, item.Key)
	}
}

