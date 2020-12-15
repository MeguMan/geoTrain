package client

import (
	"fmt"
	"net/http"
)

func (c *Client) Login(password string) (interface{}, error) {
	path := fmt.Sprintf("/login?password=%s", password)
	resp, err := http.Get(c.BaseURL+path)
	fmt.Println(err)
	if err != nil {
		return resp.StatusCode, err
	}
	c.SessionCookie = resp.Cookies()[0].Name + "=" + resp.Cookies()[0].Value
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func (c *Client) Set(key string, value string, ttl int64) (interface{}, error) {
	path := fmt.Sprintf("/rows?key=%s&value=%s&ttl=%d", key, value, ttl)
	req, err := http.NewRequest("POST", c.BaseURL+path, nil)
	req.Header.Set("Cookie", c.SessionCookie)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	return  resp.Body, nil
}

func (c *Client) Get(key string) (interface{}, error) {
	path := fmt.Sprintf("/rows/%s", key)
	req, err := http.NewRequest("GET", c.BaseURL+path, nil)
	req.Header.Set("Cookie", c.SessionCookie)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	return resp.Body, nil
}

func (c *Client) HSet(hash string, field, value interface{}) (interface{}, error) {
	path := fmt.Sprintf("/rows/hash?hash=%s&field=%s&value=%s", hash, field, value)
	req, err := http.NewRequest("POST", c.BaseURL+path, nil)
	req.Header.Set("Cookie", c.SessionCookie)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	return resp.Body, nil
}


func (c *Client) HGet(hash string, field interface{}) (interface{}, error) {
	path := fmt.Sprintf("/rows/hash/%s/%s", hash, field)
	req, err := http.NewRequest("GET", c.BaseURL+path, nil)
	req.Header.Set("Cookie", c.SessionCookie)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	return resp.Body, nil
}

func (c *Client) Delete(key string) (interface{}, error) {
	path := fmt.Sprintf("/rows/%s", key)
	req, err := http.NewRequest("DELETE", c.BaseURL+path, nil)
	req.Header.Set("Cookie", c.SessionCookie)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	return resp.Body, nil
}

func (c *Client) Save() error {
	path := fmt.Sprintf("/save")
	req, err := http.NewRequest("GET", c.BaseURL+path, nil)
	if err != nil {
		return err
	}
	_, err = c.HTTPClient.Do(req)
	return nil
}
