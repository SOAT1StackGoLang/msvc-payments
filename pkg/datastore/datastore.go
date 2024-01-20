package datastore

import (
	"context"
	"time"

	logger "github.com/SOAT1StackGoLang/msvc-payments/pkg/middleware"
	"github.com/redis/go-redis/v9"
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

	ping, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	if ping == "PONG" {
		logger.Info("Connected to Redis: PONG!")
	}

	return &RedisStore{Client: client}, nil
}

// Set adds a key-value pair to the store
func (s *RedisStore) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := s.Client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// Get retrieves a value from the store by its key
func (s *RedisStore) Get(ctx context.Context, key string) (string, error) {
	value, err := s.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return value, nil
}

// Implement your pub/sub methods here
// Publish sends a message to a channel
func (s *RedisStore) Publish(ctx context.Context, channel string, message interface{}) error {
	err := s.Client.Publish(ctx, channel, message).Err()
	if err != nil {
		return err
	}
	return nil
}

// Subscribe subscribes to a channel and returns a channel that receives messages
func (s *RedisStore) Subscribe(ctx context.Context, channel string) (<-chan *redis.Message, error) {
	pubsub := s.Client.Subscribe(ctx, channel)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return nil, err
	}

	return pubsub.Channel(), nil
}

// Create a Subscribe method that log all messages received from the channel
func (s *RedisStore) SubscribeLog(ctx context.Context) (<-chan *redis.Message, error) {
	pubsub := s.Client.Subscribe(ctx, "log")
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return nil, err
	}

	return pubsub.Channel(), nil
}
