// endpoints.go
package endpoint

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SOAT1StackGoLang/msvc-payments/internal/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type Endpoints struct {
	CreatePayment endpoint.Endpoint
	GetPayment    endpoint.Endpoint
	UpdatePayment endpoint.Endpoint
	// Add other endpoints here
}

// Implement MakeCreatePaymentHandler
func MakeCreatePaymentHandler(e endpoint.Endpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		request := service.CreatePaymentRequest{} // Use the CreatePaymentRequest type from the service package
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if request.Payment.ID == uuid.Nil {
			err = errors.New("error on decoding request body")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := e(r.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Cast the response to the CreatePaymentResponse type from the service package
		createPaymentResponse := response.(service.CreatePaymentResponse)

		// Encode the response
		if err := json.NewEncoder(w).Encode(createPaymentResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Implement MakeGetPaymentHandler
func MakeGetPaymentHandler(e endpoint.Endpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := service.GetPaymentRequest{} // Use the GetPaymentRequest type from the service package
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if request.PaymentID == uuid.Nil {
			err := errors.New("error on decoding request body")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := e(r.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Cast the response to the GetPaymentResponse type from the service package
		getPaymentResponse := response.(service.GetPaymentResponse)

		// Encode the response
		if err := json.NewEncoder(w).Encode(getPaymentResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Implement MakeUpdatePaymentHandler
func MakeUpdatePaymentHandler(e endpoint.Endpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := service.UpdatePaymentRequest{} // Use the UpdatePaymentRequest type from the service package
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := e(r.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Cast the response to the UpdatePaymentResponse type from the service package
		updatePaymentResponse := response.(service.UpdatePaymentResponse)

		// Encode the response
		if err := json.NewEncoder(w).Encode(updatePaymentResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		CreatePayment: makeCreatePaymentEndpoint(s),
		GetPayment:    makeGetPaymentEndpoint(s),
		UpdatePayment: makeUpdatePaymentEndpoint(s),
		// Initialize other endpoints here
	}
}

// Implement makeCreatePaymentEndpoint
func makeCreatePaymentEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.CreatePaymentRequest)
		resp, err := s.CreatePayment(ctx, req)
		return resp, err
	}
}

// Implement makeGetPaymentEndpoint
func makeGetPaymentEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.GetPaymentRequest)
		resp, err := s.GetPayment(ctx, req)
		return resp, err
	}
}

// Implement makeUpdatePaymentEndpoint
func makeUpdatePaymentEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.UpdatePaymentRequest)
		resp, err := s.UpdatePayment(ctx, req)
		return resp, err
	}
}
