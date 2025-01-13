package main

import (
	"math/rand"
	"strings"
)

type Codec struct {
	base    string
	charset string
	lookup  map[string]string
	mapping map[string]string
}

func Constructor() Codec {
	return Codec{
		base:    "http://tinyurl.com/",
		charset: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		lookup:  make(map[string]string),
		mapping: make(map[string]string),
	}
}

func (c *Codec) encode(url string) string {
	if tiny, exists := c.lookup[url]; exists {
		return tiny
	}

	var key string
	for {
		key = c.genKey()
		if _, exists := c.mapping[key]; !exists {
			break
		}
	}

	tiny := c.base + key
	c.mapping[key] = url
	c.lookup[url] = tiny

	return tiny
}

func (c *Codec) decode(url string) string {
	key := strings.TrimPrefix(url, c.base)
	if original, exists := c.mapping[key]; exists {
		return original
	}
	return ""
}

func (c *Codec) genKey() string {
	key := make([]byte, 6)
	for i := range key {
		key[i] = c.charset[rand.Intn(len(c.charset))]
	}
	return string(key)
}

func main() {
	codec := Constructor()

	tiny := codec.encode("https://leetcode.com/problems/design-tinyurl")
	original := codec.decode(tiny)

	println("Encoded:", tiny)
	println("Decoded:", original)
}
