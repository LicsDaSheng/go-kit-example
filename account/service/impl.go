package service

import (
	"context"
	"oncekey/go-kit-example/account/dao"
	"oncekey/go-kit-example/account/model"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gofrs/uuid"
)

type service struct {
	dao    dao.Dao
	logger log.Logger
}

// NewService create Service with Res
func NewService(dao dao.Dao, logger log.Logger) Service {
	return &service{
		dao:    dao,
		logger: logger,
	}
}

func (s service) CreateUser(ctx context.Context, email string, password string) (string, error) {
	logger := log.With(s.logger, "method", "CreateUser")
	uuid, _ := uuid.NewV4()
	id := uuid.String()
	logger.Log("uuid", uuid)
	user := model.User{
		ID:       id,
		Email:    email,
		Password: password,
	}
	if err := s.dao.CreateUser(ctx, user); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}
	logger.Log("create user", id)
	return "Success", nil

}

func (s service) GetUser(ctx context.Context, id string) (string, error) {
	logger := log.With(s.logger, "method", "GetUser")

	email, err := s.dao.GetUser(ctx, id)

	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("Get User", id)

	return email, nil
}
