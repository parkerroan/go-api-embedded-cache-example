package storage

import (
	"time"

	"github.com/dgraph-io/ristretto"
)

type Cache interface {
	Set(id string, item Item) error
	Retreive(id string) (*Item, error)
	Remove(id string) error
}

type CacheClient struct {
	cache *ristretto.Cache
}

func NewCacheClient(cache *ristretto.Cache) *CacheClient {
	return &CacheClient{
		cache: cache,
	}
}

func (c *CacheClient) Set(id string, item Item) error {

	// set item in cache
	c.cache.SetWithTTL(id, item, 1, time.Duration(12*time.Hour))

	return nil
}

func (c *CacheClient) Retreive(id string) (*Item, error) {

	// get item from cache
	if item, found := c.cache.Get(id); found {
		item := item.(Item)
		return &item, nil
	}

	return nil, nil
}

func (c *CacheClient) Remove(id string) error {
	// remove item from cache
	c.cache.Del(id)

	return nil
}
