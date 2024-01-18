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
	redisStore, err := datastore.NewRedisStore("localhost:6379", "", 0)
	if err != nil {
		// handle error
		log.Println(err)
	}

	svc := service.NewService(redisStore)
	endpoints := endpoint.MakeEndpoints(svc)

	httpHandler := transport.NewHTTPHandler(endpoints)

	// Start the HTTP server
	transport.NewHTTPServer(":8080", httpHandler)
}
