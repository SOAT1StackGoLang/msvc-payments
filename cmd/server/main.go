package main

import (
	"os"

	"github.com/SOAT1StackGoLang/msvc-payments/internal/endpoint"
	"github.com/SOAT1StackGoLang/msvc-payments/internal/service"
	"github.com/SOAT1StackGoLang/msvc-payments/internal/transport"
	logger "github.com/SOAT1StackGoLang/msvc-payments/pkg/middleware"
)

// main is the entry point of the program.
// It initializes the Redis store, creates the service, sets up the endpoints,
// creates an HTTP handler, and starts the HTTP server.
func main() {
	redisStore, err := initializeApp()
	if err != nil {
		os.Exit(1)
	}

	// Create the service
	svc := service.NewService(redisStore)

	// Start processing payments background service
	go svc.StartProcessingPayments()

	// Create the endpoints using MakeEndpoints and CreatePaymentEndpoint from the service package
	endpoints := endpoint.MakeEndpoints(svc)

	httpHandler := transport.NewHTTPHandler(endpoints)

	// Start the HTTP server
	logger.Info("Starting HTTP server...")
	transport.NewHTTPServer(":8080", httpHandler)
}
