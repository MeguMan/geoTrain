package memcache

import (
	"container/list"
	"testing"
)

func TestLru(t *testing.T) *LRU {
	t.Helper()
	l := LRU{
		capacity: 3,
		items: make(map[string]*list.Element),
		queue:    list.New(),
		password: "strongpassword",
	}
	item := &Item{
		Key:   "firstKey",
		Value: "ItemsValue",
		TTL:   0,
	}
	hItem := &HashItem{
		Key:   "testHash",
		Value: map[interface{}]interface{}{
			"testField": "testValue",
		},
	}
	element := l.queue.PushFront(item)
	l.items[item.Key] = element
	element2 := l.queue.PushFront(hItem)
	l.items[hItem.Key] = element2
	return &l
}
