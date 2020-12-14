package memcache

import (
	"container/list"
	"reflect"
)

type HashItem struct {
	Key   string
	Value map[interface{}] interface{}
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

func NewHashItem(key string, value map[interface{}] interface{}) *HashItem {
	return &HashItem{
		Key:   key,
		Value: value,
	}
}

func (l *LRU) purge() {
	if element := l.queue.Back(); element != nil {
		if reflect.TypeOf(element.Value) != reflect.TypeOf(new(Item)) {
			item := l.queue.Remove(element).(*HashItem)
			delete(l.items, item.Key)
		} else {
			item := l.queue.Remove(element).(*Item)
			delete(l.items, item.Key)
		}
	}
}

