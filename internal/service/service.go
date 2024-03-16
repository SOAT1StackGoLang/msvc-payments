package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SOAT1StackGoLang/msvc-payments/pkg/messages"
	"math/rand"
	"time"

	"github.com/SOAT1StackGoLang/msvc-payments/pkg/datastore"
	logger "github.com/SOAT1StackGoLang/msvc-payments/pkg/middleware"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// define the Request and Response types here
type Request1 struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

type Response1 struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

type Request2 struct {
	Message string `json:"message"`
}

type Response2 struct {
	Message string `json:"message"`
}

type Request3 struct {
}

type Response3 struct {
}

type Service interface {
	CreatePayment(ctx context.Context, request CreatePaymentRequest) (CreatePaymentResponse, error)
	UpdatePayment(ctx context.Context, request UpdatePaymentRequest) (UpdatePaymentResponse, error)
	GetPayment(ctx context.Context, request GetPaymentRequest) (GetPaymentResponse, error)
	StartProcessingPayments()
	StartConsumingPaymentsRequests()
}

type serviceImpl struct {
	redisClient datastore.RedisStore
}

func NewService(redisStore datastore.RedisStore) Service {
	return &serviceImpl{redisClient: redisStore}
}

type Payment struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Price     decimal.Decimal
	OrderID   uuid.UUID
	Status    PaymentStatus
}

func PaymentStatusChangedMessageFromPayment(p Payment) messages.PaymentStatusChangedMessage {
	return messages.PaymentStatusChangedMessage{
		ID:        p.ID.String(),
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
		OrderID:   p.OrderID.String(),
		Status:    string(p.Status),
	}

}
func PaymentFromPaymentCreationRequestMessage(p messages.PaymentCreationRequestMessage) (*Payment, error) {
	id, err := uuid.Parse(p.ID)
	if err != nil {
		return nil, err
	}

	orderID, err := uuid.Parse(p.OrderID)
	if err != nil {
		return nil, err
	}

	createdAt, err := time.Parse(time.RFC3339, p.CreatedAt)
	if err != nil {
		return nil, err
	}

	status := paymentStatusFromOrderStatus(p.Status)

	return &Payment{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
		Price:     decimal.NewFromFloat(p.Price),
		OrderID:   orderID,
		Status:    status,
	}, nil
}

func paymentStatusFromOrderStatus(status string) PaymentStatus {
	switch status {
	case "Aberto":
		return PaymentStatusPending
	case "Aguardando Pagamento":
		return PaymentStatusPending
	case "Recebido":
		return PaymentStatusPaid
	default:
		return PaymentStatusClosed
	}
}

type PaymentStatus string

const (
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusFailed  PaymentStatus = "failed"
	PaymentStatusClosed  PaymentStatus = "closed"
)

type CreatePaymentRequest struct {
	Payment Payment `json:"payment"`
}

type CreatePaymentResponse struct {
	PaymentID uuid.UUID     `json:"payment_id"`
	Status    PaymentStatus `json:"status"`
}

type UpdatePaymentRequest struct {
	PaymentID     uuid.UUID     `json:"payment_id"`
	PaymentStatus PaymentStatus `json:"payment_status"`
}

type UpdatePaymentResponse struct {
	PaymentID    uuid.UUID     `json:"payment_id"`
	Status       PaymentStatus `json:"status"`
	PaymentError string        `json:"payment_error,omitempty"`
}

type GetPaymentRequest struct {
	PaymentID uuid.UUID `json:"payment_id"`
}

type GetPaymentResponse struct {
	Payment      Payment       `json:"payment"`
	Status       PaymentStatus `json:"status"`
	PaymentError string        `json:"payment_error,omitempty"`
}

// Implement the Service interface here

// CreatePayment creates a new payment
func (s *serviceImpl) CreatePayment(ctx context.Context, request CreatePaymentRequest) (CreatePaymentResponse, error) {
	// Validate UUID
	if request.Payment.ID == uuid.Nil {
		logger.Error("Invalid UUID")
		return CreatePaymentResponse{}, fmt.Errorf("invalid uuid")
	}
	// set the payment status to pending
	request.Payment.Status = PaymentStatusPending
	request.Payment.CreatedAt = time.Now()
	request.Payment.UpdatedAt = time.Now()
	// store the payment in the datastore
	// Convert the payment to a JSON string
	jsonString, err := json.Marshal(request.Payment)
	if err != nil {
		logger.Error(err.Error())
		return CreatePaymentResponse{}, err
	}
	// check if the payment already exists
	exists, err := s.redisClient.Exists(ctx, request.Payment.ID.String())
	if exists {
		logger.Error("Payment already exists")
		return CreatePaymentResponse{}, fmt.Errorf("payment already exists")
	} else if err != nil {
		logger.Error(err.Error())
		return CreatePaymentResponse{}, err
	}

	err = s.redisClient.Set(ctx, request.Payment.ID.String(), jsonString, 0)
	if err != nil {
		logger.Error(err.Error())
		return CreatePaymentResponse{}, err
	}
	// place the payment in the redis queue
	err = s.redisClient.LPush(ctx, "payment_pending_queue", request.Payment.ID.String())
	if err != nil {
		// delete the payment from the datastore
		pusherr := err
		delerr := s.redisClient.Delete(ctx, request.Payment.ID.String())
		if delerr != nil {
			// concat the errors
			err = fmt.Errorf("PUSH: %s ---- DELETE: %s", pusherr.Error(), delerr.Error())
			logger.Error(err.Error())
			return CreatePaymentResponse{}, err
		}
		return CreatePaymentResponse{}, err
	}
	return CreatePaymentResponse{PaymentID: request.Payment.ID, Status: PaymentStatusPending}, nil
}

