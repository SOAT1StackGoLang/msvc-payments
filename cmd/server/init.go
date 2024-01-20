package main

import (
	"context"

	"github.com/SOAT1StackGoLang/msvc-payments/pkg/datastore"
	logger "github.com/SOAT1StackGoLang/msvc-payments/pkg/middleware"
)

// initializeApp initializes the application by loading the configuration, connecting to the datastore,
// and subscribing to the Redis channel for receiving messages.
// It returns a pointer to the RedisStore and an error if any.

func initializeApp() (*datastore.RedisStore, error) {
	logger.InitializeLogger()

	// Load the configuration
	logger.Info("Loading configuration...")
	configs, err := LoadConfig()
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("Connecting to datastore...")
	redisStore, err := datastore.NewRedisStore(configs.KVSURI, "", 0)
	if err != nil {
		// handle error
		logger.Error(err.Error())
		return nil, err
	}

	// Subscribe to the Redis channel
	ch, err := redisStore.SubscribeLog(context.Background())
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	go func() {
		for msg := range ch {
			logger.Info("channel msg: " + msg.String())
		}
	}()

	return redisStore, nil
}
