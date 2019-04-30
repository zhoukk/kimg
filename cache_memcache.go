package main

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

type kimgMemcacheCache struct {
	client *memcache.Client
}

// NewKimgMemcacheCache create a memcache cache instance.
func NewKimgMemcacheCache(config *KimgConfig) (KimgCache, error) {
	addr := fmt.Sprintf("%s:%d", config.Cache.MemcacheHost, config.Cache.MemcachePort)
	client := memcache.New(addr)

	return &kimgMemcacheCache{
		client: client,
	}, nil
}

func (cache *kimgMemcacheCache) Release() {}

func (cache *kimgMemcacheCache) Set(key string, data []byte) error {
	it := &memcache.Item{Key: key, Value: data}
	return cache.client.Set(it)
}

func (cache *kimgMemcacheCache) Get(key string) ([]byte, error) {
	it, err := cache.client.Get(key)
	if err != nil {
		return nil, err
	}
	return it.Value, nil
}

func (cache *kimgMemcacheCache) Del(key string) error {
	return cache.client.Delete(key)
}
