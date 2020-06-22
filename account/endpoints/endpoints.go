package endpoints

import (
	"context"
	"oncekey/go-kit-example/account/service"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints endpoints
type Endpoints struct {
	CreateUser endpoint.Endpoint
	GetUser    endpoint.Endpoint
}

// MakeEndpoints  make a endpoints for Service
func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		CreateUser: makeCreateUserEndpoint(s),
		GetUser:    makeGetUserEndpoint(s),
	}
}

func makeCreateUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)
		ok, err := s.CreateUser(ctx, req.Email, req.Password)

		return CreateUserResponse{
			OK: ok,
		}, err
	}
}

func makeGetUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)
		email, err := s.GetUser(ctx, req.ID)

		return GetUserResponse{
			Email: email,
		}, err
	}
}
