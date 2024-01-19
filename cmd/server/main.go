// main.go
package main

import (
	"log"

	"github.com/SOAT1StackGoLang/msvc-payments/pkg/datastore"
	"github.com/SOAT1StackGoLang/msvc-payments/pkg/endpoint"
	"github.com/SOAT1StackGoLang/msvc-payments/pkg/service"
	"github.com/SOAT1StackGoLang/msvc-payments/pkg/transport"
)

func main() {
	// Load the configuration
	log.Println("Loading configuration...")
	configs, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connecting to datastore...")
	redisStore, err := datastore.NewRedisStore(configs.KVSURI, "", 0)
	if err != nil {
		// handle error
		log.Println(err)
	}

	// Create the service
	svc := service.NewService(redisStore)
	endpoints := endpoint.MakeEndpoints(svc)

	httpHandler := transport.NewHTTPHandler(endpoints)

	// Start the HTTP server
	log.Println("Starting HTTP server...")
	transport.NewHTTPServer(":8080", httpHandler)
}