// MockPaymentProcess is a function that simulates a payment process and returns a mock payment status.
// It takes a pointer to a PaymentStatus struct as input and returns a PaymentStatus value.
// The function generates a random payment status between paid and failed, with a preference for paid.
func MockPaymentProcess(p PaymentStatus) PaymentStatus {
	// return a random payment status between paid and failed but prefer paid
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if r.Float64() < 0.8 {
		return PaymentStatusPaid
	}
	return PaymentStatusFailed
}

// ProcessPayment processes a payment
func (s *serviceImpl) ProcessPayment(ctx context.Context, paymentID uuid.UUID) (Payment, error) {
	// get the payment from the datastore
	paymentstored, err := s.redisClient.Get(ctx, paymentID.String())
	if err != nil {
		return Payment{}, err
	}
	// convert paymentstored to Payment type
	var payment Payment
	err = json.Unmarshal([]byte(paymentstored), &payment)
	if err != nil {
		logger.Error(fmt.Errorf("Error unmarshalling payment: %s", err.Error()).Error())
		return Payment{}, err
	}
	// simulate a payment process
	payment.Status = MockPaymentProcess(payment.Status)
	return payment, nil
}

// UpdatePayment updates a payment
func (s *serviceImpl) UpdatePayment(ctx context.Context, request UpdatePaymentRequest) (UpdatePaymentResponse, error) {
	// get the payment from the datastore
	if request.PaymentStatus == PaymentStatusClosed {
		paymentStored, err := s.redisClient.Get(ctx, request.PaymentID.String())
		if err != nil || paymentStored == "" {
			return UpdatePaymentResponse{}, nil
		}

		// Removing from queues
		_ = s.redisClient.LREM(ctx, "payment_pending_queue", 0, request.PaymentID.String())
		_ = s.redisClient.LREM(ctx, "payment_processing_queue", 0, request.PaymentID.String())

		// payment string to bytes
		paymentBytes, err := json.Marshal(paymentStored)

		err = s.redisClient.Set(ctx, request.PaymentID.String(), paymentBytes, 0)

		return UpdatePaymentResponse{
			PaymentID: request.PaymentID,
			Status:    request.PaymentStatus,
		}, nil
	}

	payment, err := s.ProcessPayment(ctx, request.PaymentID)
	paymentResp := PaymentStatusChangedMessageFromPayment(payment)
	pRespBytes, err := json.Marshal(paymentResp)
	if err != nil {
		return UpdatePaymentResponse{}, err
	}
	// convert payment
	paymentBytes, err := json.Marshal(payment)
	if err != nil {
		logger.Error(fmt.Errorf("Error marshalling payment: %s", err.Error()).Error())
		return UpdatePaymentResponse{}, err
	}
	// notify channels of the payment status
	switch payment.Status {
	case PaymentStatusPaid:
		err = s.redisClient.Publish(ctx, messages.PaymentStatusResponseChannel, pRespBytes)
		if err != nil {
			logger.Error(err.Error())
			return UpdatePaymentResponse{}, err
		}
		// add the payment to the paid queue
		err = s.redisClient.LPush(ctx, "payment_paid_queue", payment.ID.String())
		if err != nil {
			logger.Error(err.Error())
			return UpdatePaymentResponse{}, err
		}
	case PaymentStatusFailed:
		err = s.redisClient.Publish(ctx, messages.PaymentStatusResponseChannel, pRespBytes)
		if err != nil {
			logger.Error(err.Error())
			return UpdatePaymentResponse{}, err
		}
		// add the payment to the failed queue
		err = s.redisClient.LPush(ctx, "payment_failed_queue", payment.ID.String())
		if err != nil {
			logger.Error(err.Error())
			return UpdatePaymentResponse{}, err
		}
	}

	// store the payment in the datastore
	err = s.redisClient.Set(ctx, request.PaymentID.String(), paymentBytes, 0)
	if err != nil {
		return UpdatePaymentResponse{}, err
	}
	return UpdatePaymentResponse{PaymentID: request.PaymentID, Status: request.PaymentStatus}, nil
}

// GetPayment gets a payment
func (s *serviceImpl) GetPayment(ctx context.Context, request GetPaymentRequest) (GetPaymentResponse, error) {
	// get the payment from the datastore
	paymentStored, err := s.redisClient.Get(ctx, request.PaymentID.String())
	if err != nil {
		return GetPaymentResponse{}, err
	}
	// convert paymentstored to Payment type
	var payment Payment
	err = json.Unmarshal([]byte(paymentStored), &payment)
	if err != nil {
		logger.Error(fmt.Errorf("Error unmarshalling payment: %s", err.Error()).Error())
		return GetPaymentResponse{}, err
	}
	return GetPaymentResponse{Payment: payment, Status: payment.Status}, nil
}
