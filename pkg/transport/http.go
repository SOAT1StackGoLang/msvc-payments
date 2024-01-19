// http.go
package transport

import (
	"net/http"

	"github.com/SOAT1StackGoLang/msvc-payments/pkg/endpoint"
	"github.com/gorilla/mux"
)

func NewHTTPHandler(endpoints endpoint.Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Methods("POST").Path("/endpoint1").Handler(endpoint.MakeEndpoint1Handler(endpoints.Endpoint1))
	r.Methods("POST").Path("/endpoint2").Handler(endpoint.MakeEndpoint2Handler(endpoints.Endpoint2))
	r.Methods("GET").Path("/endpoint3").Handler(endpoint.MakeEndpoint3Handler(endpoints.Endpoint3))
	// Add other endpoints here
	return r
}

// NewHTTPServer starts an HTTP server with the given handler
func NewHTTPServer(addr string, handler http.Handler) {
	http.ListenAndServe(addr, handler)
}
