package dao

import (
	"context"
	"database/sql"
	"errors"
	"oncekey/go-kit-example/account/model"

	"github.com/go-kit/kit/log"
)

// ErrRepo err an error with repository error
var ErrRepo = errors.New("Unable to handle Repo Request")

// Dao dao
type Dao interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, id string) (string, error)
}

// Dao dao层处理
type dao struct {
	db     *sql.DB
	logger log.Logger
}

// NewDao create a repository with db and logger ,return repository
func NewDao(db *sql.DB, logger log.Logger) Dao {
	return &dao{
		db:     db,
		logger: logger,
	}
}

// CreateUser 创建用户
func (d dao) CreateUser(ctx context.Context, user model.User) error {
	sql := `
		INSERT INTO users (id,email,password)
		VALUES($1, $2, $3)
	`
	if user.Email == "" || user.Password == "" {
		return ErrRepo
	}

	_, err := d.db.ExecContext(ctx, sql, user.ID, user.Email, user.Password)

	if err != nil {
		return err
	}
	return nil
}

// GetUser 获取用户
func (d dao) GetUser(ctx context.Context, id string) (string, error) {
	var email string
	err := d.db.QueryRow("SELECT email FROM users WHERE id = $1", id).Scan(&email)
	if err != nil {
		return "", ErrRepo
	}
	return email, nil
}
