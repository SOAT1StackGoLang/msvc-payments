package datastore

import (
	"github.com/go-redis/redis"
)

type RedisStore struct {
	Client *redis.Client
}

func NewRedisStore(addr string, password string, db int) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &RedisStore{Client: client}, nil
}
