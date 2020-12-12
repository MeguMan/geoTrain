package memcache

import (
	"errors"
	"fmt"
	"os"
	"time"
)

var (
	errNotFound = errors.New("row with this key now wasn't found")
)

func (c *LRU) Set(key string, value string, expiration int64) {
	if element, exists := c.items[key]; exists == true {
		c.queue.MoveToFront(element)
		element.Value.(*Item).Value = value
		return
	}

	if c.queue.Len() == c.capacity {
		c.purge()
	}

	if expiration != 0 {
		expiration += time.Now().Unix()
	}

	item := NewItem(key, value, expiration)

	quit := make(chan bool)
	go func() {
		for _ = range time.Tick(time.Second) {
			select {
			case <- quit:
				return
			default:
				if item.expired() {
					delete(c.items, item.Key)
					quit <- true
				}
			}
		}
	}()

	element := c.queue.PushFront(item)
	c.items[item.Key] = element
}

func (c *LRU) Get(key string) (string, error) {
	element, exists := c.items[key]
	if exists == false {
		return "", errNotFound
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