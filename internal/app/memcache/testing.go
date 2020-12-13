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
	element := l.queue.PushFront(item)
	l.items[item.Key] = element
	return &l
}
