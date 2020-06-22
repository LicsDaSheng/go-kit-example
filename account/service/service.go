package service

import "context"

// Service service
type Service interface {
	CreateUser(ctx context.Context, emial string, password string) (string, error)
	GetUser(ctx context.Context, id string) (string, error)
}
