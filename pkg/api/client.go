package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	kitlog "github.com/go-kit/log"
	"net/http"
	"strings"
)

type client struct {
	baseURL string
	logger  kitlog.Logger
}

//go:generate mockgen -destination=../mocks/api_mocks.go -package=mocks github.com/SOAT1StackGoLang/msvc-payments/pkg/api PaymentAPI
type PaymentAPI interface {
	CreatePayment(request CreatePaymentRequest) (CreatePaymentResponse, error)
	GetPayment(request GetPaymentRequest) (GetPaymentResponse, error)
}

func NewClient(baseURL string, logger kitlog.Logger) PaymentAPI {
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "http://" + baseURL
	}
	return &client{
		baseURL: baseURL,
		logger:  logger,
	}
}

func (c *client) CreatePayment(request CreatePaymentRequest) (CreatePaymentResponse, error) {

	url := fmt.Sprintf("%s/payments", c.baseURL)

	payload, err := json.Marshal(request)
	if err != nil {
		return CreatePaymentResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return CreatePaymentResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return CreatePaymentResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CreatePaymentResponse{}, errors.New("unexpected status code")
	}

	var responseBody CreatePaymentResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return CreatePaymentResponse{}, err
	}

	return responseBody, nil
}

func (c *client) GetPayment(request GetPaymentRequest) (GetPaymentResponse, error) {
	url := fmt.Sprintf("%s/payments", c.baseURL)

	payload, err := json.Marshal(request)
	if err != nil {
		return GetPaymentResponse{}, err
	}

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(payload))
	if err != nil {
		return GetPaymentResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return GetPaymentResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GetPaymentResponse{}, errors.New("unexpected status code")
	}

	var responseBody GetPaymentResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return GetPaymentResponse{}, err
	}

	return responseBody, nil
}
