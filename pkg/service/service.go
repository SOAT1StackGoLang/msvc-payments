package service

import (
	"context"

	"github.com/SOAT1StackGoLang/msvc-payments/pkg/datastore"
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
	Endpoint1(ctx context.Context, request Request1) (Response1, error)
	Endpoint2(ctx context.Context, request Request2) (Response2, error)
	Endpoint3(ctx context.Context, request Request3) (Response3, error)
}

type service struct {
	redisClient *datastore.RedisStore
}

func NewService(redisStore *datastore.RedisStore) Service {
	return &service{redisClient: redisStore}
}

// Implement the Service interface here
func (s *service) Endpoint1(ctx context.Context, request Request1) (Response1, error) {
	// store the message in the datastore
	err := s.redisClient.Set(ctx, request.Id, request.Message, 0)
	if err != nil {
		return Response1{}, err
	}
	get, err := s.redisClient.Client.Get(ctx, request.Id).Result()
	if err != nil {
		return Response1{}, err
	}
	// parse get to response
	return Response1{Id: request.Id, Message: get}, nil
}

func (s *service) Endpoint2(ctx context.Context, request Request2) (Response2, error) {
	// publish the message to the channel
	err := s.redisClient.Publish(ctx, "log", request.Message)
	if err != nil {
		return Response2{}, err
	}
	return Response2{}, nil
}

func (s *service) Endpoint3(ctx context.Context, request Request3) (Response3, error) {

	return Response3{}, nil
	// ...
}
