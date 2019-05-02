package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type kimgRedisCache struct {
	pool *redis.Pool
}

// NewKimgRedisCache create a redis cache instance.
func NewKimgRedisCache(config *KimgConfig) (KimgCache, error) {
	addr := fmt.Sprintf("%s:%d", config.Cache.RedisHost, config.Cache.RedisPort)
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return &kimgRedisCache{
		pool: pool,
	}, nil
}

func (cache *kimgRedisCache) Release() {}

func (cache *kimgRedisCache) getConnect() (redis.Conn, error) {
	conn := cache.pool.Get()
	if conn == nil {
		return nil, errors.New("can not connect to redis server")
	}
	return conn, nil
}

func (cache *kimgRedisCache) Set(key string, data []byte) error {
	conn, err := cache.getConnect()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = redis.Bytes(conn.Do("SET", key, data))
	return err
}

func (cache *kimgRedisCache) Get(key string) ([]byte, error) {
	conn, err := cache.getConnect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (cache *kimgRedisCache) Del(key string) error {
	conn, err := cache.getConnect()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = redis.Bytes(conn.Do("DEL", key))
	return err
}
