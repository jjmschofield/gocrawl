package caches

import (
	"github.com/go-redis/redis"
	"log"
	"sync"
)

type StrRedis struct { // TODO - the cache interface swallows errors - an extension would be beneficial to allow consumers to respond
	cacheMutex sync.Mutex
	client     *redis.Client
	setKey     string
	cache      map[string]bool
}

func NewStrRedis(setKey string, options *redis.Options) StrRedis {
	client := redis.NewClient(options)

	_, err := client.Ping().Result()
	if err != nil {
		log.Panicf("couldn't connect to reddis %s", err)
	}

	return StrRedis{
		client: client,
		setKey: setKey,
	}
}

func (c *StrRedis) Add(str string) {
	defer c.cacheMutex.Unlock()
	c.cacheMutex.Lock()

	c.client.SAdd(c.setKey, str)
}

func (c *StrRedis) Remove(str string) {
	defer c.cacheMutex.Unlock()
	c.cacheMutex.Lock()

	c.client.SRem(c.setKey, str)
}

func (c *StrRedis) Has(str string) bool {
	defer c.cacheMutex.Unlock()
	c.cacheMutex.Lock()

	hasKey, err := c.client.SIsMember(c.setKey, str).Result()
	if err != nil {
		log.Panicf("couldn't check is member for set for %s %v", c.setKey, err)
	}

	return hasKey
}

func (c *StrRedis) Count() int {
	defer c.cacheMutex.Unlock()
	c.cacheMutex.Lock()

	count, err := c.client.SCard(c.setKey).Result()
	if err != nil {
		log.Panicf("couldn't get count for set for %s %v", c.setKey, err)
	}

	return int(count) // TODO interface should be extended to int64 to handle larger numbers
}
