package caches

import "sync"

type Str struct {
	cacheMutex sync.Mutex
	cache      map[string]bool
}

func NewStr() Str {
	return Str{
		cache: make(map[string]bool),
	}
}

func (c *Str) Add(str string) {
	defer c.cacheMutex.Unlock()
	c.cacheMutex.Lock()

	c.cache[str] = true
}

func (c *Str) Remove(str string) {
	defer c.cacheMutex.Unlock()
	c.cacheMutex.Lock()

	delete(c.cache, str)
}

func (c *Str) Has(str string) bool {
	defer c.cacheMutex.Unlock()
	c.cacheMutex.Lock()

	_, hasKey := c.cache[str]
	return hasKey
}

func (c *Str) Count() int {
	defer c.cacheMutex.Unlock()
	c.cacheMutex.Lock()

	return len(c.cache)
}