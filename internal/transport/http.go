// http.go
package transport

import (
	"net/http"

	"github.com/SOAT1StackGoLang/msvc-payments/internal/endpoint"
	"github.com/gorilla/mux"
)

// NewHTTPHandler returns a new HTTP handler that routes incoming requests to the appropriate endpoints.
// It takes an `endpoints` parameter of type `endpoint.Endpoints` which contains the implementation of various endpoints.
// The handler is responsible for mapping the incoming HTTP requests to the corresponding endpoint functions.
// It returns an `http.Handler` that can be used to serve the HTTP requests.
func NewHTTPHandler(endpoints endpoint.Endpoints) http.Handler {
	r := mux.NewRouter()
	// Add other endpoints here

	// Create Payment endpoint
	r.Methods("POST").Path("/payments").Handler(endpoint.MakeCreatePaymentHandler(endpoints.CreatePayment))
	// Get Payment endpoint
	r.Methods("GET").Path("/payments").Handler(endpoint.MakeGetPaymentHandler(endpoints.GetPayment))
	// Update Payment endpoint
	r.Methods("PUT").Path("/payments").Handler(endpoint.MakeUpdatePaymentHandler(endpoints.UpdatePayment))
	return r
}

// NewHTTPServer creates a new HTTP server that listens on the specified address
// and handles requests using the provided handler.
//
// Parameters:
// - addr: The address to listen on (e.g., "localhost:8080").
// - handler: The http.Handler to handle incoming requests.
//
// Example usage:
//
//	NewHTTPServer("localhost:8080", myHandler)
//
// Note: This function blocks indefinitely, so it should typically be called in a
// separate goroutine.
func NewHTTPServer(addr string, handler http.Handler) {
	http.ListenAndServe(addr, handler)
}
