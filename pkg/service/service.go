package service

import (
	"context"

	"github.com/SOAT1StackGoLang/msvc-payments/pkg/datastore"
)

// define the Request and Response types here
type Request1 struct {
	Message string
}

type Response1 struct {
	Message string
}

type Request2 struct {
}

type Response2 struct {
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
	// lets mock something just to build
	return Response1{}, nil
	// ...
}

func (s *service) Endpoint2(ctx context.Context, request Request2) (Response2, error) {

	return Response2{}, nil
	// ...
}

func (s *service) Endpoint3(ctx context.Context, request Request3) (Response3, error) {

	return Response3{}, nil
	// ...
}
