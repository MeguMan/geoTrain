package memcache

import (
	"errors"
	"fmt"
	"os"
	"time"
)

var (
	errNotFound = errors.New("row with this key wasn't found")
)

func (c *LRU) Set(key string, value string, ttl int64) {
	if element, exists := c.items[key]; exists == true {
		item := element.Value.(*Item)
		c.queue.MoveToFront(element)
		item.Value = value

		if ttl != 0 {
			item.TTL = ttl
			go c.deleteAfterExpiration(item)
			return
		}
		return
	}

	if c.queue.Len() == c.capacity {
		c.purge()
	}

	if ttl != 0 {
		item := NewItem(key, value, ttl)
		go c.deleteAfterExpiration(item)
		element := c.queue.PushFront(item)
		c.items[item.Key] = element
	} else {
		item := NewItem(key, value, ttl)
		element := c.queue.PushFront(item)
		c.items[item.Key] = element
	}
}

func (c *LRU) Hset(key string, value interface{}, expiration int64) {

}

func (c *LRU) Get(key string) (interface{}, error) {
	element, exists := c.items[key]
	if exists == false {
		return nil, errNotFound
	}
	c.queue.MoveToFront(element)

	return element.Value.(*Item).Value, nil
}

func (c *LRU) GetAllKeys() []string {
	var keys []string
	for k, _ := range c.items {
		keys = append(keys, k)
	}
	return keys
}

func (c *LRU) Save() error{
	var data string
	for k, v := range c.items {
		data += fmt.Sprintf("%s - %s \n", k, v.Value.(*Item).Value)
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

func (c *LRU) Delete(key string) error{
	_, exists := c.items[key]
	if exists == false {
		return errNotFound
	}

	delete(c.items, key)
	return nil
}

func (c *LRU) CheckPassword(password string) bool {
	if c.password == password {
		return true
	}
	return false
}

func (c *LRU) deleteAfterExpiration(item *Item) {
	for _ = range time.Tick(time.Second) {
		if item.TTL == 0 {
			delete(c.items, item.Key)
			break
		} else {
			item.TTL -= 1
		}
	}
}