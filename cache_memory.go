package main

import (
	"container/list"
	"errors"
	"sync"
)

type kimgMemoryCache struct {
	mtx      sync.Mutex
	list     *list.List
	table    map[string]*list.Element
	size     int64
	capacity int64
}

type cacheEntry struct {
	key  string
	data []byte
	size int64
}

// NewKimgMemoryCache create a memory cache instance.
func NewKimgMemoryCache(config *KimgConfig) (KimgCache, error) {
	return &kimgMemoryCache{
		list:     list.New(),
		table:    make(map[string]*list.Element),
		capacity: config.Cache.Memory.Capacity,
	}, nil
}

func (cache *kimgMemoryCache) Set(key string, data []byte) error {
	cache.mtx.Lock()
	defer cache.mtx.Unlock()

	if ele := cache.table[key]; ele != nil {
		cache.updateInplace(ele, data)
	} else {
		cache.addNew(key, data)
	}
	return nil
}

func (cache *kimgMemoryCache) Get(key string) ([]byte, error) {
	cache.mtx.Lock()
	defer cache.mtx.Unlock()

	ele, ok := cache.table[key]
	if ele == nil || !ok {
		return nil, errors.New("memory cache miss")
	}
	cache.moveToFront(ele)

	return ele.Value.(*cacheEntry).data, nil
}

func (cache *kimgMemoryCache) Del(key string) error {
	cache.mtx.Lock()
	defer cache.mtx.Unlock()

	ele, ok := cache.table[key]
	if ele == nil || !ok {
		return errors.New("memory cache miss")
	}

	cache.list.Remove(ele)
	delete(cache.table, key)
	cache.size -= ele.Value.(*cacheEntry).size

	return nil
}

func (cache *kimgMemoryCache) updateInplace(ele *list.Element, data []byte) {
	cache.size += int64(len(data)) - ele.Value.(*cacheEntry).size
	ele.Value.(*cacheEntry).data = data
	ele.Value.(*cacheEntry).size = int64(len(data))
	cache.moveToFront(ele)
	cache.checkCapacity()
}

func (cache *kimgMemoryCache) moveToFront(ele *list.Element) {
	cache.list.MoveToFront(ele)
}

func (cache *kimgMemoryCache) addNew(key string, data []byte) {
	entry := &cacheEntry{key, data, int64(len(data))}
	ele := cache.list.PushFront(entry)
	cache.table[key] = ele
	cache.size += entry.size
	cache.checkCapacity()
}

func (cache *kimgMemoryCache) checkCapacity() {
	for cache.size > cache.capacity {
		ele := cache.list.Back()
		entry := ele.Value.(*cacheEntry)
		cache.list.Remove(ele)
		delete(cache.table, entry.key)
		cache.size -= entry.size
	}
}
