package messages

type PaymentCreationRequestMessage struct {
	ID        string  `json:"id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Price     float64 `json:"price"`
	OrderID   string  `json:"order_id"`
	Status    string  `json:"status"`
}

type PaymentStatusChangedMessage struct {
	ID        string `json:"id"`
	OrderID   string `json:"order_id"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
}
