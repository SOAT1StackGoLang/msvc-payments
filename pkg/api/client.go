package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	kitlog "github.com/go-kit/log"
	"net/http"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     kitlog.Logger
}

//go:generate mockgen -destination=../mocks/api_mocks.go -package=mocks github.com/SOAT1StackGoLang/msvc-payments/pkg/api PaymentAPI
type PaymentAPI interface {
	CreatePayment(request CreatePaymentRequest) (CreatePaymentResponse, error)
	GetPayment(request GetPaymentRequest) (GetPaymentResponse, error)
}

func NewClient(baseURL string, httpClient *http.Client, logger kitlog.Logger) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
		logger:     logger,
	}
}

func (c *Client) CreatePayment(request CreatePaymentRequest) (CreatePaymentResponse, error) {
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

	resp, err := c.httpClient.Do(req)
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

func (c *Client) GetPayment(request GetPaymentRequest) (GetPaymentResponse, error) {
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

	resp, err := c.httpClient.Do(req)
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
