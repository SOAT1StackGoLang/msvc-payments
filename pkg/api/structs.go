package api

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Payment struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Price     decimal.Decimal
	OrderID   uuid.UUID
	Status    PaymentStatus
}

type PaymentStatus string

const (
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusFailed  PaymentStatus = "failed"
	// TO-DO PaymentStatusClosed  PaymentStatus = "closed"
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
