package kimg

import "log"

// KimgCache is a interface to provide cache in kimg.
type KimgCache interface {
	Set(key string, data []byte) error
	Get(key string) ([]byte, error)
	Del(key string) error
}

// NewKimgCache create a cache instance according to cache mode in config.
func NewKimgCache(config *KimgConfig) (KimgCache, error) {
	switch config.Cache.Mode {
	case "none":
		log.Println("[INFO] cache disabled")
		return nil, nil
	case "memory":
		log.Println("[INFO] cache [memory] used")
		return NewKimgMemoryCache(config)
	case "memcache":
		log.Println("[INFO] cache [memcache] used")
		return NewKimgMemcacheCache(config)
	case "redis":
		log.Println("[INFO] cache [redis] used")
		return NewKimgRedisCache(config)
	default:
		log.Printf("[WARN] unsupported cache mode :%s\n", config.Cache.Mode)
		return nil, nil
	}
}
