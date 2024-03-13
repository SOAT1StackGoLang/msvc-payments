package messages

import (
	"github.com/SOAT1StackGoLang/msvc-payments/internal/service"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type PaymentCreationRequestMessage struct {
	ID        string  `json:"id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Price     float64 `json:"price"`
	OrderID   string  `json:"order_id"`
	Status    string  `json:"status"`
}

func (p PaymentCreationRequestMessage) ToCreatePaymentRequest() (*service.CreatePaymentRequest, error) {
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

	status := orderStatusToPaymentStatus(p.Status)

	return &service.CreatePaymentRequest{
		Payment: service.Payment{
			ID:        id,
			CreatedAt: createdAt,
			UpdatedAt: time.Now(),
			Price:     decimal.NewFromFloat(p.Price),
			OrderID:   orderID,
			Status:    status,
		},
	}, err
}

func orderStatusToPaymentStatus(status string) service.PaymentStatus {
	switch status {
	case "Aberto":
		return service.PaymentStatusPending
	case "Aguardando Pagamento":
		return service.PaymentStatusPending
	case "Recebido":
		return service.PaymentStatusPaid
	default:
		return service.PaymentStatusClosed
	}
}

func paymentStatusToOrderStatus(status service.PaymentStatus) string {
	switch status {
	case service.PaymentStatusPending:
		return "Aguardando Pagamento"
	case service.PaymentStatusPaid:
		return "Recebido"
	case service.PaymentStatusFailed:
		return "Falha no Pagamento"
	default:
		return "Fechado"
	}

}

type PaymentStatusChangedMessage struct {
	ID        string `json:"id"`
	OrderID   string `json:"order_id"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
}

func ToPaymentStatusChangedMessage(p service.Payment) *PaymentStatusChangedMessage {
	return &PaymentStatusChangedMessage{
		ID:        p.ID.String(),
		OrderID:   p.OrderID.String(),
		Status:    paymentStatusToOrderStatus(p.Status),
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
	}
}
