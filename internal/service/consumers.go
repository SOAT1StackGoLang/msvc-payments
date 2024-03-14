package service

import (
	"context"
	"encoding/json"
	"github.com/SOAT1StackGoLang/msvc-payments/pkg/messages"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	logger "github.com/SOAT1StackGoLang/msvc-payments/pkg/middleware"
	"github.com/google/uuid"
)

// initPaymentProccess and updatePaymentStatus are the functions that will be used by the goroutines
// process payments from the payments pending queue using RPopLPush that will use ProcessPayment be used by a goroutine

type BackgroundService interface {
	StartProcessingPayments()
}

func (s *serviceImpl) paymentProccess(ctx context.Context) error {
	logger.Info("Initializing payments processing...")

	// Listen for a shutdown signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-shutdown:
			// Stop processing payments and return
			logger.Info("Shutting down payment processing...")
			return nil
		default:
			// Process a payment
			payment_id, err := s.processPayment(ctx)
			if err != nil {
				logger.Error("Error while processing payments: ", err.Error())
				// Retry the operation with exponential backoff
				for i := 0; i < 3; i++ {
					time.Sleep(time.Second * time.Duration(math.Pow(2, float64(i))))
					payment_id, err = s.processPayment(ctx)
					if err == nil {
						break
					}
				}
				// If the operation still fails, move the payment to a dead-letter queue
				if err != nil {
					err = s.redisClient.LPush(ctx, "payments_deadletter", payment_id)
					if err != nil {
						logger.Error("Error while moving payment to dead-letter queue: ", err.Error())
					}
					// remove from the payments processing queue
					err = s.redisClient.LREM(ctx, "payments_processing", 0, payment_id)
					if err != nil {
						logger.Error("Error while cleaning up payment: ", err.Error())
					}
				}

			}
		}
	}
}
func (s *serviceImpl) processPayment(ctx context.Context) (string, error) {
	logger.Info("Initializing payments processing...")
	for {
		// get the payment from the payments pending queue

		payment_id, err := s.redisClient.BLMOVE(ctx, "payment_pending_queue", "payments_processing")
		if err != nil {
			logger.Error("Error while processing payments: ", err.Error())
			return "", err
		}
		// validate if the payment is valid uuid.UUID
		var payment_id_valid uuid.UUID
		payment_id_valid, err = uuid.Parse(payment_id)
		if err != nil {
			logger.Error("Error while parsing payment id: ", err.Error())
			return "", err
		}
		// process the payment
		payment, err := s.ProcessPayment(ctx, payment_id_valid)
		if err != nil {
			logger.Error("Error while processing payment: ", err.Error())
			return "", err
		}
		// update the payment status
		_, err = s.UpdatePayment(ctx, UpdatePaymentRequest{PaymentID: payment_id_valid, PaymentStatus: payment.Status})
		if err != nil {
			logger.Error("Error while updating payment: ", err.Error())
			return "", err
		}
		// cleanup the payment from the payments processing queue
		err = s.redisClient.LREM(ctx, "payments_processing", 0, payment.ID.String())
		if err != nil {
			logger.Error("Error while cleaning up payment: ", err.Error())
			return "", err
		}
		// return the payment id
		return payment.ID.String(), nil
	}
}

func (s *serviceImpl) StartProcessingPayments() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		retryInterval := time.Second * 5 // Retry every 5 seconds
		maxRetries := 3                  // Maximum number of retries

		for i := 0; i < maxRetries; i++ {
			err := s.paymentProccess(context.Background())
			if err != nil {
				logger.Info("Error while processing payments: ", err.Error())
				logger.Info("Retrying in ", retryInterval.String(), " seconds...")

				// Wait for retryInterval before retrying
				timer := time.NewTimer(retryInterval)
				<-timer.C // Directly receive from the timer's channel
			} else {
				// If no error, break the loop
				break
			}
		}
	}()

	wg.Wait()
}

func (s *serviceImpl) StartConsumingPayments() {
	ctx := context.Background()
	sub, err := s.redisClient.Subscribe(ctx, messages.OrderPaymentCreationRequestChannel)
	if err != nil {
		logger.Error("failed subscribing to payment creation requests")
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-sub:
			s.handlePaymentCreationRequest(msg.Payload)
		}
	}
}

func (s *serviceImpl) handlePaymentCreationRequest(payload string) {
	var paymentRequest messages.PaymentCreationRequestMessage
	err := json.Unmarshal([]byte(payload), &paymentRequest)
	if err != nil {
		logger.Error("failed unmarshalling payment creation request")
		return
	}

	pR, err := PaymentFromPaymentCreationRequestMessage(paymentRequest)
	if err != nil {
		logger.Error("failed converting payment creation request")
		return
	}

	if pR.Status == PaymentStatusClosed {
		_, err := s.UpdatePayment(context.Background(), UpdatePaymentRequest{
			PaymentID:     pR.ID,
			PaymentStatus: PaymentStatusClosed,
		})
		if err != nil {
			logger.Error("failed updating payment")
		}

		return
	}
	_, err = s.CreatePayment(context.Background(), CreatePaymentRequest{Payment: Payment{
		ID:        pR.ID,
		CreatedAt: pR.CreatedAt,
		UpdatedAt: pR.UpdatedAt,
		Price:     pR.Price,
		OrderID:   pR.OrderID,
		Status:    pR.Status,
	}})
	if err != nil {
		logger.Error("failed creating payment")
		return
	}
}
