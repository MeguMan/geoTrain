package memcache

import (
	"fmt"
	"os"
	"time"
)

func (c *LRU) Set(key string, value interface{}, expiration int64) {
	if element, exists := c.items[key]; exists == true {
		fmt.Println("EXISTS")
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

func (c *LRU) Get(key string) interface{} {
	element, exists := c.items[key]
	if exists == false {
		return nil
	}
	c.queue.MoveToFront(element)

	return element.Value.(*Item).Value
}

func (c *LRU) Save() {
	var data string
	for k, v := range c.items {
		data += fmt.Sprintf("%s - %s \n", k, v.Value.(*Item).Value)
	}

	file, err := os.Create("data.txt")
	if err != nil{
		fmt.Println("Unable to create file:", err)
	}
	defer file.Close()

	file.WriteString(data)
}

func (c *LRU) CheckPassword(password string) bool {
	if c.password == password {
		return true
	}
	return false
}