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

// NewRedisStore creates a new RedisStore instance with the given address, password, and database number.
// It establishes a connection to the Redis server and returns the RedisStore object.
// If the connection is successful, it logs a message indicating the successful connection.
// If an error occurs during the connection, it returns nil and the error.

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

// Exists checks if a key exists in the store
func (s *RedisStore) Exists(ctx context.Context, key string) (bool, error) {
	value, err := s.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if value == 0 {
		return false, nil
	}
	return true, nil
}

// Delete removes a key-value pair from the store
func (s *RedisStore) Delete(ctx context.Context, key string) error {
	err := s.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
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

// Create a LPUSH method that will add a message to a list
func (s *RedisStore) LPush(ctx context.Context, key string, value interface{}) error {
	err := s.Client.LPush(ctx, key, value).Err()
	if err != nil {
		return err
	}
	return nil
}

// Create a BRPOP/BLPOP method that will remove and return the first element of a list
func (s *RedisStore) BRPop(ctx context.Context, key string) (string, error) {
	value, err := s.Client.BRPop(ctx, 0, key).Result()
	if err != nil {
		return "", err
	}
	return value[1], nil
}

// Create a BLMOVE method that will move an element from a list to another list atomically
func (s *RedisStore) BLMOVE(ctx context.Context, source string, destination string) (string, error) {
	value, err := s.Client.BLMove(ctx, source, destination, "RIGHT", "LEFT", 1*time.Hour).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

// Create a LREM method that will remove the first count occurrences of elements equal to value from the list stored at key
func (s *RedisStore) LREM(ctx context.Context, key string, count int64, value interface{}) error {
	err := s.Client.LRem(ctx, key, count, value).Err()
	if err != nil {
		return err
	}
	return nil
}
