package datastore

import (
	"time"

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

// Set adds a key-value pair to the store
func (s *RedisStore) Set(key string, value interface{}, expiration time.Duration) error {
	err := s.Client.Set(key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// Get retrieves a value from the store by its key
func (s *RedisStore) Get(key string) (string, error) {
	value, err := s.Client.Get(key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return value, nil
}

// Implement your pub/sub methods here
// Publish sends a message to a channel
func (s *RedisStore) Publish(channel string, message interface{}) error {
	err := s.Client.Publish(channel, message).Err()
	if err != nil {
		return err
	}
	return nil
}

// Subscribe subscribes to a channel and returns a channel that receives messages
func (s *RedisStore) Subscribe(channel string) (<-chan *redis.Message, error) {
	pubsub := s.Client.Subscribe(channel)
	_, err := pubsub.Receive()
	if err != nil {
		return nil, err
	}

	return pubsub.Channel(), nil
}
