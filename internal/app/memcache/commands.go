package memcache

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"
)

var (
	errNotFound = errors.New("row with this key wasn't found")
	errWrongType = errors.New("operation against a key holding the wrong kind of value")
)

func (l *LRU) Set(key string, value string, ttl int64) {
	if element, exists := l.items[key]; exists {
		if reflect.TypeOf(element.Value) == reflect.TypeOf(new(HashItem)) {
			item := NewItem(key, value, ttl)
			element.Value = item
			if ttl != 0 {
				go l.deleteAfterExpiration(item)
			}
			l.items[item.Key] = element
			l.queue.MoveToFront(element)
			return
		}
		item := element.Value.(*Item)
		l.queue.MoveToFront(element)
		item.Value = value

		if ttl != 0 {
			item.TTL = ttl
			go l.deleteAfterExpiration(item)
			return
		}
		return
	}
	if l.queue.Len() == l.capacity {
		l.purge()
	}

	if ttl != 0 {
		item := NewItem(key, value, ttl)
		go l.deleteAfterExpiration(item)
		element := l.queue.PushFront(item)
		l.items[item.Key] = element
	} else {
		item := NewItem(key, value, ttl)
		element := l.queue.PushFront(item)
		l.items[item.Key] = element
	}
}

func (l *LRU) HSet(hash string, field, value interface{}) error{
	if element, exists := l.items[hash]; exists {
		if reflect.TypeOf(element.Value) != reflect.TypeOf(new(HashItem)) {
			return errWrongType
		}
		element.Value.(*HashItem).Value[field] = value
		l.queue.MoveToFront(element)
		return nil
	}

	if l.queue.Len() == l.capacity {
		l.purge()
	}
	m := map[interface{}]interface{}{
		field: value,
	}
	item := NewHashItem(hash, m, 0)
	element := l.queue.PushFront(item)
	l.items[item.Key] = element
	return nil
}

func (l *LRU) Get(key string) (interface{}, error) {
	element, exists := l.items[key]
	if !exists {
		return nil, errNotFound
	}
	if reflect.TypeOf(element.Value) != reflect.TypeOf(new(Item)) {
		return nil, errWrongType
	}
	l.queue.MoveToFront(element)

	return element.Value.(*Item).Value, nil
}

func (l *LRU) HGet(hash string, field interface{}) (interface{}, error) {
	element, exists := l.items[hash]
	if !exists {
		return nil, errNotFound
	}
	if reflect.TypeOf(element.Value) != reflect.TypeOf(new(HashItem)) {
		return nil, errWrongType
	}
	value, exists := element.Value.(*HashItem).Value[field]
	if !exists {
		return nil, errNotFound
	}

	l.queue.MoveToFront(element)
	return value, nil
}

func (l *LRU) GetAllKeys() []string {
	var keys []string
	for k, _ := range l.items {
		keys = append(keys, k)
	}
	return keys
}

func (l *LRU) Save() error{
	var data string
	for k, v := range l.items {
		if reflect.TypeOf(v.Value) == reflect.TypeOf(new(Item)) {
			data += fmt.Sprintf("%s - %s \n", k, v.Value.(*Item).Value)
		} else {
			data += fmt.Sprintf("%s - %s \n", k, v.Value.(*HashItem).Value)
		}
	}

	file, err := os.Create("data.txt")
	if err != nil{
		fmt.Println("Unable to create file")
		return err
	}
	defer file.Close()

	file.WriteString(data)
	return nil
}

func (l *LRU) Delete(key string) error{
	_, exists := l.items[key]
	if !exists {
		return errNotFound
	}

	delete(l.items, key)
	return nil
}

func (l *LRU) CheckPassword(password string) bool {
	if l.password == password {
		return true
	}
	return false
}

func (l *LRU) deleteAfterExpiration(item *Item) {
	for _ = range time.Tick(time.Second) {
		if item.TTL == 0 {
			delete(l.items, item.Key)
			break
		} else {
			item.TTL -= 1
		}
	}
}