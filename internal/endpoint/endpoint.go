// endpoints.go
package endpoint

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/SOAT1StackGoLang/msvc-payments/internal/service"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	Endpoint1 endpoint.Endpoint
	Endpoint2 endpoint.Endpoint
	Endpoint3 endpoint.Endpoint
	// Add other endpoints here
}

func MakeEndpoint1Handler(e endpoint.Endpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := service.Request1{} // Use the Request1 type from the service package
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := e(r.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func MakeEndpoint2Handler(e endpoint.Endpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := service.Request2{} // Use the Request2 type from the service package
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := e(r.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func MakeEndpoint3Handler(e endpoint.Endpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := service.Request3{} // Use the Request3 type from the service package
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := e(r.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		Endpoint1: makeEndpoint1(s),
		Endpoint2: makeEndpoint2(s),
		Endpoint3: makeEndpoint3(s),
		// Initialize other endpoints here
	}
}

func makeEndpoint1(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.Request1)
		resp, err := s.Endpoint1(ctx, req)
		return resp, err
	}
}

func makeEndpoint2(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.Request2)
		resp, err := s.Endpoint2(ctx, req)
		return resp, err
	}
}

func makeEndpoint3(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.Request3)
		resp, err := s.Endpoint3(ctx, req)
		return resp, err
	}
}

// Implement other makeEndpoint functions similarly
