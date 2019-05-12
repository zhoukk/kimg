package main

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type kimgMemcacheCache struct {
	client *memcache.Client
}

// NewKimgMemcacheCache create a memcache cache instance.
func NewKimgMemcacheCache(config *KimgConfig) (KimgCache, error) {
	client := memcache.New(config.Cache.Memcache.URL)

	return &kimgMemcacheCache{
		client: client,
	}, nil
}

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
