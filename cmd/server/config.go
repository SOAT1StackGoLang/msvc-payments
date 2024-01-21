// Create a function to load the configuration from environment variables using go-kit
package main

import (
	"os"
	"strconv"
)

// Config is a struct to hold the configuration
type Config struct {
	KVSHost string `envconfig:"KVSTORE_HOST"`
	KVSPort int    `envconfig:"KVSTORE_PORT"`
	// optional configs
	KVSURI string `envconfig:"KVSTORE_URI"`
	KVSDB  int    `envconfig:"KVSTORE_DB"`
}

// LoadConfig loads the configuration values for the server.
// It retrieves the values from environment variables and sets default values if necessary.
// The configuration includes the KVStore host, port, and URI.
// If the environment variables are not set or invalid, default values are used.
// The function returns the loaded configuration and an error if any.
func LoadConfig() (Config, error) {
	var cfg Config

	// Load KVSHost
	cfg.KVSHost = os.Getenv("KVSTORE_HOST")
	if cfg.KVSHost == "" {
		cfg.KVSHost = "localhost" // Set default value for KVSTORE_HOST
	}

	// Load KVSPort
	portStr := os.Getenv("KVSTORE_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		// Set default value if KVSTORE_PORT is not set or invalid
		port = 6379
	}
	cfg.KVSPort = port

	// Load KVSURI (optional)
	cfg.KVSURI = os.Getenv("KVSTORE_URI")

	// Set default value for KVSURI if not set
	if cfg.KVSURI == "" {
		cfg.KVSURI = cfg.KVSHost + ":" + strconv.Itoa(cfg.KVSPort)
	}

	// KVSDB is not used in this example, but it can be added as an optional config
	kvsStr := os.Getenv("KVSTORE_DB")
	kvsDB, err := strconv.Atoi(kvsStr)
	if err != nil {
		// Set default value if KVSTORE_DB is not set or invalid
		kvsDB = 0
	}
	cfg.KVSDB = kvsDB

	return cfg, nil
}
