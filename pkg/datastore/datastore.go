package datastore

import (
	"context"
	"time"

	logger "github.com/SOAT1StackGoLang/msvc-payments/pkg/middleware"
	"github.com/redis/go-redis/v9"
)

type redisStore struct {
	Client *redis.Client
}

//go:generate mockgen -destination=../mocks/datastore_mocks.go -package=mocks github.com/SOAT1StackGoLang/msvc-payments/pkg/datastore RedisStore
type RedisStore interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
	Publish(ctx context.Context, channel string, message interface{}) error
	Subscribe(ctx context.Context, channel string) (<-chan *redis.Message, error)
	SubscribeLog(ctx context.Context) (<-chan *redis.Message, error)
	LPush(ctx context.Context, key string, value interface{}) error
	RPush(ctx context.Context, key string, value any) error
	RetrieveAllFromList(ctx context.Context, key string) ([]string, error)
	BRPop(ctx context.Context, key string) (string, error)
	BLMOVE(ctx context.Context, source string, destination string) (string, error)
	LREM(ctx context.Context, key string, count int64, value interface{}) error
}

// NewRedisStore creates a new RedisStore instance with the given address, password, and database number.
// It establishes a connection to the Redis server and returns the RedisStore object.
// If the connection is successful, it logs a message indicating the successful connection.
// If an error occurs during the connection, it returns nil and the error.

func NewRedisStore(addr string, password string, db int) (RedisStore, error) {
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

	return &redisStore{Client: client}, nil
}

// Set adds a key-value pair to the store
func (s *redisStore) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := s.Client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// Get retrieves a value from the store by its key
func (s *redisStore) Get(ctx context.Context, key string) (string, error) {
	value, err := s.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return value, nil
}

// Exists checks if a key exists in the store
func (s *redisStore) Exists(ctx context.Context, key string) (bool, error) {
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
func (s *redisStore) Delete(ctx context.Context, key string) error {
	err := s.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

// Implement your pub/sub methods here
// Publish sends a message to a channel
func (s *redisStore) Publish(ctx context.Context, channel string, message interface{}) error {
	err := s.Client.Publish(ctx, channel, message).Err()
	if err != nil {
		return err
	}
	return nil
}

// Subscribe subscribes to a channel and returns a channel that receives messages
func (s *redisStore) Subscribe(ctx context.Context, channel string) (<-chan *redis.Message, error) {
	pubsub := s.Client.Subscribe(ctx, channel)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return nil, err
	}

	return pubsub.Channel(), nil
}

// Create a Subscribe method that log all messages received from the channel
func (s *redisStore) SubscribeLog(ctx context.Context) (<-chan *redis.Message, error) {
	pubsub := s.Client.Subscribe(ctx, "log")
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return nil, err
	}

	return pubsub.Channel(), nil
}

// Create a LPUSH method that will add a message to a list
func (s *redisStore) LPush(ctx context.Context, key string, value interface{}) error {
	err := s.Client.LPush(ctx, key, value).Err()
	if err != nil {
		return err
	}
	return nil
}

// RPush appends items into the end of a list defined by its key
func (s *redisStore) RPush(ctx context.Context, key string, value any) error {
	_, err := s.Client.RPush(ctx, key, value).Result()
	return err
}

// RetrieveAllFromList uses LRange to retrieve every value stored in a list defined by its key
func (s *redisStore) RetrieveAllFromList(ctx context.Context, key string) ([]string, error) {
	return s.Client.LRange(ctx, key, 0, -1).Result()
}

// Create a BRPOP/BLPOP method that will remove and return the first element of a list
func (s *redisStore) BRPop(ctx context.Context, key string) (string, error) {
	value, err := s.Client.BRPop(ctx, 0, key).Result()
	if err != nil {
		return "", err
	}
	return value[1], nil
}

// Create a BLMOVE method that will move an element from a list to another list atomically
func (s *redisStore) BLMOVE(ctx context.Context, source string, destination string) (string, error) {
	value, err := s.Client.BLMove(ctx, source, destination, "RIGHT", "LEFT", 1*time.Hour).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

// Create a LREM method that will remove the first count occurrences of elements equal to value from the list stored at key
func (s *redisStore) LREM(ctx context.Context, key string, count int64, value interface{}) error {
	err := s.Client.LRem(ctx, key, count, value).Err()
	if err != nil {
		return err
	}
	return nil
}
