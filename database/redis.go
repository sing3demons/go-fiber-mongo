package database

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sing3demons/go-fiber-mongo/models"
)

type redisCache struct {
	host    string
	db      int
	expires time.Duration
}

type RedisCache interface {
	Set(key string, value interface{}) error
	GetProducts(key string) ([]models.Product, error)
}

func NewRedisCache(host string, db int, expires time.Duration) RedisCache {
	return &redisCache{
		host:    host,
		db:      db,
		expires: expires,
	}
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: os.Getenv("RDB_PASSWORD"),
		DB:       cache.db,
	})
}

func (cache *redisCache) Set(key string, value interface{}) error {
	rdb := cache.getClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err := rdb.Set(ctx, key, json, cache.expires*time.Second).Err(); err != nil {
		return err
	}

	return nil
}
func (cache *redisCache) GetProducts(key string) ([]models.Product, error) {
	rdb := cache.getClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	product := []models.Product{}
	if err := json.Unmarshal([]byte(val), &product); err != nil {
		panic(err)
	}

	return product, nil
}
